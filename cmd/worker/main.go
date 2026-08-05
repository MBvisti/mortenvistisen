package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	mailclients "mortenvistisen/clients/email"
	objectstorage "mortenvistisen/clients/storage"
	"mortenvistisen/config"
	"mortenvistisen/database"
	"mortenvistisen/email"
	"mortenvistisen/internal/server"
	"mortenvistisen/queue"
	"mortenvistisen/telemetry"

	"go.uber.org/fx"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app := fx.New(
		fx.Provide(
			func() context.Context { return ctx },
			func(cfg config.Config) (email.TransactionalSender, email.MarketingSender) {
				if config.Env == server.ProdEnvironment {
					sender := mailclients.NewAwsSes(cfg)
					return sender, sender
				}

				return mailclients.NewMailpit(cfg), mailclients.NewMailpit(cfg)
			},
		),

		config.Module,
		objectstorage.Module,
		database.Module,
		telemetry.Module,
		queue.Module,
		queue.WorkersModule,

		fx.Invoke(startQueueProcessor),
	)

	if err := app.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Stop(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func startQueueProcessor(lc fx.Lifecycle, appCtx context.Context, p queue.Processor) {
	var done <-chan struct{}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			slog.InfoContext(appCtx, "starting queue processor")
			done = startInBackground(appCtx, "queue processor", p.Start)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return stopAndWait(ctx, p.Stop, done)
		},
	})
}

func startInBackground(
	ctx context.Context,
	name string,
	start func(context.Context) error,
) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := start(ctx); err != nil {
			slog.Error(name+" error", "error", err)
		}
	}()
	return done
}

func stopAndWait(
	ctx context.Context,
	stop func(context.Context) error,
	done <-chan struct{},
) error {
	stopErr := stop(ctx)
	select {
	case <-done:
		return stopErr
	case <-ctx.Done():
		return errors.Join(stopErr, ctx.Err())
	}
}
