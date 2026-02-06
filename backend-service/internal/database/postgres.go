package database

import (
	"fmt"

	"github.com/beanspect/backend-service/internal/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	log.Info().
		Str("host", cfg.DBHost).
		Int("port", cfg.DBPort).
		Str("database", cfg.DBName).
		Msg("Connecting to PostgreSQL...")

	var gormConfig *gorm.Config
	if cfg.Env == "development" {
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info().Msg("Connected to PostgreSQL successfully")
	return db, nil
}

// Get returns the database connection
func Get() *gorm.DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
