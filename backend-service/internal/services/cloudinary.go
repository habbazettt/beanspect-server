package services

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/rs/zerolog/log"
)

// CloudinaryService handles image uploads to Cloudinary
type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryService creates a new Cloudinary service
func NewCloudinaryService() (*CloudinaryService, error) {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		return nil, fmt.Errorf("CLOUDINARY_URL environment variable is not set")
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	log.Info().Msg("Cloudinary service initialized")
	return &CloudinaryService{cld: cld}, nil
}

// UploadImage uploads an image to Cloudinary
func (s *CloudinaryService) UploadImage(ctx context.Context, file io.Reader, filename string, folder string) (*UploadResult, error) {
	unique := true
	overwrite := false

	uploadParams := uploader.UploadParams{
		Folder:         folder,
		PublicID:       filename,
		UniqueFilename: &unique,
		Overwrite:      &overwrite,
	}

	result, err := s.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("Failed to upload image to Cloudinary")
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	log.Info().
		Str("public_id", result.PublicID).
		Str("url", result.SecureURL).
		Msg("Image uploaded to Cloudinary")

	return &UploadResult{
		PublicID: result.PublicID,
		URL:      result.SecureURL,
		Format:   result.Format,
		Width:    result.Width,
		Height:   result.Height,
		Bytes:    result.Bytes,
	}, nil
}

// UploadResult represents the result of an upload
type UploadResult struct {
	PublicID string `json:"public_id"`
	URL      string `json:"url"`
	Format   string `json:"format"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Bytes    int    `json:"bytes"`
}

// GetOptimizedURL returns an optimized URL for an image
func (s *CloudinaryService) GetOptimizedURL(publicID string, width, height int) string {
	asset, err := s.cld.Image(publicID)
	if err != nil {
		log.Error().Err(err).Str("public_id", publicID).Msg("Failed to get image asset")
		return ""
	}

	url, err := asset.String()
	if err != nil {
		log.Error().Err(err).Str("public_id", publicID).Msg("Failed to generate URL")
		return ""
	}

	// Append transformation parameters
	return fmt.Sprintf("%s?w=%d&h=%d&c=fill&q=auto&f=auto", url, width, height)
}

// DeleteImage deletes an image from Cloudinary
func (s *CloudinaryService) DeleteImage(ctx context.Context, publicID string) error {
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	log.Info().Str("public_id", publicID).Msg("Image deleted from Cloudinary")
	return nil
}
