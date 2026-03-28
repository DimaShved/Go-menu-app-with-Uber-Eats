package crud

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
