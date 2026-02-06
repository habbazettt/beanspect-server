package models

import (
	"time"

	"gorm.io/gorm"
)

// SpeciesOrigin represents coffee species origin data with GIS information
type SpeciesOrigin struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Species        string `gorm:"uniqueIndex;size:50;not null" json:"species"`
	CommonName     string `gorm:"size:100" json:"common_name"`
	ScientificName string `gorm:"size:150" json:"scientific_name"`

	// Origin Location
	Country   string  `gorm:"size:100;not null" json:"country"`
	Region    string  `gorm:"size:100" json:"region"`
	Latitude  float64 `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude float64 `gorm:"type:decimal(10,7)" json:"longitude"`

	// Characteristics
	Description   string `gorm:"type:text" json:"description"`
	TasteProfile  string `gorm:"type:text" json:"taste_profile"`
	CaffeineLevel string `gorm:"size:50" json:"caffeine_level"` // low, medium, high
	Altitude      string `gorm:"size:50" json:"altitude"`       // e.g., "1000-2000m"

	// Media
	ImageURL string `gorm:"size:500" json:"image_url"`

	// Metadata
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (SpeciesOrigin) TableName() string {
	return "species_origins"
}
