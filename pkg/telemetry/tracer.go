package telemetry

import (
	"context"

	"github.com/MBvisti/mortenvistisen/pkg/config"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

type NoopSpan struct {
	embedded.Span
}

// AddEvent implements trace.Span.
func (n NoopSpan) AddEvent(name string, options ...trace.EventOption) {
}

// AddLink implements trace.Span.
func (n NoopSpan) AddLink(link trace.Link) {
}

// End implements trace.Span.
func (n NoopSpan) End(options ...trace.SpanEndOption) {
}

// IsRecording implements trace.Span.
func (n NoopSpan) IsRecording() bool {
	return true
}

// RecordError implements trace.Span.
func (n NoopSpan) RecordError(err error, options ...trace.EventOption) {
}

// SetAttributes implements trace.Span.
func (n NoopSpan) SetAttributes(kv ...attribute.KeyValue) {
}

// SetName implements trace.Span.
func (n NoopSpan) SetName(name string) {
}

// SetStatus implements trace.Span.
func (n NoopSpan) SetStatus(code codes.Code, description string) {
}

// SpanContext implements trace.Span.
func (n NoopSpan) SpanContext() trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{})
}

// TracerProvider implements trace.Span.
func (n NoopSpan) TracerProvider() trace.TracerProvider {
	return tracesdk.NewTracerProvider()
}

var _ trace.Span = new(NoopSpan)

type NoopTracer struct {
	embedded.Tracer
}

// Start implements trace.Tracer.
func (n NoopTracer) Start(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return context.Background(), NoopSpan{}
}

var _ trace.Tracer = new(NoopTracer)

type Tracer struct {
	tracer trace.Tracer
	name   string
}

type Otel struct {
	cfg           config.Cfg
	traceProvider *tracesdk.TracerProvider
}

func NewOtel(cfg config.Cfg) Otel {
	sampler := tracesdk.WithSampler(tracesdk.NeverSample())
	if cfg.App.Environment == config.PROD_ENVIRONMENT {
		sampler = tracesdk.WithSampler(tracesdk.AlwaysSample())
	}

	tp := tracesdk.NewTracerProvider(
		sampler,
	)

	return Otel{
		cfg,
		tp,
	}
}

func (o Otel) NewTracer(name string) Tracer {
	var t trace.Tracer
	if o.cfg.App.Environment == config.PROD_ENVIRONMENT {
		t = o.traceProvider.Tracer(name)
	}
	if o.cfg.App.Environment == config.DEV_ENVIRONMENT {
		t = NoopTracer{}
	}

	return Tracer{
		t,
		name,
	}
}

func (o Otel) Shutdown() error {
	return o.traceProvider.Shutdown(context.Background())
}

func (t Tracer) CreateSpan(
	ctx context.Context,
	name string,
) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name)
}

func (t Tracer) CreateChildSpan(
	ctx context.Context,
	span trace.Span,
	name string,
) (context.Context, trace.Span) {
	return span.TracerProvider().Tracer(t.name).Start(ctx, name)
}
