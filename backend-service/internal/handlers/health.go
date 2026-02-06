package handlers

import (
	"github.com/beanspect/backend-service/internal/config"
	"github.com/beanspect/backend-service/internal/database"
	"github.com/gofiber/fiber/v2"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status      string `json:"status"`
	Service     string `json:"service"`
	Version     string `json:"version"`
	DBConnected bool   `json:"db_connected"`
}

// Health returns the health status of the service
func Health(c *fiber.Ctx) error {
	cfg := config.Get()

	// Check database connection
	dbConnected := false
	if db := database.Get(); db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Ping(); err == nil {
				dbConnected = true
			}
		}
	}

	return c.JSON(HealthResponse{
		Status:      "healthy",
		Service:     cfg.AppName,
		Version:     cfg.AppVersion,
		DBConnected: dbConnected,
	})
}
