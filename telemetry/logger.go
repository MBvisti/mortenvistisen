package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

type LogExporter interface {
	Name() string
	GetSlogHandler(ctx context.Context, res *resource.Resource, cfg *telemetryOptions) (slog.Handler, error)
	Shutdown(ctx context.Context) error
}

type traceLogHandler struct {
	handler slog.Handler
}

func (h *traceLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *traceLogHandler) Handle(ctx context.Context, record slog.Record) error {
	traceAttrs := traceAttrsFromContext(ctx)
	for _, attr := range traceAttrs {
		record.AddAttrs(attr)
	}

	return h.handler.Handle(ctx, record)
}

func (h *traceLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceLogHandler{
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h *traceLogHandler) WithGroup(name string) slog.Handler {
	return &traceLogHandler{
		handler: h.handler.WithGroup(name),
	}
}

func traceAttrsFromContext(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		attrs = append(attrs,
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	}

	return attrs
}

// OtlpHttpLogExporter exports structured slog records using OTLP over HTTP.
type OtlpHttpLogExporter struct {
	endpoint       string
	headers        map[string]string
	loggerProvider *sdklog.LoggerProvider
}

func NewOtlpLogExporter(endpoint string, headers map[string]string) *OtlpHttpLogExporter {
	if endpoint == "" {
		return nil
	}

	return &OtlpHttpLogExporter{
		endpoint: endpoint,
		headers:  headers,
	}
}

func (o *OtlpHttpLogExporter) Name() string {
	return "otlp-http-logs"
}

func (o *OtlpHttpLogExporter) GetSlogHandler(
	ctx context.Context,
	res *resource.Resource,
	cfg *telemetryOptions,
) (slog.Handler, error) {
	endpoint, err := otlpHTTPURL(o.endpoint, "logs")
	if err != nil {
		return nil, err
	}

	exporter, err := otlploghttp.New(
		ctx,
		otlploghttp.WithEndpointURL(endpoint),
		otlploghttp.WithHeaders(o.headers),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP log exporter: %w", err)
	}

	o.loggerProvider = sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(
			exporter,
			sdklog.WithMaxQueueSize(cfg.queueSize),
			sdklog.WithExportMaxBatchSize(cfg.batchSize),
			sdklog.WithExportInterval(cfg.batchTimeout),
		)),
	)

	return otelslog.NewHandler(
		cfg.serviceName,
		otelslog.WithLoggerProvider(o.loggerProvider),
		otelslog.WithSource(true),
	), nil
}

func (o *OtlpHttpLogExporter) Shutdown(ctx context.Context) error {
	if o.loggerProvider == nil {
		return nil
	}

	if err := o.loggerProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown OTLP HTTP log exporter: %w", err)
	}

	return nil
}

var _ LogExporter = (*OtlpHttpLogExporter)(nil)

func otlpHTTPURL(endpoint, signal string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", fmt.Errorf("OTLP endpoint cannot be empty")
	}

	if !strings.HasPrefix(endpoint, "https://") && !strings.HasPrefix(endpoint, "http://") {
		return "", fmt.Errorf("OTLP HTTP endpoint must include an http or https scheme")
	}

	return strings.TrimRight(endpoint, "/") + "/v1/" + signal, nil
}
