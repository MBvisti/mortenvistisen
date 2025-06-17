package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mbvisti/mortenvistisen/clients"
	"github.com/mbvisti/mortenvistisen/config"
	"github.com/mbvisti/mortenvistisen/handlers"
	"github.com/mbvisti/mortenvistisen/handlers/middleware"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/psql/queue"
	"github.com/mbvisti/mortenvistisen/psql/queue/workers"
	"github.com/mbvisti/mortenvistisen/router"
	"github.com/mbvisti/mortenvistisen/server"
	"github.com/mbvisti/mortenvistisen/telemetry"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
	"riverqueue.com/riverui"
)

var AppVersion string

func migrate(ctx context.Context) error {
	slog.Info("STARTING TO MIGRATE")

	cfg := config.NewConfig()

	gooseLock, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return err
	}

	fsys, err := fs.Sub(psql.Migrations, "migrations")
	if err != nil {
		return err
	}

	pool, err := psql.CreatePooledConnection(ctx, cfg.GetDatabaseURL())
	if err != nil {
		return err
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)

	gooseProvider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		fsys,
		goose.WithVerbose(true),
		goose.WithSessionLocker(gooseLock),
	)
	if err != nil {
		return err
	}
	_, err = gooseProvider.Up(ctx)
	if err != nil {
		return err
	}

	return nil
}

func run(ctx context.Context) error {
	cfg := config.NewConfig()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var tel *telemetry.Telemetry
	if cfg.Environment == config.PROD_ENVIRONMENT {
		if err := migrate(ctx); err != nil {
			return err
		}

		t, err := telemetry.New(
			ctx,
			AppVersion,
			cfg.ServiceName,
			&telemetry.LokiExporter{
				LogLevel:   slog.LevelInfo,
				WithTraces: true,
				URL:        "https://telemetry-loki.mbvlabs.com",
				Labels: map[string]string{
					"env":     cfg.Environment,
					"service": cfg.ServiceName,
				},
			},
			telemetry.NewOtlpHttpTraceExporter(
				"telemetry-alloy.mbvlabs.com",
				false,
				map[string]string{
					"Authorization": "Basic YWRtaW46U2Ftc3VuZzIwNjE=",
				},
			),
			telemetry.NewOtlpHttpMetricExporter(
				"telemetry-alloy.mbvlabs.com",
				false,
				map[string]string{
					"Authorization": "Basic YWRtaW46U2Ftc3VuZzIwNjE=",
				},
			),
		)
		if err != nil {
			return fmt.Errorf("failed to initialize telemetry: %w", err)
		}

		defer func() {
			if err := t.Shutdown(ctx); err != nil {
				slog.Error("Failed to shutdown telemetry", "error", err)
			}
		}()

		tel = t
	}
	if cfg.Environment == config.STAGING_ENVIRONMENT {
		if err := migrate(ctx); err != nil {
			return err
		}

		t, err := telemetry.New(
			ctx,
			AppVersion,
			cfg.ServiceName,
			&telemetry.LokiExporter{
				LogLevel:   slog.LevelInfo,
				WithTraces: true,
				URL:        "https://telemetry-loki.mbvlabs.com",
				Labels: map[string]string{
					"env":     "staging",
					"service": "blog-staging",
				},
			},
			telemetry.NewOtlpHttpTraceExporter(
				"telemetry-alloy.mbvlabs.com",
				false,
				map[string]string{
					"Authorization": "Basic YWRtaW46U2Ftc3VuZzIwNjE=",
				},
			),
			telemetry.NewOtlpHttpMetricExporter(
				"telemetry-alloy.mbvlabs.com",
				false,
				map[string]string{
					"Authorization": "Basic YWRtaW46U2Ftc3VuZzIwNjE=",
				},
			),
		)
		if err != nil {
			return fmt.Errorf("failed to initialize telemetry: %w", err)
		}

		defer func() {
			if err := t.Shutdown(ctx); err != nil {
				slog.Error("Failed to shutdown telemetry", "error", err)
			}
		}()

		tel = t
	}

	if cfg.Environment == config.DEV_ENVIRONMENT {
		t, err := telemetry.New(
			ctx,
			AppVersion,
			"dev",
			&telemetry.StdoutExporter{
				LogLevel:   slog.LevelDebug,
				WithTraces: true,
			},
			&telemetry.NoopTraceExporter{},
			&telemetry.NoopMetricExporter{},
		)
		if err != nil {
			return fmt.Errorf("failed to initialize telemetry: %w", err)
		}

		tel = t
	}

	if cfg.Environment == config.STAGING_ENVIRONMENT ||
		cfg.Environment == config.PROD_ENVIRONMENT {
		defer func() {
			if err := tel.Shutdown(ctx); err != nil {
				slog.Error("Failed to shutdown telemetry", "error", err)
			}
		}()
	}

	if err := telemetry.SetupRuntimeMetricsInCallback(telemetry.GetMeter()); err != nil {
		return fmt.Errorf("failed to setup callback metrics: %w", err)
	}

	conn, err := psql.CreatePooledConnection(
		ctx,
		cfg.GetDatabaseURL(),
	)
	if err != nil {
		return err
	}

	emailClient := clients.NewEmail()

	queueWorkers, err := workers.SetupWorkers(workers.WorkerDependencies{
		DB:          conn,
		EmailClient: emailClient,
	})
	if err != nil {
		return err
	}

	queueLogger, queueLoggerShutdown := telemetry.NewLogger(
		ctx,
		&telemetry.StdoutExporter{
			LogLevel:   slog.LevelError,
			WithTraces: true,
		},
	)
	if cfg.Environment == config.STAGING_ENVIRONMENT ||
		cfg.Environment == config.PROD_ENVIRONMENT {
		defer func() {
			if err := queueLoggerShutdown(ctx); err != nil {
				slog.Error("Failed to shutdown telemetry", "error", err)
			}
		}()
	}

	psql := psql.NewPostgres(conn, nil)
	psql.NewQueue(
		queue.WithLogger(queueLogger),
		queue.WithWorkers(queueWorkers),
	)

	opts := &riverui.ServerOpts{
		Client: psql.Queue(),
		DB:     conn,
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.SetLogLoggerLevel(slog.LevelError),
		})),
		Prefix: "/river",
	}
	riverUI, err := riverui.NewServer(opts)
	if err != nil {
		return err
	}

	if err := riverUI.Start(ctx); err != nil {
		return err
	}

	handlers := handlers.NewHandlers(
		psql,
		emailClient,
	)

	mw, err := middleware.New(tel.AppTracerProvider)
	if err != nil {
		return err
	}

	routes := router.New(
		ctx,
		handlers,
		mw,
		riverUI,
		tel.AppTracerProvider,
	)

	router, c := routes.SetupRoutes(ctx, handlers.App.NotFoundPage)

	server := server.NewHttp(c, router)

	if err := psql.Queue().Start(ctx); err != nil {
		return err
	}

	return server.Start(c)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
