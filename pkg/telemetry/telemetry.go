package telemetry

import (
	"log/slog"
	"os"
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/grafana/loki-client-go/loki"
	"github.com/lmittmann/tint"
	slogloki "github.com/samber/slog-loki/v3"
)

func NewTelemetry(cfg config.Cfg, release, service string) {
	switch cfg.App.Environment {
	case config.PROD_ENVIRONMENT:
		logger := productionLogger(
			cfg.Telemetry.SinkURL,
			cfg.Telemetry.TenantID,
			cfg.App.Environment,
			release,
			service,
		)
		slog.SetDefault(logger)
	case config.DEV_ENVIRONMENT:
		logger := developmentLogger()
		slog.SetDefault(logger)
	default:
		logger := developmentLogger()
		slog.SetDefault(logger)
	}
}

func productionLogger(url, tenantID, release, _, service string) *slog.Logger {
	config, _ := loki.NewDefaultConfig(url)
	config.TenantID = tenantID
	client, err := loki.New(config)
	if err != nil {
		panic(err)
	}

	defer client.Stop()

	logger := slog.New(
		slogloki.Option{Level: slog.LevelInfo, Client: client, AddSource: true}.NewLokiHandler(),
	)
	logger = logger.
		With("release", release).
		With("service", service)

	return logger
}

func developmentLogger() *slog.Logger {
	return slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			AddSource:  true,
		}),
	)
}
