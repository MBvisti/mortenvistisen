package telemetry

import (
	"log/slog"
	"os"
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/grafana/loki-client-go/loki"
	"github.com/lmittmann/tint"
)

func NewTelemetry(cfg config.Cfg, release, projectName string) *loki.Client {
	switch cfg.App.Environment {
	case config.PROD_ENVIRONMENT:
		logger := developmentLogger()
		slog.SetDefault(logger)
		return nil
	case config.DEV_ENVIRONMENT:
		logger := developmentLogger()
		slog.SetDefault(logger)
		return nil
	default:
		logger := developmentLogger()
		slog.SetDefault(logger)

		return nil
	}
}

// func productionLogger(url, tenantID, release, projectName string) (*slog.Logger, *loki.Client) {
// 	cfg, _ := loki.NewDefaultConfig(url)
// 	cfg.TenantID = tenantID
// 	client, err := loki.New(cfg)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	logger := slog.New(
// 		slogloki.Option{
// 			Level:  slog.LevelInfo,
// 			Client: client,
// 			AttrFromContext: []func(ctx context.Context) []slog.Attr{
// 				slogotel.ExtractOtelAttrFromContext(
// 					[]string{"parent"},
// 					"trace_id",
// 					"span_id",
// 				),
// 			},
// 		}.NewLokiHandler(),
// 	)
// 	logger = logger.
// 		With(
// 			"release",
// 			release,
// 		).With("env", config.PROD_ENVIRONMENT).With("service_name", projectName)
//
// 	return logger, client
// }

func developmentLogger() *slog.Logger {
	return slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	)
}
