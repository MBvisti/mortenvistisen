package main

import (
	"context"
	"log/slog"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/http"
	mw "github.com/MBvisti/mortenvistisen/http/middleware"
	"github.com/MBvisti/mortenvistisen/http/router"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/pkg/mail_client"
	"github.com/MBvisti/mortenvistisen/pkg/queue"
	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/repository/psql"
	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/MBvisti/mortenvistisen/services"
)

func main() {
	logger := telemetry.SetupLogger()
	slog.SetDefault(logger)

	cfg := config.New()

	conn, err := psql.CreatePooledConnection(
		context.Background(),
		cfg.Db.GetUrlString(),
	)
	if err != nil {
		panic(err)
	}

	db := database.New(conn)
	psql := psql.NewPostgres(conn)

	// postmark := mail.NewPostmark(cfg.ExternalProviders.PostmarkApiToken)
	awsSes := mail_client.NewAwsSimpleEmailService()

	tokenManager := tokens.NewManager(cfg.Auth.TokenSigningKey)

	riverClient := queue.NewClient(conn, queue.WithLogger(logger))

	postManager := posts.NewPostManager()

	// validator := validator.New()
	// validator.RegisterStructValidation(
	// 	services.PasswordMatchValidation,
	// 	services.NewUserValidation{},
	// )
	// validator.RegisterStructValidation(
	// 	services.ResetPasswordMatchValidation,
	// 	services.UpdateUserValidation{},
	// )

	mailService := services.NewEmailSvc(cfg, &awsSes)
	// newsletterUsecase := usecases.NewNewsletter(*db, validator, mailService)

	tknService := services.NewTokenSvc(psql, cfg.Auth.TokenSigningKey)
	authService := services.NewAuth(cfg, psql)

	newsletterModel := models.NewNewsletterSvc(psql, psql, tknService, &mailService)
	subModel := models.NewSubscriberSvc(&mailService, tknService, psql)
	userModel := models.NewUserSvc(authService, psql)
	tagModel := models.NewTagSvc()
	articleModel := models.NewArticleSvc(psql)

	cookieStore := controllers.NewCookieStore(cfg.Auth.SessionKey)

	controllerDeps := controllers.NewDependencies(
		*db,
		*tokenManager,
		riverClient,
		postManager,
		mailService,
		*tknService,
		authService,
		newsletterModel,
		subModel,
		articleModel,
		tagModel,
		userModel,
		psql,
		cookieStore,
	)

	middleware := mw.NewMiddleware(authService)
	router := router.NewRouter(controllerDeps, middleware, cfg, logger)
	router.LoadInRoutes()

	server := http.NewServer(router, logger, cfg)

	server.Start()
}
