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
	"github.com/MBvisti/mortenvistisen/queue/jobs"
	"github.com/MBvisti/mortenvistisen/queue/workers"
	"github.com/MBvisti/mortenvistisen/routes"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/telemetry"
	"github.com/a-h/templ"
	"github.com/dromara/carbon/v2"
	"github.com/maypok86/otter"
	"github.com/riverqueue/river"
	"riverqueue.com/riverui"
)

var appRelease string

func run(ctx context.Context) error {
	cfg := config.NewConfig()
	carbon.SetDefault(carbon.Default{
		Layout:       carbon.DateTimeLayout,
		Timezone:     carbon.UTC,
		WeekStartsAt: carbon.Monday,
		Locale:       "en",
	})

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

	if cfg.Environment == config.DEV_ENVIRONMENT {
		slog.SetDefault(telemetry.DevelopmentLogger())
	}
	if cfg.Environment == config.PROD_ENVIRONMENT {
		slog.SetDefault(telemetry.ProductionLogger())
	}

	conn, err := psql.CreatePooledConnection(
		ctx,
		cfg.GetDatabaseURL(),
	)
	if err != nil {
		return err
	}

	emailSvc := services.NewMail()

	queueWorkers, err := workers.SetupWorkers(workers.WorkerDependencies{
		Emailer: emailSvc,
		Conn:    conn,
	})
	if err != nil {
		return err
	}

	periodicJobs := []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(24*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.SubscriberCleanupJobArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
	}

	riverClient := queue.NewClient(
		conn,
		queue.WithWorkers(queueWorkers),
		queue.WithPeriodicJobs(periodicJobs),
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

	cacheBuilder, err := otter.NewBuilder[string, templ.Component](20)
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

	if err := riverClient.Start(c); err != nil {
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
