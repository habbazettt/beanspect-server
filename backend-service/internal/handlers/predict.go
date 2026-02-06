package handlers

import (
	"io"

	"github.com/beanspect/backend-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// PredictHandler handles image prediction requests
type PredictHandler struct {
	inferenceClient *services.InferenceClient
}

// NewPredictHandler creates a new predict handler
func NewPredictHandler() *PredictHandler {
	return &PredictHandler{
		inferenceClient: services.NewInferenceClient(),
	}
}

// Predict proxies prediction requests to the inference service
func (h *PredictHandler) Predict(c *fiber.Ctx) error {
	// Get file from form
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

	// Send to inference service
	prediction, err := h.inferenceClient.Predict(file.Filename, content)
	if err != nil {
		log.Error().Err(err).Msg("Inference service error")
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   true,
			"code":    "INFERENCE_ERROR",
			"message": err.Error(),
		})
	}

	return c.JSON(prediction)
}
