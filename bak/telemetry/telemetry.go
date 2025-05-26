package telemetry

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func DevelopmentLogger() *slog.Logger {
	return slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	)
}

func ProductionLogger() *slog.Logger {
	return slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelError,
			TimeFormat: time.Kitchen,
		}),
	)
}
