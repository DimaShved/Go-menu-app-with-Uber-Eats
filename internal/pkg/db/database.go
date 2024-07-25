package db

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"time"
	"uber-go-menu-copy/internal/config"
	"uber-go-menu-copy/internal/domain"
)

var DB *gorm.DB

func Connect(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DbName,
		cfg.Port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = database.WithContext(ctx).AutoMigrate(
		&domain.Restaurant{},
		&domain.MenuSections{},
		&domain.MenuCategory{},
		&domain.MenuItem{},
	)

	if err != nil {
		slog.Error("Failed to migrate database", "error", err)
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = database
	slog.Info("Connection to DB established.")
	return nil
}
