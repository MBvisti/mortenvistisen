package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/MBvisti/mortenvistisen/config"
	"github.com/MBvisti/mortenvistisen/http"
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/MBvisti/mortenvistisen/queue"
	"github.com/MBvisti/mortenvistisen/queue/workers"
	"github.com/MBvisti/mortenvistisen/routes"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/lmittmann/tint"
	"github.com/maypok86/otter"
	"riverqueue.com/riverui"
)

var appRelease string

func developmentLogger() *slog.Logger {
	return slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	)
}

func productionLogger() *slog.Logger {
	return slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelError,
			TimeFormat: time.Kitchen,
		}),
	)
}

func run(ctx context.Context) error {
	cfg := config.NewConfig()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// otel := telemetry.NewOtel()
	// defer func() {
	// 	if err := otel.Shutdown(); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// appTracer := otel.NewTracer("app/tracer")

	// client := telemetry.NewTelemetry(
	// 	cfg,
	// 	appRelease,
	// 	strings.ToLower(cfg.ProjectName),
	// )
	// if client != nil {
	// 	defer client.Stop()
	// }

	slog.SetDefault(productionLogger())

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
	riverClient := queue.NewClient(
		conn,
		queue.WithWorkers(queueWorkers),
	)
	psql := psql.NewPostgres(conn, riverClient)

	opts := &riverui.ServerOpts{
		Client: riverClient,
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

	// Start the server to initialize background processes for caching and periodic queries:
	if err := riverUI.Start(ctx); err != nil {
		return err
	}

	authSvc := services.NewAuth(psql)
	emailSvc := services.NewMail()

	cacheBuilder, err := otter.NewBuilder[string, string](20)
	if err != nil {
		return err
	}

	pageCacher, err := cacheBuilder.WithVariableTTL().Build()
	if err != nil {
		return err
	}

	postManager := posts.NewManager()

	handlers := handlers.NewHandlers(
		psql,
		pageCacher,
		authSvc,
		emailSvc,
		postManager,
	)

	routes := routes.NewRoutes(
		handlers,
		riverUI,
	)

	router, c := routes.SetupRoutes(ctx)

	server := http.NewServer(c, router)

	return server.Start(c)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
