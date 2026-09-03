package utils

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {
	handler := slog.NewJSONHandler(
		os.Stdout,
		nil,
	)

	logger := slog.New(handler)

	return logger.With(
		slog.String("service", "task-api"),
		slog.String("environment", "development"),
	)
}
