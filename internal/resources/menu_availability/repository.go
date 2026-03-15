package menu_availability

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

type Repository struct {
	*crud.GormRepository[domain.MenuAvailability]
}

func NewRepository() *Repository {
	return &Repository{GormRepository: crud.NewGormRepository[domain.MenuAvailability]()}
}

func (r *Repository) UpsertBatch(ctx context.Context, tx *gorm.DB, availabilities []domain.MenuAvailability) ([]domain.MenuAvailability, error) {
	valuePlaceholders := make([]string, 0, len(availabilities))
	args := make([]any, 0, len(availabilities)*4)

	for _, availability := range availabilities {
		valuePlaceholders = append(valuePlaceholders, "(?, ?, ?, ?, NOW(), NOW())")
		args = append(args, availability.MenuSectionId, availability.DayOfWeek, availability.OpenTime, availability.CloseTime)
	}

	query := fmt.Sprintf(`
		INSERT INTO menu_availabilities (menu_section_id, day_of_week, open_time, close_time, created_at, updated_at)
		VALUES %s
		ON CONFLICT (menu_section_id, day_of_week) DO UPDATE SET
			open_time = EXCLUDED.open_time,
			close_time = EXCLUDED.close_time,
			updated_at = EXCLUDED.updated_at,
			deleted_at = NULL
		RETURNING *
	`, strings.Join(valuePlaceholders, ", "))

	rows, err := tx.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, errorx.ErrDatabaseQuery.WithDetails(err.Error())
	}
	defer rows.Close()

	result := make([]domain.MenuAvailability, 0, len(availabilities))
	for rows.Next() {
		var availability domain.MenuAvailability
		if err := tx.ScanRows(rows, &availability); err != nil {
			return nil, errorx.ErrDatabaseQuery.WithDetails(err.Error())
		}
		result = append(result, availability)
	}
	if err := rows.Err(); err != nil {
		return nil, errorx.ErrDatabaseQuery.WithDetails(err.Error())
	}
	return result, nil
}
