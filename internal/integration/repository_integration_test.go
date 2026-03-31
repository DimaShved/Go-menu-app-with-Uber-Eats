package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	container, err := tcpostgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		tcpostgres.WithDatabase("menu_test"),
		tcpostgres.WithUsername("menu_test"),
		tcpostgres.WithPassword("menu_test"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic(err)
	}

	code := 1
	defer func() {
		if testDB != nil {
			if sqlDB, err := testDB.DB(); err == nil {
				_ = sqlDB.Close()
			}
		}
		_ = container.Terminate(context.Background())
		os.Exit(code)
	}()

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := testDB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
		panic(err)
	}
	if err := migrateTestSchema(testDB); err != nil {
		panic(err)
	}

	code = m.Run()
}

func migrateTestSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Restaurant{},
		&domain.MenuSection{},
		&domain.MenuCategory{},
		&domain.MenuItem{},
		&domain.MenuAvailability{},
		&domain.Variation{},
		&domain.VariationOption{},
		&domain.Modifier{},
		&domain.ModifierOption{},
	)
}
