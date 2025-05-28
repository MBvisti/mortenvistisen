package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"time"

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

var appVersion string

func migrate(ctx context.Context) error {
	ctx, cancel := context.WithTimeoutCause(
		ctx,
		5*time.Minute,
		errors.New("migration timeout of 5 minutes reached"),
	)
	defer cancel()

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

	ops := flag.String("cmd", "", "")
	version := flag.String("version", "", "")
	flag.Parse()

	if *version != "" {
		v, err := strconv.Atoi(*version)
		if err != nil {
			return err
		}

		switch *ops {
		case "up":
			_, err = gooseProvider.UpTo(ctx, int64(v))
			if err != nil {
				return err
			}
		case "down":
			_, err = gooseProvider.DownTo(ctx, int64(v))
			if err != nil {
				return err
			}
		}
	}

	if *version == "" {
		switch *ops {
		case "up":
			_, err = gooseProvider.Up(ctx)
			if err != nil {
				return err
			}
		case "down":
			_, err = gooseProvider.Down(ctx)
			if err != nil {
				return err
			}
		case "upbyone":
			_, err = gooseProvider.UpByOne(ctx)
			if err != nil {
				return err
			}
		case "reset":
			_, err = gooseProvider.DownTo(ctx, 0)
			if err != nil {
				return err
			}
		case "status":
			statuses, err := gooseProvider.Status(ctx)
			if err != nil {
				return err
			}

			for _, status := range statuses {
				slog.Info(
					"database status",
					"version",
					status.Source.Version,
					"file_name",
					status.Source.Path,
					"state",
					status.State,
					"applied_at",
					status.AppliedAt,
				)
			}
		}
	}

	return nil
}

func run(ctx context.Context) error {
	cfg := config.NewConfig()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := migrate(ctx); err != nil {
		panic(err)
	}

	tel, err := telemetry.New(
		ctx,
		appVersion,
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
	if cfg.Environment == config.PROD_ENVIRONMENT {
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
	queueWorkers, err := workers.SetupWorkers(workers.WorkerDependencies{})
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
	if cfg.Environment == config.PROD_ENVIRONMENT {
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

	emailClient := clients.NewEmail()

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

	router, c := routes.SetupRoutes(ctx)

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
