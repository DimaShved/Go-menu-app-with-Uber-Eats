package app

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"uber-go-menu/internal/pkg/errorx"
)

func TestAPIErrorForKeepsAPIErrorContract(t *testing.T) {
	apiErr := errorx.ErrInvalidInput.WithDetails("Invalid JSON body")

	got := apiErrorFor(apiErr)

	if got != apiErr {
		t.Fatalf("expected the original API error to be returned")
	}
}

func TestAPIErrorForMapsRecordNotFound(t *testing.T) {
	err := errorx.ErrRecordNotFound.Wrap(errors.New("menu item is gone"))

	got := apiErrorFor(err)

	if got.Code != errorx.ErrNotFound.Code {
		t.Fatalf("expected code %d, got %d", errorx.ErrNotFound.Code, got.Code)
	}
	if got.HTTPStatus != errorx.ErrNotFound.HTTPStatus {
		t.Fatalf("expected status %d, got %d", errorx.ErrNotFound.HTTPStatus, got.HTTPStatus)
	}
	if got.Message != errorx.ErrNotFound.Message {
		t.Fatalf("expected message %q, got %q", errorx.ErrNotFound.Message, got.Message)
	}
	if !strings.Contains(got.Details, "menu item is gone") {
		t.Fatalf("expected original error detail, got %q", got.Details)
	}
}

func TestAPIErrorForMapsUnexpectedErrorsToInternal(t *testing.T) {
	got := apiErrorFor(errors.New("database fell over"))

	if got.Code != errorx.ErrInternal.Code {
		t.Fatalf("expected code %d, got %d", errorx.ErrInternal.Code, got.Code)
	}
	if got.HTTPStatus != errorx.ErrInternal.HTTPStatus {
		t.Fatalf("expected status %d, got %d", errorx.ErrInternal.HTTPStatus, got.HTTPStatus)
	}
	if got.Message != errorx.ErrInternal.Message {
		t.Fatalf("expected message %q, got %q", errorx.ErrInternal.Message, got.Message)
	}
	if got.Details != "An unexpected error occurred" {
		t.Fatalf("expected generic detail, got %q", got.Details)
	}
}

func TestErrorHandlerWritesAPIErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   int
		wantMsg    string
		wantDetail string
	}{
		{
			name:       "invalid input",
			err:        errorx.ErrInvalidInput.WithDetails("Invalid JSON body"),
			wantStatus: errorx.ErrInvalidInput.HTTPStatus,
			wantCode:   errorx.ErrInvalidInput.Code,
			wantMsg:    errorx.ErrInvalidInput.Message,
			wantDetail: "Invalid JSON body",
		},
		{
			name:       "record not found",
			err:        errorx.ErrRecordNotFound.Wrap(errors.New("restaurant not found")),
			wantStatus: errorx.ErrNotFound.HTTPStatus,
			wantCode:   errorx.ErrNotFound.Code,
			wantMsg:    errorx.ErrNotFound.Message,
			wantDetail: "restaurant not found",
		},
		{
			name:       "unexpected error",
			err:        errors.New("panic without the panic"),
			wantStatus: errorx.ErrInternal.HTTPStatus,
			wantCode:   errorx.ErrInternal.Code,
			wantMsg:    errorx.ErrInternal.Message,
			wantDetail: "An unexpected error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
			app.Get("/boom", func(fiber.Ctx) error {
				return tt.err
			})

			response, err := app.Test(httptest.NewRequest("GET", "/boom", nil))
			if err != nil {
				t.Fatalf("expected request to be handled, got %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, response.StatusCode)
			}

			var body errorx.APIError
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatalf("expected API error JSON, got %v", err)
			}

			if body.Code != tt.wantCode {
				t.Fatalf("expected code %d, got %d", tt.wantCode, body.Code)
			}
			if body.HTTPStatus != tt.wantStatus {
				t.Fatalf("expected body status %d, got %d", tt.wantStatus, body.HTTPStatus)
			}
			if body.Message != tt.wantMsg {
				t.Fatalf("expected message %q, got %q", tt.wantMsg, body.Message)
			}
			if !strings.Contains(body.Details, tt.wantDetail) {
				t.Fatalf("expected details to contain %q, got %q", tt.wantDetail, body.Details)
			}
		})
	}
}
