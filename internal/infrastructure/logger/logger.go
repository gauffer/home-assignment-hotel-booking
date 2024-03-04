package logger

import (
	"log/slog"
	"time"

	"github.com/go-chi/httplog/v2"
)

func SetupLogger() *httplog.Logger {
	logger := httplog.NewLogger("booking", httplog.Options{
		JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		Tags: map[string]string{
			"version": "0.0.1",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
	})
	return logger
}
