package main

import (
	"log/slog"
	"os"
	"time"
	"uber-go-menu/internal/app"
	"uber-go-menu/internal/config"
	"uber-go-menu/internal/pkg/db"
	"uber-go-menu/internal/pkg/validator"
)

func main() {
	logLevel := slog.LevelInfo
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339Nano))
			}
			return a
		},
	}))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	database, err := db.Connect(&cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	vld := validator.Validate()
	server := app.NewHTTPServer(database, vld)

	slog.Info("Starting server", slog.String("port", cfg.App.PORT))
	if err := server.Listen(":" + cfg.App.PORT); err != nil {
		slog.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
