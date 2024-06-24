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
)

func main() {
	// logger := telemetry.SetupLogger()
	// setup loki client

	// tp := trace.NewTracerProvider(
	// 	trace.WithSampler(trace.AlwaysSample()),
	// )
	// tracer := tp.Tracer("hello/world")
	//
	// ctx, span := tracer.Start(context.Background(), "foo")
	// defer span.End()

	// lokiCfg, _ := loki.NewDefaultConfig("https://monitoring.mbv-labs.com/loki/api/v1/push")
	// lokiCfg.TenantID = "local-test"
	// client, err := loki.New(lokiCfg)
	// if err != nil {
	// 	panic(err)
	// }
	// defer client.Stop()
	//
	// logger := slog.New(slogloki.Option{
	// 	Level:  slog.LevelInfo,
	// 	Client: client,
	// 	AttrFromContext: []func(ctx context.Context) []slog.Attr{
	// 		slogotel.ExtractOtelAttrFromContext([]string{"tracing"}, "trace_id", "span_id"),
	// 	},
	// }.NewLokiHandler())
	// logger = logger.
	// 	// With("environment", "dev").
	// 	// With("release", "v1.0.0").
	// 	// With("container", "mbv").
	// 	With("service_name", "mortenvistisen_blog")
	//
	// log.Print("yo cunt")
	// logger.Info("yoyoyo this shit whack")

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

	riverClient := queue.NewClient(conn, queue.WithLogger(logger))

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

	subModel := models.NewSubscriber(*db)

	controllerDeps := controllers.NewDependencies(
		*db,
		*tokenManager,
		riverClient,
		validator,
		postManager,
		mailClient,
		authSessionStore,
		newsletterUsecase,
		subModel,
	)

	middleware := mw.NewMiddleware(authSessionStore)
	router := router.NewRouter(controllerDeps, middleware, cfg, logger)
	router.LoadInRoutes()

	server := http.NewServer(router, logger, cfg)

	server.Start()
}
