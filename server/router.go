package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	if strings.HasPrefix(c.Path(), "/internal_api/validate") {
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
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
	}))

	// Routes
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.GET("/check-time", func(c echo.Context) error {
		return c.String(http.StatusOK, strconv.FormatInt(time.Now().Unix(), 10))
	})
	//e.POST("/login", AuthenticateHandler(s))
	api := e.Group("/internal_api")
	{
		// api
		api.POST("/object", CreateSimpleObjectHandler(s))
		api.PUT("/object", UpdateObjectHandler(s))
		api.GET("/object", GetCertForServiceHandler(s))
		api.GET("/validate", ValidateObjectHandler(s))
		api.POST("/cert", GetCertForClientHandler(s))
		api.POST("/policy", CreatePolicyHandler(s))
		api.PUT("/policy", UpdatePolicyHandler(s))
		api.POST("/authorize", AuthenticateHandler(s))
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

				ctx := c.Request().Context()

				ctx2 := context.WithValue(ctx, ContextKeyTransaction, tx)

				c.SetRequest(c.Request().WithContext(ctx2))

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

				ctx := c.Request().Context()

				ctx2 := context.WithValue(ctx, ContextKeyTransaction, txi)

				c.SetRequest(c.Request().WithContext(ctx2))

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
			log.Logger().Error(err)
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

func UpdatePolicyHandler(s *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type UpdatePolicyRequest struct {
			Id         string         `json:"id" validate:"required"`
			Token      string         `json:"token" validate:"required"`
			Name       string         `json:"name"`
			ServiceID  string         `json:"service_id"`
			Permission map[string]int `json:"permission"`
			Status     string         `json:"status"`
			ApplyFrom  int64          `json:"apply_from"`
			ApplyTo    int64          `json:"apply_to"`
		}

		type UpdatePolicyResponse struct {
			Code    int               `json:"code"`
			Message string            `json:"message"`
			Policy  *gokontrol.Policy `json:"policy"`
		}

		pr := new(UpdatePolicyRequest)
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
			ID:         pr.Id,
			Name:       pr.Name,
			ServiceID:  pr.ServiceID,
			Permission: pr.Permission,
			Status:     pr.Status,
			ApplyFrom:  pr.ApplyFrom,
			ApplyTo:    pr.ApplyTo,
		}
		err := s.Kontrol.UpdatePolicy(c.Request().Context(), pr.Token, policy)
		if err != nil {
			log.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, UpdatePolicyResponse{Code: http.StatusOK, Message: "ok", Policy: policy})
	}
}

//ValidateObjectHandler quick check if token is valid
func ValidateObjectHandler(s *Service) echo.HandlerFunc {
	return func(c echo.Context) error {

		type ValidateObjectResponse struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		if c.Request().Header.Get("X-Forwarded-Method") == http.MethodOptions {
			return c.JSON(http.StatusOK, ValidateObjectResponse{Code: http.StatusOK, Message: "ok"})
		}

		// verify Access-Token header exist
		if _, ok := c.Request().Header["Authorization"]; !ok {
			return c.JSON(http.StatusUnauthorized, errors.New("Header 'Authorization' is empty "))
		}
		reqToken := c.Request().Header["Authorization"][0]
		reqToken = strings.Trim(strings.Replace(reqToken, "Bearer", "", 1), " ")

		_, err := s.Kontrol.ValidateToken(c.Request().Context(), reqToken, c.Request().Header["X-Forwarded-Uri"][0], c.Request().Method)
		if err != nil {
			log.Logger().Debug(err)
			return c.JSON(http.StatusForbidden, CommonError.FORBIDDEN)
		}
		return c.JSON(http.StatusOK, ValidateObjectResponse{Code: http.StatusOK, Message: "ok"})
	}
}

//GetCertForClientHandler return object permission after successful authen
func GetCertForClientHandler(s *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type GetCertForClientRequest struct {
			ObjectID  string `json:"object_id" validate:"required"`
			ServiceID string `json:"service_id" validate:"required"`
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
			log.Logger().Error(err)
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

//AuthenticateHandler Authenticate user --> call REST API cert to get request
func AuthenticateHandler(s *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type AuthenticateRequest struct {
			ServiceID string `json:"service_id" validate:"required"`
			UserName  string `json:"user_name" validate:"required"`
			Password  string `json:"password" validate:"required"`
		}

		type AuthenticateResponse struct {
			Code             int                         `json:"code"`
			Message          string                      `json:"message"`
			ObjectPermission *gokontrol.ObjectPermission `json:"object_permission"`
		}
		type User struct {
			ExternalId string `json:"external_id"`
			UserName   string `json:"user_name"`
			Password   string `json:"password"`
		}
		// this is mock user for demo authenticate step in external service
		users := map[string]User{
			// adt - user
			"adtuser1":  {ExternalId: "adt_id_1", UserName: "adtuser1", Password: "pass1"},
			"adtuser2":  {ExternalId: "adt_id_2", UserName: "adtuser2", Password: "pass2"},
			"adtuser3":  {ExternalId: "adt_id_3", UserName: "adtuser3", Password: "pass3"},
			"adtuser4":  {ExternalId: "adt_id_4", UserName: "adtuser4", Password: "pass4"},
			"adtuser5":  {ExternalId: "adt_id_5", UserName: "adtuser5", Password: "pass5"},
			"adtuser6":  {ExternalId: "adt_id_6", UserName: "adtuser6", Password: "pass6"},
			"adtuser7":  {ExternalId: "adt_id_7", UserName: "adtuser7", Password: "pass7"},
			"adtuser8":  {ExternalId: "adt_id_8", UserName: "adtuser8", Password: "pass8"},
			"adtuser9":  {ExternalId: "adt_id_9", UserName: "adtuser9", Password: "pass9"},
			"adtuser10": {ExternalId: "adt_id_19", UserName: "adtuser10", Password: "pass10"},

			//idt user for login
			"idtuser1":  {ExternalId: "idt_id_1", UserName: "idtuser1", Password: "pass1"},
			"idtuser2":  {ExternalId: "idt_id_2", UserName: "idtuser2", Password: "pass2"},
			"idtuser3":  {ExternalId: "idt_id_3", UserName: "idtuser3", Password: "pass3"},
			"idtuser4":  {ExternalId: "idt_id_4", UserName: "idtuser4", Password: "pass4"},
			"idtuser5":  {ExternalId: "idt_id_5", UserName: "idtuser5", Password: "pass5"},
			"idtuser6":  {ExternalId: "idt_id_6", UserName: "idtuser6", Password: "pass6"},
			"idtuser7":  {ExternalId: "idt_id_7", UserName: "idtuser7", Password: "pass7"},
			"idtuser8":  {ExternalId: "idt_id_8", UserName: "idtuser8", Password: "pass8"},
			"idtuser9":  {ExternalId: "idt_id_9", UserName: "idtuser9", Password: "pass9"},
			"idtuser10": {ExternalId: "idt_id_19", UserName: "idtuser10", Password: "pass10"},

			//hrd user for login
			"hrduser1":  {ExternalId: "hrd_id_1", UserName: "hrduser1", Password: "pass1"},
			"hrduser2":  {ExternalId: "hrd_id_2", UserName: "hrduser2", Password: "pass2"},
			"hrduser3":  {ExternalId: "hrd_id_3", UserName: "hrduser3", Password: "pass3"},
			"hrduser4":  {ExternalId: "hrd_id_4", UserName: "hrduser4", Password: "pass4"},
			"hrduser5":  {ExternalId: "hrd_id_5", UserName: "hrduser5", Password: "pass5"},
			"hrduser6":  {ExternalId: "hrd_id_6", UserName: "hrduser6", Password: "pass6"},
			"hrduser7":  {ExternalId: "hrd_id_7", UserName: "hrduser7", Password: "pass7"},
			"hrduser8":  {ExternalId: "hrd_id_8", UserName: "hrduser8", Password: "pass8"},
			"hrduser9":  {ExternalId: "hrd_id_9", UserName: "hrduser9", Password: "pass9"},
			"hrduser10": {ExternalId: "hrd_id_19", UserName: "hrduser10", Password: "pass10"},
		}
		pr := new(AuthenticateRequest)
		c.Bind(pr)
		if err := c.Validate(pr); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		// authenticate -- for demo :)
		if _, ok := users[pr.UserName]; ok == true {
			if users[pr.UserName].Password != pr.Password {
				return c.JSON(http.StatusForbidden, errors.New("Invalid username or password "))
			}
		} else {
			return c.JSON(http.StatusForbidden, errors.New("User is not existed "))
		}
		cert, err := s.getServerCert(pr.ServiceID, users[pr.UserName].ExternalId)
		if err != nil {
			log.Logger().Error(err)
			return c.JSON(http.StatusUnprocessableEntity, err)
		}
		return c.JSON(http.StatusOK, AuthenticateResponse{Code: http.StatusOK, Message: "ok", ObjectPermission: cert})
	}
}

// getCert call API `cert` to token  and permissions
func (s *Service) getServerCert(serviceId, externalId string) (*gokontrol.ObjectPermission, error) {
	type GetCertForClientResponse struct {
		Code             int                         `json:"code"`
		Message          string                      `json:"message"`
		ObjectPermission *gokontrol.ObjectPermission `json:"object_permission"`
	}
	type GetCertForClientRequest struct {
		ObjectID  string `json:"object_id" validate:"required"`
		ServiceID string `json:"service_id" validate:"required"`
	}
	data := GetCertForClientRequest{ObjectID: externalId, ServiceID: serviceId}
	bodyData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(fmt.Sprintf("http://sso_service:%s/internal_api/cert", s.Config.HTTPPort), "application/json", bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var apiResponse GetCertForClientResponse
	err = json.Unmarshal(responseData, &apiResponse)
	if err != nil {
		return nil, err
	}
	if apiResponse.Code == http.StatusOK {
		return apiResponse.ObjectPermission, nil
	} else {
		return nil, errors.New(apiResponse.Message)
	}
}
