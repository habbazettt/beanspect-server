package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/beanspect/backend-service/internal/config"
	"github.com/beanspect/backend-service/internal/database"
	"github.com/beanspect/backend-service/internal/handlers"
	"github.com/beanspect/backend-service/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Load configuration
	cfg := config.Load()
	log.Info().
		Str("app", cfg.AppName).
		Str("version", cfg.AppVersion).
		Str("env", cfg.Env).
		Msg("Starting BeanSpect Backend Service")

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		// Continue without database for now - will retry on requests
	} else {
		defer database.Close()
		// Run migrations
		if err := database.Migrate(db); err != nil {
			log.Error().Err(err).Msg("Failed to run migrations")
		}
		// Seed initial data
		if err := database.SeedSpeciesOrigins(db); err != nil {
			log.Error().Err(err).Msg("Failed to seed species origins")
		}
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: errorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(cors.New(cors.Config{
		AllowOrigins: joinOrigins(cfg.CORSOrigins),
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Routes
	setupRoutes(app)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Info().Msg("Shutting down gracefully...")
		if err := app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("Error during shutdown")
		}
	}()

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Info().Str("address", addr).Msg("Server starting")
	if err := app.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func setupRoutes(app *fiber.App) {
	// Root
	app.Get("/", func(c *fiber.Ctx) error {
		cfg := config.Get()
		return c.JSON(fiber.Map{
			"service": cfg.AppName,
			"version": cfg.AppVersion,
			"docs":    "/api/docs",
			"health":  "/health",
		})
	})

	// Health check
	app.Get("/health", handlers.Health)

	// API routes
	api := app.Group("/api")

	// Predict handler
	predictHandler := handlers.NewPredictHandler()
	api.Post("/predict", predictHandler.Predict)

	// Origin handler
	originHandler := handlers.NewOriginHandler()
	api.Get("/origins", originHandler.GetAllOrigins)
	api.Get("/origins/geojson", originHandler.GetOriginGeoJSON)
	api.Get("/origin/:species", originHandler.GetOriginBySpecies)

	// Analyze handler
	analyzeHandler := handlers.NewAnalyzeHandler()
	api.Post("/analyze", analyzeHandler.Analyze)
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	log.Error().Err(err).Int("status", code).Msg("Request error")

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"code":    "INTERNAL_ERROR",
		"message": err.Error(),
	})
}

func joinOrigins(origins []string) string {
	result := ""
	for i, origin := range origins {
		if i > 0 {
			result += ","
		}
		result += origin
	}
	return result
}
