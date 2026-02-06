package handlers

import (
	"io"

	"github.com/beanspect/backend-service/internal/database"
	"github.com/beanspect/backend-service/internal/models"
	"github.com/beanspect/backend-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// AnalyzeHandler handles the combined analyze requests
type AnalyzeHandler struct {
	inferenceClient *services.InferenceClient
}

// NewAnalyzeHandler creates a new analyze handler
func NewAnalyzeHandler() *AnalyzeHandler {
	return &AnalyzeHandler{
		inferenceClient: services.NewInferenceClient(),
	}
}

// AnalyzeResponse represents the combined response
type AnalyzeResponse struct {
	Prediction PredictionData `json:"prediction"`
	Origin     *OriginData    `json:"origin"`
}

// PredictionData contains prediction results
type PredictionData struct {
	Species        string                     `json:"species"`
	Confidence     float64                    `json:"confidence"`
	AllPredictions []services.ClassPrediction `json:"all_predictions"`
}

// OriginData contains species origin information
type OriginData struct {
	ID             uint    `json:"id"`
	Species        string  `json:"species"`
	CommonName     string  `json:"common_name"`
	ScientificName string  `json:"scientific_name"`
	Country        string  `json:"country"`
	Region         string  `json:"region"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Description    string  `json:"description"`
	TasteProfile   string  `json:"taste_profile"`
	CaffeineLevel  string  `json:"caffeine_level"`
	Altitude       string  `json:"altitude"`
	ImageURL       string  `json:"image_url"`
}

// Analyze receives an image, gets prediction, and returns combined data with origin
func (h *AnalyzeHandler) Analyze(c *fiber.Ctx) error {
	// Step 1: Receive image from frontend
	file, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get file from form")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"code":    "FILE_REQUIRED",
			"message": "Image file is required",
		})
	}

	// Open file
	f, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msg("Failed to open file")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"code":    "FILE_OPEN_ERROR",
			"message": "Failed to open uploaded file",
		})
	}
	defer f.Close()

	// Read file content
	content, err := io.ReadAll(f)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read file")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"code":    "FILE_READ_ERROR",
			"message": "Failed to read uploaded file",
		})
	}

	log.Info().Str("filename", file.Filename).Int("size", len(content)).Msg("Received image for analysis")

	// Step 2 & 3: Forward to inference service and receive prediction
	prediction, err := h.inferenceClient.Predict(file.Filename, content)
	if err != nil {
		log.Error().Err(err).Msg("Inference service error")
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   true,
			"code":    "INFERENCE_ERROR",
			"message": err.Error(),
		})
	}

	log.Info().
		Str("species", prediction.PredictedClass).
		Float64("confidence", prediction.Confidence).
		Msg("Received prediction from inference service")

	// Step 4: Fetch GIS origin data
	db := database.Get()
	var origin *OriginData

	if db != nil {
		var speciesOrigin models.SpeciesOrigin
		result := db.Where("species = ?", prediction.PredictedClass).First(&speciesOrigin)
		if result.Error == nil {
			origin = &OriginData{
				ID:             speciesOrigin.ID,
				Species:        speciesOrigin.Species,
				CommonName:     speciesOrigin.CommonName,
				ScientificName: speciesOrigin.ScientificName,
				Country:        speciesOrigin.Country,
				Region:         speciesOrigin.Region,
				Latitude:       speciesOrigin.Latitude,
				Longitude:      speciesOrigin.Longitude,
				Description:    speciesOrigin.Description,
				TasteProfile:   speciesOrigin.TasteProfile,
				CaffeineLevel:  speciesOrigin.CaffeineLevel,
				Altitude:       speciesOrigin.Altitude,
				ImageURL:       speciesOrigin.ImageURL,
			}
			log.Info().Str("species", origin.Species).Str("country", origin.Country).Msg("Fetched origin data")
		} else {
			log.Warn().Str("species", prediction.PredictedClass).Msg("Origin data not found for species")
		}
	} else {
		log.Warn().Msg("Database not connected, skipping origin data fetch")
	}

	// Step 5: Return combined response
	response := AnalyzeResponse{
		Prediction: PredictionData{
			Species:        prediction.PredictedClass,
			Confidence:     prediction.Confidence,
			AllPredictions: prediction.AllPredictions,
		},
		Origin: origin,
	}

	log.Info().Msg("Analysis complete, returning combined response")
	return c.JSON(response)
}
