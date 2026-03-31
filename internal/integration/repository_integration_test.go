package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
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

func TestGormRepositoryPersistsAndSoftDeletesRestaurants(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	repository := crud.NewGormRepository[domain.Restaurant]()

	restaurant := &domain.Restaurant{
		Name:    "Best Burger",
		Address: "Main Street 1",
	}
	if err := repository.Create(ctx, testDB, restaurant); err != nil {
		t.Fatalf("create restaurant: %v", err)
	}
	if restaurant.ID == uuid.Nil {
		t.Fatal("expected database-generated ID")
	}

	restaurant.Name = "Best Burger Updated"
	if err := repository.Update(ctx, testDB, restaurant); err != nil {
		t.Fatalf("update restaurant: %v", err)
	}

	found, err := repository.GetByID(ctx, testDB, restaurant.ID)
	if err != nil {
		t.Fatalf("get updated restaurant: %v", err)
	}
	if found.Name != "Best Burger Updated" {
		t.Fatalf("expected updated name, got %q", found.Name)
	}

	list, err := repository.List(ctx, testDB)
	if err != nil {
		t.Fatalf("list restaurants: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected one restaurant before delete, got %d", len(list))
	}

	if err := repository.Delete(ctx, testDB, restaurant.ID); err != nil {
		t.Fatalf("soft delete restaurant: %v", err)
	}

	if _, err := repository.GetByID(ctx, testDB, restaurant.ID); !errorx.IsAppError(err, errorx.ErrRecordNotFound) {
		t.Fatalf("expected deleted restaurant to be filtered, got %v", err)
	}

	list, err = repository.List(ctx, testDB)
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected deleted restaurant to be omitted from list, got %d rows", len(list))
	}

	var deleted domain.Restaurant
	if err := testDB.Unscoped().First(&deleted, "id = ?", restaurant.ID).Error; err != nil {
		t.Fatalf("expected row to remain after soft delete: %v", err)
	}
	if deleted.DeletedAt == nil {
		t.Fatal("expected deleted_at to be set")
	}
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

func cleanDatabase(t *testing.T) {
	t.Helper()

	err := testDB.Exec(`
		TRUNCATE TABLE
			item_categories,
			modifier_options,
			modifiers,
			variation_options,
			variations,
			menu_availabilities,
			menu_items,
			menu_categories,
			menu_sections,
			restaurants
		RESTART IDENTITY CASCADE
	`).Error
	if err != nil {
		t.Fatalf("clean database: %v", err)
	}
}
