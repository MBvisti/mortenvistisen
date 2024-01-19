package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/pkg/mail"
	"github.com/MBvisti/mortenvistisen/pkg/queue"
	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

func main() {
	router := echo.New()

	logger := telemetry.SetupLogger()

	slog.SetDefault(logger)

	// Middleware
	router.Use(slogecho.New(logger))
	router.Use(middleware.Recover())

	dbCtx := context.Background()
	conn := database.SetupDatabasePool(dbCtx, config.Cfg.GetDatabaseURL())
	defer conn.Close()

	db := database.New(conn)

	q := queue.New(db)
	if err := q.InitilizeRepeatingJobs(context.Background(), nil); err != nil {
		panic(err)
	}

	postmark := mail.NewPostmark(os.Getenv("POSTMARK_API_TOKEN"))

	mailClient := mail.NewMail(&postmark)
	tokenManager := tokens.NewManager()
	postManager := posts.NewPostManager()

	controllers := controllers.NewController(*db, mailClient, *tokenManager, *q, postManager)

	server := routes.NewServer(router, controllers, logger)

	server.Start()
}
