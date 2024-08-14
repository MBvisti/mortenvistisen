package main

import (
	"context"
	"log/slog"

	"github.com/MBvisti/mortenvistisen/http"
	"github.com/MBvisti/mortenvistisen/http/handlers"
	mw "github.com/MBvisti/mortenvistisen/http/middleware"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/pkg/mail_client"
	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/MBvisti/mortenvistisen/psql/database"
	"github.com/MBvisti/mortenvistisen/queue"
	"github.com/MBvisti/mortenvistisen/routes"
	"github.com/MBvisti/mortenvistisen/services"
)

// version is the latest commit sha at build time
var version string

func main() {
	cfg := config.New(version)

	otel := telemetry.NewOtel(cfg)
	defer func() {
		if err := otel.Shutdown(); err != nil {
			panic(err)
		}
	}()

	blogTracer := otel.NewTracer("blog/tracer")

	client := telemetry.NewTelemetry(cfg, version, cfg.App.ProjectName)
	if client != nil {
		defer client.Stop()
	}

	conn, err := psql.CreatePooledConnection(
		context.Background(),
		cfg.Db.GetUrlString(),
	)
	if err != nil {
		panic(err)
	}

	db := database.New(conn)
	psql := psql.NewPostgres(conn)

	awsSes := mail_client.NewAwsSimpleEmailService()

	riverClient := queue.NewClient(conn, queue.WithLogger(slog.Default()))

	postManager := posts.NewPostManager()

	mailService := services.NewEmailSvc(cfg, &awsSes)

	tknService := services.NewTokenSvc(psql, cfg.Auth.TokenSigningKey)
	authService := services.NewAuth(cfg, psql)

	newsletterSvc := models.NewNewsletterSvc(psql, psql, tknService, &mailService)
	subscriberSvc := models.NewSubscriberSvc(&mailService, tknService, psql)
	userSvc := models.NewUserSvc(authService, psql)
	tagSvc := models.NewTagSvc(psql)
	articleSvc := models.NewArticleSvc(psql)

	cookieStore := handlers.NewCookieStore(cfg.Auth.SessionKey)

	baseHandlers := handlers.NewDependencies(
		*db,
		cfg,
		riverClient,
		cookieStore,
		blogTracer,
	)

	apiHandlers := handlers.NewApi(baseHandlers)
	appHandlers := handlers.NewApp(
		baseHandlers,
		articleSvc,
		subscriberSvc,
		postManager,
		*tknService,
	)
	authHandlers := handlers.NewAuthentication(
		baseHandlers,
		authService,
		userSvc,
		mailService,
		*tknService,
		cfg,
	)
	dashboardHandlers := handlers.NewDashboard(
		baseHandlers,
		articleSvc,
		tagSvc,
		postManager,
		newsletterSvc,
		subscriberSvc,
		*tknService,
		mailService,
	)
	registerHanlders := handlers.NewRegistration(
		baseHandlers,
		authService,
		userSvc,
		*tknService,
		mailService,
	)

	middleware := mw.NewMiddleware(authService)
	router := routes.NewRouter(
		middleware,
		cfg,
		apiHandlers,
		appHandlers,
		authHandlers,
		dashboardHandlers,
		registerHanlders,
		baseHandlers,
	)
	router.LoadInRoutes()

	server := http.NewServer(router, cfg)

	server.Start()
}
