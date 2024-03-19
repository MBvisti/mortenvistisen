package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/pkg/mail"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/riverqueue/river"
)

type Controller struct {
	db               database.Queries
	mail             mail.Mail
	validate         *validator.Validate
	tknManager       tokens.Manager
	cfg              config.Cfg
	queueClient      *river.Client[pgx.Tx]
	authSessionStore *sessions.CookieStore
	postManager      posts.PostManager
}

func NewController(
	db database.Queries, mail mail.Mail, tknManager tokens.Manager, cfg config.Cfg, qc *river.Client[pgx.Tx], pm posts.PostManager,
	authSessionStore *sessions.CookieStore,
) Controller {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return Controller{
		db,
		mail,
		validate,
		tknManager,
		cfg,
		qc,
		authSessionStore,
		pm,
	}
}

func (c *Controller) AppHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, []byte("app is healthy and running"))
}

func (c *Controller) InternalError(ctx echo.Context) error {
	var from string

	return views.InternalServerErr(ctx, views.InternalServerErrData{
		FromLocation: from,
	})
}

func (c *Controller) Redirect(ctx echo.Context) error {
	toLocation := ctx.QueryParam("to")
	if toLocation == "" {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		return c.InternalError(ctx)
	}

	ctx.Response().Writer.Header().Add("HX-Redirect", fmt.Sprintf("/%s", toLocation))

	return nil
}

func (c *Controller) formatArticleSlug(slug string) string {
	return fmt.Sprintf("posts/%s", slug)
}

func (c *Controller) buildURLFromSlug(slug string) string {
	return fmt.Sprintf("%s://%s/%s", os.Getenv("APP_SCHEME"), os.Getenv("APP_HOST"), slug)
}
