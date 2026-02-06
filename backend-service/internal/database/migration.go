package database

import (
	"github.com/beanspect/backend-service/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Info().Msg("Running database migrations...")

	err := db.AutoMigrate(
		&models.SpeciesOrigin{},
	)
	if err != nil {
		return err
	}

	log.Info().Msg("Database migrations completed")
	return nil
}

// SeedSpeciesOrigins seeds the initial coffee species data
func SeedSpeciesOrigins(db *gorm.DB) error {
	log.Info().Msg("Seeding species origins data...")

	// Check if data already exists
	var count int64
	db.Model(&models.SpeciesOrigin{}).Count(&count)
	if count > 0 {
		log.Info().Int64("count", count).Msg("Species origins data already exists, skipping seed")
		return nil
	}

	species := []models.SpeciesOrigin{
		{
			Species:        "arabica",
			CommonName:     "Arabica Coffee",
			ScientificName: "Coffea arabica",
			Country:        "Ethiopia",
			Region:         "Kaffa Province",
			Latitude:       7.0000,
			Longitude:      36.0000,
			Description:    "Arabica coffee is considered the most superior species of coffee. It originated in the highlands of Ethiopia and is known for its smooth, complex flavor profile with notes of fruit, berries, and wine-like acidity.",
			TasteProfile:   "Sweet, soft, fruity with notes of berries, chocolate, and caramel. Complex acidity ranging from citrus to wine-like.",
			CaffeineLevel:  "Low to Medium (1.2-1.5%)",
			Altitude:       "1000-2000m",
			ImageURL:       "https://images.unsplash.com/photo-1514432324607-a09d9b4aefdd?w=800&q=80",
		},
		{
			Species:        "robusta",
			CommonName:     "Robusta Coffee",
			ScientificName: "Coffea canephora",
			Country:        "Vietnam",
			Region:         "Central Highlands",
			Latitude:       12.0000,
			Longitude:      108.0000,
			Description:    "Robusta coffee is known for its strong, bold flavor and high caffeine content. Originally from central and western sub-Saharan Africa, it is now primarily grown in Vietnam and Indonesia.",
			TasteProfile:   "Strong, bold, earthy with notes of dark chocolate, nuts, and grain. Low acidity with a heavy body.",
			CaffeineLevel:  "High (2.2-2.7%)",
			Altitude:       "200-800m",
			ImageURL:       "https://images.unsplash.com/photo-1559056199-641a0ac8b55e?w=800&q=80",
		},
		{
			Species:        "liberica",
			CommonName:     "Liberica Coffee",
			ScientificName: "Coffea liberica",
			Country:        "Philippines",
			Region:         "Batangas",
			Latitude:       13.7500,
			Longitude:      121.0000,
			Description:    "Liberica coffee has large, irregular-shaped beans with a unique aroma. Originally from Liberia, West Africa, it is now primarily grown in the Philippines and Malaysia. Known locally as 'Kapeng Barako'.",
			TasteProfile:   "Bold, smoky, woody with floral and fruity notes. Unique aroma described as jackfruit-like.",
			CaffeineLevel:  "Medium (1.2-1.5%)",
			Altitude:       "200-400m",
			ImageURL:       "https://images.unsplash.com/photo-1447933601403-0c6688de566e?w=800&q=80",
		},
		{
			Species:        "excelsa",
			CommonName:     "Excelsa Coffee",
			ScientificName: "Coffea excelsa (Coffea liberica var. dewevrei)",
			Country:        "Philippines",
			Region:         "Southeast Asia",
			Latitude:       7.5000,
			Longitude:      124.0000,
			Description:    "Excelsa coffee is a rare variety often classified as a variant of Liberica. It has a distinctive tart, fruity, and mysterious flavor profile. Primarily grown in Southeast Asia.",
			TasteProfile:   "Tart, fruity, complex with dark roast notes. Has a wine-like, popcorn, or fruity aftertaste.",
			CaffeineLevel:  "Low to Medium (1.0-1.4%)",
			Altitude:       "300-600m",
			ImageURL:       "https://images.unsplash.com/photo-1611854779393-1b2da9d400fe?w=800&q=80",
		},
	}

	for _, s := range species {
		if err := db.Create(&s).Error; err != nil {
			log.Error().Err(err).Str("species", s.Species).Msg("Failed to seed species")
			return err
		}
		log.Info().Str("species", s.Species).Msg("Seeded species origin")
	}

	log.Info().Msg("Species origins data seeded successfully")
	return nil
}
