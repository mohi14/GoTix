package server

import (
	"fmt"
	"gotickets/internal/config"
	"gotickets/internal/domain/booking"
	"gotickets/internal/domain/event"
	"gotickets/internal/domain/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func StartServer(db *gorm.DB, cfg *config.Config) {
	db.AutoMigrate(&user.User{}, &event.Event{}, &booking.Booking{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Running")
	})

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})

	//routes
	user.RegisterRoutes(e, db, cfg)
	event.RegisterRoutes(e, db)
	booking.RegisterRoutes(e, db, cfg)

	port := fmt.Sprintf(":%s", cfg.PORT)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
