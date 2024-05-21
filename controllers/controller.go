package controllers

import (
	"github.com/MBvisti/mortenvistisen/pkg/mail"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

// Actions:
// - index | GET
// - create | GET
// - store | POST
// - show | GET
// - edit | GET
// - update | PUT/PATCH
// - destroy | DELETE

type Dependencies struct {
	DB          database.Queries
	TknManager  tokens.Manager
	QueueClient *river.Client[pgx.Tx]
	Validate    *validator.Validate
	PostManager posts.PostManager
	Mail        mail.Mail
	AuthStore   *sessions.CookieStore
}
