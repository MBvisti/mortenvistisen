package middleware

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func (m MW) Logging() echo.MiddlewareFunc {
	otelMiddleware := otelecho.Middleware(
		"mortenvistisen",
		otelecho.WithTracerProvider(m.tp),
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.Contains(c.Request().URL.Path, "/assets/") {
				return next(c)
			}
			if strings.Contains(c.Request().URL.Path, "/api/v1/health") {
				return next(c)
			}

			var ctx context.Context
			start := time.Now()

			m.httpInFlight.Add(ctx, 1)

			wrappedNext := func(c echo.Context) error {
				ctx = c.Request().Context()

				err := next(c)

				return err
			}

			err := otelMiddleware(wrappedNext)(c)

			m.httpInFlight.Add(ctx, -1)

			statusCode := c.Response().Status

			attrs := []attribute.KeyValue{
				attribute.String("method", c.Request().Method),
				attribute.String("route", c.Path()),
				attribute.Int("status_code", statusCode),
			}

			m.httpRequestsTotal.Add(
				ctx,
				1,
				metric.WithAttributes(attrs...),
			)

			requestDuration := time.Since(start)
			m.httpDuration.Record(
				ctx,
				requestDuration.Seconds(),
				metric.WithAttributes(attrs...),
			)

			slog.InfoContext(ctx, "HTTP request completed",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", statusCode,
				"duration", requestDuration.Seconds(),
				"remote_addr", c.RealIP(),
				"user_agent", c.Request().UserAgent(),
			)

			return err
		}
	}
}
