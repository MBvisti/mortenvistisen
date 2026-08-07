package config

import "github.com/caarlos0/env/v11"

type telemetry struct {
	ServiceName     string  `env:"TELEMETRY_SERVICE_NAME"     envDefault:"mortenvistisen"`
	ServiceVersion  string  `env:"TELEMETRY_SERVICE_VERSION"  envDefault:"1.0.0"`
	OtlpEndpoint    string  `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:""`
	OtlpHeaders     string  `env:"OTEL_EXPORTER_OTLP_HEADERS"  envDefault:""`
	OtlpProtocol    string  `env:"OTEL_EXPORTER_OTLP_PROTOCOL" envDefault:""`
	TraceSampleRate float64 `env:"TRACE_SAMPLE_RATE"          envDefault:"1.0"`
	BatchSize       int     `env:"TELEMETRY_BATCH_SIZE"       envDefault:"512"`
	BatchTimeoutMs  int     `env:"TELEMETRY_BATCH_TIMEOUT_MS" envDefault:"5000"`
}

func newTelemetryConfig() telemetry {
	telemetryCfg := telemetry{}

	if err := env.ParseWithOptions(&telemetryCfg, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		panic(err)
	}

	return telemetryCfg
}
