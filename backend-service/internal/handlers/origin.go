package handlers

import (
	"github.com/beanspect/backend-service/internal/database"
	"github.com/beanspect/backend-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// OriginHandler handles species origin requests
type OriginHandler struct{}

// NewOriginHandler creates a new origin handler
func NewOriginHandler() *OriginHandler {
	return &OriginHandler{}
}

// GetAllOrigins returns all species origins
func (h *OriginHandler) GetAllOrigins(c *fiber.Ctx) error {
	db := database.Get()
	if db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   true,
			"code":    "DB_NOT_CONNECTED",
			"message": "Database connection not available",
		})
	}

	var origins []models.SpeciesOrigin
	if err := db.Find(&origins).Error; err != nil {
		log.Error().Err(err).Msg("Failed to fetch species origins")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"code":    "FETCH_ERROR",
			"message": "Failed to fetch species origins",
		})
	}

	return c.JSON(fiber.Map{
		"data":  origins,
		"count": len(origins),
	})
}

// GetOriginBySpecies returns origin data for a specific species
func (h *OriginHandler) GetOriginBySpecies(c *fiber.Ctx) error {
	species := c.Params("species")
	if species == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"code":    "SPECIES_REQUIRED",
			"message": "Species parameter is required",
		})
	}

	db := database.Get()
	if db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   true,
			"code":    "DB_NOT_CONNECTED",
			"message": "Database connection not available",
		})
	}

	var origin models.SpeciesOrigin
	result := db.Where("species = ?", species).First(&origin)
	if result.Error != nil {
		log.Warn().Str("species", species).Msg("Species not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"code":    "SPECIES_NOT_FOUND",
			"message": "Species '" + species + "' not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": origin,
	})
}

// GetOriginGeoJSON returns origin data in GeoJSON format for mapping
func (h *OriginHandler) GetOriginGeoJSON(c *fiber.Ctx) error {
	db := database.Get()
	if db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   true,
			"code":    "DB_NOT_CONNECTED",
			"message": "Database connection not available",
		})
	}

	var origins []models.SpeciesOrigin
	if err := db.Find(&origins).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"code":    "FETCH_ERROR",
			"message": "Failed to fetch species origins",
		})
	}

	// Build GeoJSON FeatureCollection
	features := make([]fiber.Map, len(origins))
	for i, origin := range origins {
		features[i] = fiber.Map{
			"type": "Feature",
			"geometry": fiber.Map{
				"type":        "Point",
				"coordinates": []float64{origin.Longitude, origin.Latitude},
			},
			"properties": fiber.Map{
				"species":        origin.Species,
				"common_name":    origin.CommonName,
				"country":        origin.Country,
				"region":         origin.Region,
				"description":    origin.Description,
				"taste_profile":  origin.TasteProfile,
				"caffeine_level": origin.CaffeineLevel,
				"image_url":      origin.ImageURL,
			},
		}
	}

	return c.JSON(fiber.Map{
		"type":     "FeatureCollection",
		"features": features,
	})
}
