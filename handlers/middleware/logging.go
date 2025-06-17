package middleware

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mbvisti/mortenvistisen/config"
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
			if strings.Contains(c.Request().URL.Path, "/assets/") ||
				strings.Contains(c.Request().URL.Path, "/api/v1/health") {
				return next(c)
			}

			start := time.Now()
			var ctx context.Context
			var statusCode int

			m.httpInFlight.Add(c.Request().Context(), 1)
			defer m.httpInFlight.Add(c.Request().Context(), -1)

			wrappedNext := func(c echo.Context) error {
				ctx = c.Request().Context()
				err := next(c)
				statusCode = c.Response().Status
				return err
			}

			err := otelMiddleware(wrappedNext)(c)

			// Record metrics AFTER everything completes
			requestDuration := time.Since(start)
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
