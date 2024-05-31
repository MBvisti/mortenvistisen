package main

import (
	"context"
	"log/slog"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/http"
	mw "github.com/MBvisti/mortenvistisen/http/middleware"
	"github.com/MBvisti/mortenvistisen/http/router"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/pkg/mail"
	"github.com/MBvisti/mortenvistisen/pkg/queue"
	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/usecases"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	logger := telemetry.SetupLogger()
	slog.SetDefault(logger)

	cfg := config.New()

	conn := database.SetupDatabasePool(
		context.Background(),
		cfg.Db.GetUrlString(),
	)
	db := database.New(conn)

	// postmark := mail.NewPostmark(cfg.ExternalProviders.PostmarkApiToken)
	awsSes := mail.NewAwsSimpleEmailService()
	mailClient := mail.NewMail(&awsSes)

	tokenManager := tokens.NewManager(cfg.Auth.TokenSigningKey)

	authSessionStore := sessions.NewCookieStore(
		[]byte(cfg.Auth.SessionKey),
		[]byte(cfg.Auth.SessionEncryptionKey),
	)

	queueDbPool, err := pgxpool.New(
		context.Background(),
		cfg.Db.GetQueueUrlString(),
	)
	if err != nil {
		panic(err)
	}

	if err := queueDbPool.Ping(ctx); err != nil {
		panic(err)
	}

	riverClient := queue.NewClient(queueDbPool, queue.WithLogger(logger))

	postManager := posts.NewPostManager()

	validator := validator.New()
	validator.RegisterStructValidation(
		services.PasswordMatchValidation,
		services.NewUserValidation{},
	)
	validator.RegisterStructValidation(
		services.ResetPasswordMatchValidation,
		services.UpdateUserValidation{},
	)

	newsletterUsecase := usecases.NewNewsletter(*db, validator, mailClient)

	controllerDeps := controllers.NewDependencies(
		*db,
		*tokenManager,
		riverClient,
		validator,
		postManager,
		mailClient,
		authSessionStore,
		newsletterUsecase,
	)

	middleware := mw.NewMiddleware(authSessionStore)
	router := router.NewRouter(controllerDeps, middleware, cfg, logger)
	router.LoadInRoutes()

	server := http.NewServer(router, logger, cfg)

	server.Start()
}
