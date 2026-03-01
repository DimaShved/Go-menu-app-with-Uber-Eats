package repository

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
)

type MenuAvailabilityRepo interface {
	IGenericRepository[*domain.MenuAvailability]
	UpsertAvailabilities(menuAvailability []*domain.MenuAvailability) ([]*domain.MenuAvailability, error)
}

type menuAvailabilityRepo struct {
	*GenericRepository[*domain.MenuAvailability]
}

func NewMenuAvailabilityRepo(db *gorm.DB) MenuAvailabilityRepo {
	return &menuAvailabilityRepo{NewGenericRepo[*domain.MenuAvailability](db)}
}

func (r *menuAvailabilityRepo) UpsertAvailabilities(menuAvailabilities []*domain.MenuAvailability) ([]*domain.MenuAvailability, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return nil, errorx.ErrDatabase.WithDetails(tx.Error.Error())
	}

	var valuePlaceholders []string
	var args []interface{}

	for _, ma := range menuAvailabilities {
		valuePlaceholders = append(valuePlaceholders, "(?, ?, ?, ?, NOW(), NOW())")
		args = append(args, ma.MenuSectionId, ma.DayOfWeek, ma.OpenTime, ma.CloseTime)
	}

	valuesString := strings.Join(valuePlaceholders, ", ")

	query := fmt.Sprintf(`
        INSERT INTO menu_availabilities (menu_section_id, day_of_week, open_time, close_time, created_at, updated_at)
        VALUES %s
        ON CONFLICT (menu_section_id, day_of_week) DO UPDATE SET 
        open_time = EXCLUDED.open_time, close_time = EXCLUDED.close_time, updated_at = EXCLUDED.updated_at
        RETURNING *
    `, valuesString)

	rows, err := tx.Raw(query, args...).Rows()
	if err != nil {
		tx.Rollback()
		log.Printf("Error executing upsert query: %v", err)
		return nil, errorx.ErrDatabaseQuery.WithDetails(err.Error())
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error during close rows: %v", closeErr)
		}
	}()

	var returnedAvailabilities []*domain.MenuAvailability
	for rows.Next() {
		var ma domain.MenuAvailability
		if err := tx.ScanRows(rows, &ma); err != nil {
			tx.Rollback()
			log.Printf("Error scanning rows: %v", err)
			return nil, errorx.ErrDatabaseQuery.WithDetails(err.Error())
		}
		returnedAvailabilities = append(returnedAvailabilities, &ma)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, errorx.ErrDatabase.WithDetails(err.Error())
	}

	return returnedAvailabilities, nil
}
