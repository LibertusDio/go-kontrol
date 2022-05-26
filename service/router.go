package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/neko-neko/echo-logrus/v2/log"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

func NewEcho(s *Service) *echo.Echo {
	// Echo instance
	e := echo.New()
	e.Logger = s.Logger

	// Validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
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
	e.GET("/internal_api/info", func(c echo.Context) error {
		type InfoResponse struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		return c.JSON(http.StatusOK, InfoResponse{Code: http.StatusOK, Message: fmt.Sprintf("Welcome to %s service", s.Config.AppName)})
	})

	// admin := e.Group("/admin")
	// {
	// 	// admin
	// 	admin.GET("/check-permission", adminobjectpermission.GetCheckPermission(s))
	// 	admin.GET("/check-create", adminobjectpermission.GetCheckCreate(s))
	// 	admin.GET("/check-rights", adminobjectpermission.GetCheckRight(s))
	// 	admin.POST("/create", adminobjectpermission.PostCreateObject(s))
	// 	admin.POST("/chmod", adminobjectpermission.PostChmod(s))
	// 	admin.POST("/chown", adminobjectpermission.PostChown(s))
	// 	admin.POST("/check-access", adminaccesscontrol.PostCheckAccess(s))
	// 	admin.GET("/permission-list", adminobjectpermission.GetPermissionList(s))
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
