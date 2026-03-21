package app

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"uber-go-menu/internal/pkg/errorx"
)

func errorHandler(c fiber.Ctx, err error) error {
	apiErr := apiErrorFor(err)

	slog.ErrorContext(c.UserContext(), "Unhandled error in request",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.Int("status", apiErr.HTTPStatus),
		slog.Any("error", err),
	)

	return c.Status(apiErr.HTTPStatus).JSON(apiErr)
}

func apiErrorFor(err error) *errorx.APIError {
	var e *errorx.APIError
	if errors.As(err, &e) {
		return e
	}

	if errorx.IsAppError(err, errorx.ErrRecordNotFound) {
		return errorx.ErrNotFound.WithDetails(err.Error())
	}

	return errorx.ErrInternal.WithDetails("An unexpected error occurred")
}

func requestLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		status := c.Response().StatusCode()
		if err != nil {
			status = apiErrorFor(err).HTTPStatus
		}
		logEvent := slog.InfoContext

		if status >= 400 {
			logEvent = slog.WarnContext
		}
		if status >= 500 {
			logEvent = slog.ErrorContext
		}

		logEvent(c.UserContext(), "Request handled",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.Duration("latency", latency),
		)

		return err
	}
}
