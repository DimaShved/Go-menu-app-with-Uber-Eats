package crud

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"uber-go-menu/internal/pkg/errorx"
)

func TestResourceHandlerReturnsInvalidInputForBadJSON(t *testing.T) {
	app := testAppWithAPIErrorJSON(t)
	NewHandler(baseResource()).RegisterRoutes(app)

	response, err := app.Test(httptest.NewRequest("POST", "/test-resources/", strings.NewReader("{")))
	if err != nil {
		t.Fatalf("expected request to be handled, got %v", err)
	}
	defer response.Body.Close()

	assertAPIErrorResponse(t, response, errorx.ErrInvalidInput.HTTPStatus, "Invalid JSON body")
}

func TestResourceHandlerReturnsInvalidInputForBadID(t *testing.T) {
	app := testAppWithAPIErrorJSON(t)
	NewHandler(baseResource()).RegisterRoutes(app)

	response, err := app.Test(httptest.NewRequest("GET", "/test-resources/not-a-uuid", nil))
	if err != nil {
		t.Fatalf("expected request to be handled, got %v", err)
	}
	defer response.Body.Close()

	assertAPIErrorResponse(t, response, errorx.ErrInvalidInput.HTTPStatus, "Invalid UUID format")
}

func testAppWithAPIErrorJSON(t *testing.T) *fiber.App {
	t.Helper()

	return fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			apiErr, ok := err.(*errorx.APIError)
			if !ok {
				t.Fatalf("expected API error, got %T: %v", err, err)
			}
			return c.Status(apiErr.HTTPStatus).JSON(apiErr)
		},
	})
}

func assertAPIErrorResponse(t *testing.T, response *http.Response, wantStatus int, wantDetails string) {
	t.Helper()

	if response.StatusCode != wantStatus {
		t.Fatalf("expected status %d, got %d", wantStatus, response.StatusCode)
	}

	var body errorx.APIError
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("expected API error JSON, got %v", err)
	}

	if body.Code != errorx.ErrInvalidInput.Code {
		t.Fatalf("expected code %d, got %d", errorx.ErrInvalidInput.Code, body.Code)
	}
	if body.HTTPStatus != wantStatus {
		t.Fatalf("expected body status %d, got %d", wantStatus, body.HTTPStatus)
	}
	if body.Message != errorx.ErrInvalidInput.Message {
		t.Fatalf("expected message %q, got %q", errorx.ErrInvalidInput.Message, body.Message)
	}
	if body.Details != wantDetails {
		t.Fatalf("expected details %q, got %q", wantDetails, body.Details)
	}
}
