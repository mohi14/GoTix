package user

import (
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	userRepository := NewRepository(db)
	jwtService := auth.NewJWTService(cfg.JWTSECRET)
	userService := NewService(userRepository, jwtService, cfg)
	userHandler := NewHandler(userService)

	authMiddleware := middlewares.AuthMiddleware(jwtService)

	api := e.Group("/api/v1/auth")

	api.POST("/register", userHandler.CreateUser)
	api.POST("/login", userHandler.LoginUser)
	api.POST("/google", userHandler.LoginWithGoogle)
	api.GET("/me", userHandler.GetMe, authMiddleware) // protected route, requires authentication
	api.PATCH("/users/:id/role", userHandler.UpdateRole, authMiddleware, middlewares.RequireRoles(RoleAdmin, RoleSuperAdmin))

}
