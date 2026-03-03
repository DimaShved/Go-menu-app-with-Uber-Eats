package app

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"log/slog"
	"time"
	"uber-go-menu/internal/pkg/errorx"
)

func errorHandler(c fiber.Ctx, err error) error {
	slog.ErrorContext(c.UserContext(), "Unhandled error in request",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.Any("error", err),
	)

	apiErr := errorx.ErrInternal.WithDetails("An unexpected error occurred")

	var e *errorx.APIError
	if errors.As(err, &e) {
		apiErr = e
	}

	return c.Status(apiErr.HTTPStatus).JSON(fiber.Map{
		"code":    apiErr.Code,
		"message": apiErr.Message,
		"details": apiErr.Details,
	})
}

func requestLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		status := c.Response().StatusCode()
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
