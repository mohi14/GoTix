package event

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/middlewares"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	api := e.Group("/api/v1/events")
	authMiddleware := middlewares.AuthMiddleware(auth.NewJWTService(cfg.JWTSECRET))
	eventManagers := middlewares.RequireRoles("organizer", "admin", "super_admin")

	api.POST("", handler.CreateEvent, authMiddleware, eventManagers)
	api.GET("", handler.GetEvents)
	api.GET("/:id", handler.GetEventsByID)
	api.PUT("/:id", handler.UpdateEvent, authMiddleware, eventManagers)

}
