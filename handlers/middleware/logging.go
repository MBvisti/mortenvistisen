package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mbvisti/mortenvistisen/config"
	"github.com/mbvisti/mortenvistisen/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func (m MW) Logging() echo.MiddlewareFunc {
	otelMiddleware := otelecho.Middleware(
		config.Cfg.ServiceName,
		otelecho.WithTracerProvider(m.tp),
	)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Only skip health checks, but measure everything else including assets
			if strings.Contains(c.Request().URL.Path, "/api/v1/health") {
				return next(c)
			}

			// Start timing HERE, before anything else
			start := time.Now()

			// Track in-flight requests
			ctx := c.Request().Context()
			m.httpInFlight.Add(ctx, 1)
			defer m.httpInFlight.Add(ctx, -1)

			// Calculate request size before processing
			reqSize := telemetry.ComputeApproximateRequestSize(c.Request())

			// Create wrapped handler that captures the response status
			var statusCode int
			wrappedNext := func(c echo.Context) error {
				err := next(c)
				statusCode = c.Response().Status
				return err
			}

			// Execute the request through OpenTelemetry middleware
			err := otelMiddleware(wrappedNext)(c)

			// Calculate duration AFTER everything completes
			duration := time.Since(start)

			// Record metrics with enhanced labels following Echo's pattern
			url := c.Path() // route pattern like /users/:id
			if url == "" {
				// For 404 cases, use actual path to have distinction
				url = c.Request().URL.Path
			}

			attrs := []attribute.KeyValue{
				attribute.String("method", c.Request().Method),
				attribute.String("route", url),
				attribute.String("host", c.Request().Host),
				attribute.Int("status_code", statusCode),
			}

			m.httpRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
			m.httpDuration.Record(
				ctx,
				duration.Seconds(),
				metric.WithAttributes(attrs...),
			)
			m.httpRequestSize.Record(
				ctx,
				float64(reqSize),
				metric.WithAttributes(attrs...),
			)
			m.httpResponseSize.Record(
				ctx,
				float64(c.Response().Size),
				metric.WithAttributes(attrs...),
			)

			// Only log non-asset requests to avoid log spam
			if !strings.Contains(c.Request().URL.Path, "/assets/") {
				slog.InfoContext(ctx, "HTTP request completed",
					"method", c.Request().Method,
					"path", c.Request().URL.Path,
					"status", statusCode,
					"duration", duration.Seconds(),
					"remote_addr", c.RealIP(),
					"user_agent", c.Request().UserAgent(),
				)
			}

			return err
		}
	}
}

// func (m MW) Logging() echo.MiddlewareFunc {
// 	otelMiddleware := otelecho.Middleware(
// 		"mortenvistisen",
// 		otelecho.WithTracerProvider(m.tp),
// 	)
//
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			if strings.Contains(c.Request().URL.Path, "/assets/") {
// 				return next(c)
// 			}
// 			if strings.Contains(c.Request().URL.Path, "/api/v1/health") {
// 				return next(c)
// 			}
//
// 			var ctx context.Context
// 			start := time.Now()
//
// 			m.httpInFlight.Add(ctx, 1)
//
// 			wrappedNext := func(c echo.Context) error {
// 				ctx = c.Request().Context()
//
// 				err := next(c)
//
// 				// Record metrics and logging after request completion
// 				m.httpInFlight.Add(ctx, -1)
//
// 				statusCode := c.Response().Status
//
// 				attrs := []attribute.KeyValue{
// 					attribute.String("method", c.Request().Method),
// 					attribute.String("route", c.Path()),
// 					attribute.Int("status_code", statusCode),
// 				}
//
// 				m.httpRequestsTotal.Add(
// 					ctx,
// 					1,
// 					metric.WithAttributes(attrs...),
// 				)
//
// 				requestDuration := time.Since(start)
// 				m.httpDuration.Record(
// 					ctx,
// 					requestDuration.Seconds(),
// 					metric.WithAttributes(attrs...),
// 				)
//
// 				slog.InfoContext(ctx, "HTTP request completed",
// 					"method", c.Request().Method,
// 					"path", c.Request().URL.Path,
// 					"status", statusCode,
// 					"duration", requestDuration.Seconds(),
// 					"remote_addr", c.RealIP(),
// 					"user_agent", c.Request().UserAgent(),
// 				)
//
// 				return err
// 			}
//
// 			err := otelMiddleware(wrappedNext)(c)
//
// 			return err
// 		}
// 	}
// }
