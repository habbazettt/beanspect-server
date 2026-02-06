package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"github.com/beanspect/backend-service/internal/config"
	"github.com/rs/zerolog/log"
)

// ClassPrediction represents a single prediction class
type ClassPrediction struct {
	Class      string  `json:"class"`
	Confidence float64 `json:"confidence"`
}

// PredictionResponse represents the response from inference service
type PredictionResponse struct {
	PredictedClass string            `json:"predicted_class"`
	Confidence     float64           `json:"confidence"`
	AllPredictions []ClassPrediction `json:"all_predictions"`
}

// ErrorResponse represents an error from inference service
type ErrorResponse struct {
	Error   bool   `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// InferenceClient handles communication with the inference service
type InferenceClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewInferenceClient creates a new inference service client
func NewInferenceClient() *InferenceClient {
	cfg := config.Get()
	return &InferenceClient{
		baseURL: cfg.InferenceServiceURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Predict sends an image to the inference service for classification
func (c *InferenceClient) Predict(filename string, fileContent []byte) (*PredictionResponse, error) {
	url := fmt.Sprintf("%s/predict", c.baseURL)

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Determine content type based on file extension
	contentType := "image/jpeg"
	if len(filename) > 4 {
		ext := filename[len(filename)-4:]
		if ext == ".png" {
			contentType = "image/png"
		}
	}

	// Create form file with proper content-type header
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := part.Write(fileContent); err != nil {
		return nil, fmt.Errorf("failed to write file content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Info().
		Str("url", url).
		Str("filename", filename).
		Msg("Sending prediction request to inference service")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("response_body", string(respBody)).
			Msg("Inference service error response")
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error {
			return nil, fmt.Errorf("%s: %s", errResp.Code, errResp.Message)
		}
		return nil, fmt.Errorf("inference service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse prediction response
	var prediction PredictionResponse
	if err := json.Unmarshal(respBody, &prediction); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Info().
		Str("predicted_class", prediction.PredictedClass).
		Float64("confidence", prediction.Confidence).
		Msg("Received prediction from inference service")

	return &prediction, nil
}

// HealthCheck checks if the inference service is healthy
func (c *InferenceClient) HealthCheck() (bool, error) {
	url := fmt.Sprintf("%s/health", c.baseURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
