package main

import (
	"context"
	"log/slog"

	slogotel "github.com/samber/slog-otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
)

func main() {
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
	)
	tracer := tp.Tracer("hello/world")

	ctx, span := tracer.Start(context.Background(), "foo")
	defer span.End()

	span.AddEvent("bar")

	// setup loki client
	config, _ := loki.NewDefaultConfig("https://monitoring.mbv-labs.com/loki/api/v1/push")
	config.TenantID = "local-test"
	client, _ := loki.New(config)

	logger := slog.New(slogloki.Option{
		Level:  slog.LevelDebug,
		Client: client,
		AttrFromContext: []func(ctx context.Context) []slog.Attr{
			slogotel.ExtractOtelAttrFromContext([]string{"tracing"}, "trace_id", "span_id"),
		},
	}.NewLokiHandler())
	logger = logger.
		With("environment", "dev").
		With("release", "v1.0.0").
		With("service_name", "yoyoyoy").
		With("container", "yoyoyoy")

	// log error
	logger.Error("caramba!")

	// log user signup
	logger.Info("user registration")

	logger.ErrorContext(ctx, "a message")

	// stop loki client and purge buffers
	client.Stop()
}
