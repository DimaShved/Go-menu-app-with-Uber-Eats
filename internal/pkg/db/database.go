package db

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"time"
	"uber-go-menu/internal/config"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
)

func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DbName,
		cfg.Port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errorx.ErrDatabase.Wrap(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	models := []interface{}{
		&domain.Restaurant{},
		&domain.MenuSection{},
		&domain.MenuCategory{},
		&domain.MenuItem{},
		&domain.MenuAvailability{},
		&domain.Variation{},
		&domain.VariationOption{},
		&domain.Modifier{},
		&domain.ModifierOption{},
	}

	err = database.WithContext(ctx).AutoMigrate(models...)

	if err != nil {
		return nil, errorx.ErrDatabase.Wrap(fmt.Errorf("failed to migrate database: %w", err))
	}

	slog.Info("Database connection established and migrations completed.")
	return database, nil
}
