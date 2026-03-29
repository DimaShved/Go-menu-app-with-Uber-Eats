package menu_availability

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
)

func TestBuildAvailabilitiesCreatesBatchRows(t *testing.T) {
	sectionID := uuid.New()
	request := BatchUpsertRequest{
		MenuSectionID: sectionID,
		Availabilities: []AvailabilityRequest{
			{DayOfWeek: 1, OpenTime: 480, CloseTime: 900},
			{DayOfWeek: 2, OpenTime: 510, CloseTime: 930},
		},
	}

	availabilities, err := buildAvailabilities(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := []domain.MenuAvailability{
		{MenuSectionId: sectionID, DayOfWeek: 1, OpenTime: 480, CloseTime: 900},
		{MenuSectionId: sectionID, DayOfWeek: 2, OpenTime: 510, CloseTime: 930},
	}
	if !sameAvailabilities(availabilities, want) {
		t.Fatalf("expected availabilities %+v, got %+v", want, availabilities)
	}
}

func TestBuildAvailabilitiesRejectsDuplicateDays(t *testing.T) {
	sectionID := uuid.New()
	request := BatchUpsertRequest{
		MenuSectionID: sectionID,
		Availabilities: []AvailabilityRequest{
			{DayOfWeek: 4, OpenTime: 480, CloseTime: 900},
			{DayOfWeek: 4, OpenTime: 540, CloseTime: 960},
		},
	}

	availabilities, err := buildAvailabilities(request)
	if err == nil {
		t.Fatal("expected duplicate validation error")
	}
	if availabilities != nil {
		t.Fatalf("expected no availabilities, got %+v", availabilities)
	}

	apiErr, ok := err.(*errorx.APIError)
	if !ok || apiErr.Code != errorx.ErrInvalidInput.Code {
		t.Fatalf("expected invalid input error, got %T: %v", err, err)
	}
	if !strings.Contains(apiErr.Details, "duplicate availability") {
		t.Fatalf("expected duplicate detail, got %q", apiErr.Details)
	}
	if !strings.Contains(apiErr.Details, sectionID.String()) {
		t.Fatalf("expected menu section ID in details, got %q", apiErr.Details)
	}
}

func sameAvailabilities(got []domain.MenuAvailability, want []domain.MenuAvailability) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i].MenuSectionId != want[i].MenuSectionId ||
			got[i].DayOfWeek != want[i].DayOfWeek ||
			got[i].OpenTime != want[i].OpenTime ||
			got[i].CloseTime != want[i].CloseTime {
			return false
		}
	}
	return true
}
