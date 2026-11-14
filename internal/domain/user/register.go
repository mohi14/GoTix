package user

import (
	"gotickets/internal/auth"
	"gotickets/internal/config"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	userRepository := NewRepository(db)
	jwtService := auth.NewJWTService(cfg.JWTSECRET)
	userService := NewService(userRepository, jwtService)
	userHandler := NewHandler(userService)

	api := e.Group("/api/v1/auth")

	api.POST("/register", userHandler.CreateUser)
	api.POST("/login", userHandler.LoginUser)

}
