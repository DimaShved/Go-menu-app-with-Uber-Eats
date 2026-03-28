package crud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestResourcePrepareForHandlerValidatesRequiredFields(t *testing.T) {
	tests := []struct {
		name         string
		change       func(*testResource)
		missingField string
	}{
		{
			name: "missing name",
			change: func(resource *testResource) {
				resource.Name = ""
			},
			missingField: "Name",
		},
		{
			name: "missing path",
			change: func(resource *testResource) {
				resource.Path = ""
			},
			missingField: "Path",
		},
		{
			name: "missing repository",
			change: func(resource *testResource) {
				resource.Repository = nil
			},
			missingField: "Repository",
		},
		{
			name: "missing tx manager",
			change: func(resource *testResource) {
				resource.TxManager = nil
			},
			missingField: "TxManager",
		},
		{
			name: "missing create mapper",
			change: func(resource *testResource) {
				resource.MapCreate = nil
			},
			missingField: "MapCreate",
		},
		{
			name: "missing update applier",
			change: func(resource *testResource) {
				resource.ApplyUpdate = nil
			},
			missingField: "ApplyUpdate",
		},
		{
			name: "missing response mapper",
			change: func(resource *testResource) {
				resource.MapResponse = nil
			},
			missingField: "MapResponse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := baseResource()
			tt.change(&resource)

			requirePanicContains(t, func() {
				_ = resource.prepareForHandler()
			}, tt.missingField)
		})
	}
}

func TestResourcePrepareForServiceValidatesRequiredFields(t *testing.T) {
	tests := []struct {
		name         string
		change       func(*testResource)
		missingField string
	}{
		{
			name: "missing name",
			change: func(resource *testResource) {
				resource.Name = ""
			},
			missingField: "Name",
		},
		{
			name: "missing repository",
			change: func(resource *testResource) {
				resource.Repository = nil
			},
			missingField: "Repository",
		},
		{
			name: "missing tx manager",
			change: func(resource *testResource) {
				resource.TxManager = nil
			},
			missingField: "TxManager",
		},
		{
			name: "missing create mapper",
			change: func(resource *testResource) {
				resource.MapCreate = nil
			},
			missingField: "MapCreate",
		},
		{
			name: "missing update applier",
			change: func(resource *testResource) {
				resource.ApplyUpdate = nil
			},
			missingField: "ApplyUpdate",
		},
		{
			name: "missing response mapper",
			change: func(resource *testResource) {
				resource.MapResponse = nil
			},
			missingField: "MapResponse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := baseResource()
			tt.change(&resource)

			requirePanicContains(t, func() {
				_ = resource.prepareForService()
			}, tt.missingField)
		})
	}
}

func TestResourcePrepareForServiceDoesNotRequirePath(t *testing.T) {
	resource := baseResource()
	resource.Path = ""

	requireNotPanics(t, func() {
		prepared := resource.prepareForService()
		if prepared.Path != "" {
			t.Fatalf("expected Path to remain empty, got %q", prepared.Path)
		}
	})
}

func TestResourcePrepareUsesDefaultsForOptionalFields(t *testing.T) {
	tests := []struct {
		name          string
		prepare       func(testResource) testResource
		useTypedHooks bool
	}{
		{
			name: "handler",
			prepare: func(resource testResource) testResource {
				return resource.prepareForHandler()
			},
		},
		{
			name: "service",
			prepare: func(resource testResource) testResource {
				return resource.prepareForService()
			},
		},
		{
			name: "service with typed nil hooks",
			prepare: func(resource testResource) testResource {
				return resource.prepareForService()
			},
			useTypedHooks: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := baseResource()
			resource.Hooks = nil
			resource.Timeout = 0

			if tt.useTypedHooks {
				resource.Hooks = (*testHooks)(nil)
			}

			prepared := tt.prepare(resource)

			if prepared.Timeout != DefaultTimeout {
				t.Fatalf("expected Timeout %s, got %s", DefaultTimeout, prepared.Timeout)
			}
			if _, ok := prepared.Hooks.(NoopHooks[testEntity, testCreateRequest, testUpdateRequest, testResponse]); !ok {
				t.Fatalf("expected Hooks to default to NoopHooks, got %T", prepared.Hooks)
			}
		})
	}
}

func TestResourcePrepareTreatsTypedNilDependenciesAsMissing(t *testing.T) {
	t.Run("repository", func(t *testing.T) {
		var repository *testRepository
		resource := baseResource()
		resource.Repository = repository

		requirePanicContains(t, func() {
			_ = resource.prepareForService()
		}, "Repository")
	})

	t.Run("tx manager", func(t *testing.T) {
		var txManager *testTxManager
		resource := baseResource()
		resource.TxManager = txManager

		requirePanicContains(t, func() {
			_ = resource.prepareForService()
		}, "TxManager")
	})
}

type testEntity struct {
	id        uuid.UUID
	deletedAt *time.Time
}

func (e testEntity) GetID() uuid.UUID {
	return e.id
}

func (e testEntity) GetDeletedAt() *time.Time {
	return e.deletedAt
}

type testCreateRequest struct{}
type testUpdateRequest struct{}
type testResponse struct{}

type testResource = Resource[testEntity, testCreateRequest, testUpdateRequest, testResponse]

type testHooks struct {
	NoopHooks[testEntity, testCreateRequest, testUpdateRequest, testResponse]
}

func baseResource() testResource {
	return testResource{
		Name:       "test resources",
		Path:       "/test-resources",
		Repository: &testRepository{},
		TxManager:  &testTxManager{},
		MapCreate: func(testCreateRequest) (*testEntity, error) {
			return &testEntity{}, nil
		},
		ApplyUpdate: func(*testEntity, testUpdateRequest) error {
			return nil
		},
		MapResponse: func(*testEntity) (testResponse, error) {
			return testResponse{}, nil
		},
	}
}

type testRepository struct{}

var _ Repository[testEntity] = (*testRepository)(nil)

func (testRepository) Create(context.Context, *gorm.DB, *testEntity) error {
	return nil
}

func (testRepository) Update(context.Context, *gorm.DB, *testEntity) error {
	return nil
}

func (testRepository) Delete(context.Context, *gorm.DB, uuid.UUID) error {
	return nil
}

func (testRepository) GetByID(context.Context, *gorm.DB, uuid.UUID) (*testEntity, error) {
	return &testEntity{}, nil
}

func (testRepository) List(context.Context, *gorm.DB) ([]testEntity, error) {
	return nil, nil
}

type testTxManager struct{}

var _ TxManager = (*testTxManager)(nil)

func (testTxManager) DB() *gorm.DB {
	return nil
}

func (testTxManager) WithinTx(context.Context, func(*gorm.DB) error) error {
	return nil
}

func requirePanicContains(t *testing.T, fn func(), want string) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("expected panic containing %q", want)
		}
		got := fmt.Sprint(recovered)
		if !strings.Contains(got, want) {
			t.Fatalf("expected panic containing %q, got %q", want, got)
		}
	}()

	fn()
}

func requireNotPanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("expected no panic, got %v", recovered)
		}
	}()

	fn()
}
