package main

import (
	"context"
	"log/slog"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/pkg/mail"
	"github.com/MBvisti/mortenvistisen/pkg/queue"
	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/routes"
	"github.com/MBvisti/mortenvistisen/server"
	mw "github.com/MBvisti/mortenvistisen/server/middleware"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

func main() {
	ctx := context.Background()
	router := echo.New()

	logger := telemetry.SetupLogger()
	slog.SetDefault(logger)

	cfg := config.New()

	// Middleware
	router.Use(slogecho.New(logger))
	router.Use(middleware.Recover())

	conn := database.SetupDatabasePool(context.Background(), cfg.Db.GetUrlString())
	db := database.New(conn)

	postmark := mail.NewPostmark(cfg.ExternalProviders.PostmarkApiToken)
	mailClient := mail.NewMail(&postmark)

	tokenManager := tokens.NewManager(cfg.Auth.TokenSigningKey)

	authSessionStore := sessions.NewCookieStore([]byte(cfg.Auth.SessionKey), []byte(cfg.Auth.SessionEncryptionKey))

	queueDbPool, err := pgxpool.New(context.Background(), cfg.Db.GetQueueUrlString())
	if err != nil {
		panic(err)
	}

	if err := queueDbPool.Ping(ctx); err != nil {
		panic(err)
	}

	riverClient := queue.NewClient(queueDbPool, queue.WithLogger(logger))

	postManager := posts.NewPostManager()

	controllers := controllers.NewController(*db, mailClient, *tokenManager, cfg, riverClient, postManager, authSessionStore)

	serverMW := mw.NewMiddleware(authSessionStore)

	routes := routes.NewRoutes(controllers, serverMW, cfg)
	router = routes.SetupRoutes()

	server := server.NewServer(router, controllers, logger, cfg)

	server.Start()
}
