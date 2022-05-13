package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	gokontrol "github.com/LibertusDio/go-kontrol"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/neko-neko/echo-logrus/v2/log"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

func urlSkipper(c echo.Context) bool {
	if strings.HasPrefix(c.Path(), "/health") {
		return true
	}
	if strings.HasPrefix(c.Path(), "/metrics") {
		return true
	}
	if strings.HasPrefix(c.Path(), "/check-time") {
		return true
	}

	return false
}

func NewEcho(s *Service) *echo.Echo {
	// Echo instance
	e := echo.New()
	e.Logger = s.Logger
	p := prometheus.NewPrometheus("echo", urlSkipper)
	p.Use(e)

	// Validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.Gzip())
	// Fetch new store.
	e.Use(GormTransactionHandler(s.DB))

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Routes
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.GET("/check-time", func(c echo.Context) error {
		return c.String(http.StatusOK, strconv.FormatInt(time.Now().Unix(), 10))
	})

	api := e.Group("/api")
	{
		// api
		api.POST("/object", CreateSimpleObjectHandler(s))
		api.PUT("/object", UpdateObjectHandler(s))
		api.GET("/object", GetCertForServiceHandler(s))
		api.GET("/validate", ValidateObjectHandler(s))
		api.GET("/cert", GetCertForClientHandler(s))
		api.POST("/policy", CreatePolicyHandler(s))
	}

	// admin	 := e.Group("/admin")
	// {
	// 	admin.
	// }
	return e
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func GormTransactionHandler(db Database) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			if c.Request().Method != "GET" {
				txi, _ := db.Transaction()
				tx := txi.(*gorm.DB)
				c.Set(ContextKeyTransaction, tx)

				if err := next(c); err != nil {
					tx.Rollback()
					log.Logger().Debug("Transaction Rollback: ", err)
					return err
				}
				log.Logger().Debug("Transaction Commit")
				tx.Commit()
			} else {
				txi, _ := db.Session()
				c.Set(ContextKeyTransaction, txi)
				return next(c)
			}

			return nil
		})
	}
}

func CreateSimpleObjectHandler(s *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type CreateSimpleObjectRequest struct {
			ObjectID  string `json:"object_id" validate:"required"`
			Token     string `json:"token" validate:"required"`
			ServiceID string `json:"service_id" validate:"required"`
		}

		type CreateSimpleObjectResponse struct {
			Code    int                        `json:"code"`
			Message string                     `json:"message"`
			Data    gokontrol.ObjectPermission `json:"object_permission"`
		}

		pr := new(CreateSimpleObjectRequest)
		c.Bind(pr)
		if err := c.Validate(pr); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		objcert, err := s.Kontrol.AddSimpleObjectWithDefaultPolicy(c.Request().Context(), pr.ObjectID, pr.ServiceID, pr.Token)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, err)
		}

		return c.JSON(http.StatusOK, CreateSimpleObjectResponse{Code: http.StatusOK, Message: "true", Data: *objcert})
	}
}

func UpdateObjectHandler(s *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type UpdateObjectRequest struct {
			ObjectID    string   `json:"object_id" validate:"required"`
			Token       string   `json:"token" validate:"required"`
			GlobalID    string   `json:"global_id"`
			ServiceID   string   `json:"service_id" validate:"required"`
			ExternalID  string   `json:"external_id" validate:"required"`
			Status      string   `json:"status" validate:"required"`
			ApplyPolicy []string `json:"apply_policy"`
		}

		type UpdateObjectResponse struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}

		pr := new(UpdateObjectRequest)
		c.Bind(pr)
		if err := c.Validate(pr); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		ap := make([]*gokontrol.Policy, 0)
		for _, pid := range pr.ApplyPolicy {
			p, err := s.StorageKontrol.GetPolicyByID(c.Request().Context(), pid)
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			ap = append(ap, p)
		}

		err := s.Kontrol.UpdateObject(c.Request().Context(), &gokontrol.Object{
			ID:          pr.ObjectID,
			GlobalID:    pr.GlobalID,
			ExternalID:  pr.ExternalID,
			ServiceID:   pr.ServiceID,
			Status:      pr.Status,
			ApplyPolicy: ap,
		}, pr.Token)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, err)
		}

		return c.JSON(http.StatusOK, UpdateObjectResponse{Code: http.StatusOK, Message: "ok"})
	}
}

func CreatePolicyHandler(s *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type CreatePolicyRequest struct {
			Token      string         `json:"token" validate:"required"`
			Name       string         `json:"name"`
			ServiceID  string         `json:"service_id"`
			Permission map[string]int `json:"permission"`
			Status     string         `json:"status"`
			ApplyFrom  int64          `json:"apply_from"`
			ApplyTo    int64          `json:"apply_to"`
		}

		type CreatePolicyResponse struct {
			Code    int               `json:"code"`
			Message string            `json:"message"`
			Policy  *gokontrol.Policy `json:"policy"`
		}

		pr := new(CreatePolicyRequest)
		c.Bind(pr)
		if err := c.Validate(pr); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		for _, v := range pr.Permission {
			if v < 0 || v > 2 {
				return c.JSON(http.StatusBadRequest, CommonError.INVALID_PARAM)
			}
		}
		policy := &gokontrol.Policy{
			ID:         uuid.NewString(),
			Name:       pr.Name,
			ServiceID:  pr.ServiceID,
			Permission: pr.Permission,
			Status:     pr.Status,
			ApplyFrom:  pr.ApplyFrom,
			ApplyTo:    pr.ApplyTo,
		}
		err := s.Kontrol.CreatePolicy(c.Request().Context(), pr.Token, policy)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, CreatePolicyResponse{Code: http.StatusOK, Message: "ok", Policy: policy})
	}
}

//ValidateObjectHandler quick check if token is valid
func ValidateObjectHandler(s *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type ValidateObjectRequest struct {
			Token     string `query:"token" validate:"required"`
			ServiceID string `query:"service_id" validate:"required"`
		}

		type ValidateObjectResponse struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}

		pr := new(ValidateObjectRequest)
		c.Bind(pr)
		if err := c.Validate(pr); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		_, err := s.Kontrol.ValidateToken(c.Request().Context(), pr.Token, pr.ServiceID)
		if err != nil {
			return c.JSON(http.StatusForbidden, CommonError.FORBIDDEN)
		}
		return c.JSON(http.StatusOK, ValidateObjectResponse{Code: http.StatusOK, Message: "ok"})
	}
}

//GetCertForClientHandler return object permission after successful authen
func GetCertForClientHandler(s *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type GetCertForClientRequest struct {
			ObjectID  string `query:"object_id" validate:"required"`
			ServiceID string `query:"service_id" validate:"required"`
		}

		type GetCertForClientResponse struct {
			Code             int                         `json:"code"`
			Message          string                      `json:"message"`
			ObjectPermission *gokontrol.ObjectPermission `json:"object_permission"`
		}

		pr := new(GetCertForClientRequest)
		c.Bind(pr)
		if err := c.Validate(pr); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		cert, err := s.Kontrol.IssueCertForClient(c.Request().Context(), pr.ObjectID, pr.ServiceID)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, err)
		}
		return c.JSON(http.StatusOK, GetCertForClientResponse{Code: http.StatusOK, Message: "ok", ObjectPermission: cert})
	}
}

//GetCertForServiceHandler return object permission for service to cache
func GetCertForServiceHandler(s *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type GetCertForClientRequest struct {
			ObjectID  string `query:"object_id" validate:"required"`
			ServiceID string `query:"service_id" validate:"required"`
		}

		type GetCertForClientResponse struct {
			Code             int                         `json:"code"`
			Message          string                      `json:"message"`
			ObjectPermission *gokontrol.ObjectPermission `json:"object_permission"`
		}

		pr := new(GetCertForClientRequest)
		c.Bind(pr)
		if err := c.Validate(pr); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		cert, err := s.Kontrol.IssueCertForService(c.Request().Context(), pr.ObjectID, pr.ServiceID)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, err)
		}
		return c.JSON(http.StatusOK, GetCertForClientResponse{Code: http.StatusOK, Message: "ok", ObjectPermission: cert})
	}
}
