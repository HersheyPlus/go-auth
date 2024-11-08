package database

import (
	"fmt"
	"log"
	"time"
	"context"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/HersheyPlus/go-auth/config"
	"github.com/HersheyPlus/go-auth/models"
)

// instance
var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

// ConnectDatabase establishes connection to PostgreSQL database
func ConnectDatabase(cfg *config.Config) error {
	if err := connectWithRetries(cfg); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	if err := configureConnectionPool(cfg); err != nil {
		return fmt.Errorf("failed to configure connection pool: %w", err)
	}

	// Run migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func runMigrations() error {
	// Enable UUID extension for PostgreSQL
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return fmt.Errorf("failed to create uuid extension: %w", err)
	}

	// Run migrations in specific order
	if err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{}, // Move RefreshToken to models package
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func connectWithRetries(cfg *config.Config) error {
	var err error
	retryDelay := time.Second

	for attempt := 1; attempt <= cfg.Database.MaxRetries; attempt++ {
		dsn := cfg.GetDBConnString()
		db, err = gorm.Open(postgres.Open(dsn), cfg.GormConfig())
		if err == nil {
			log.Println("Successfully connected to database")
			break
		}

		if attempt == cfg.Database.MaxRetries {
			return fmt.Errorf("max retries (%d) reached: %w", cfg.Database.MaxRetries, err)
		}

		retryWait := retryDelay * time.Duration(attempt)
		log.Printf("Database connection attempt %d/%d failed: %v. Retrying in %v...",
			attempt, cfg.Database.MaxRetries, err, retryWait)
		time.Sleep(retryWait)
	}

	return nil
}

func configureConnectionPool(cfg *config.Config) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure pool
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Verify connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("connection verification failed: %w", err)
	}

	log.Println("Database connection pool configured successfully")
	return nil
}

func CloseDB() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("error closing database connection: %w", err)
	}

	log.Println("Database connection closed successfully")
	return nil
}