// Package controllers provides HTTP handlers for the web application.
package controllers

import (
	"mortenvistisen/controllers/admin"
	"mortenvistisen/controllers/api"
	"mortenvistisen/router"

	"go.uber.org/fx"
)

var otherCache = NewCacheBuilder[string]().WithSize(2).Build

var constructors = fx.Provide(
	otherCache,
	NewPages,
	NewAssets,
	NewAPI,
	api.NewArticles,
	NewSessions,
	NewRegistrations,
	NewConfirmations,
	NewSubscriberConfirmations,
	NewResetPasswords,
	admin.NewArticles,
	admin.NewTags,
	admin.NewNewsletters,
	admin.NewSubscribers,
	admin.NewProjects,
	admin.NewDashboards,
	admin.NewJobs,
	NewPosts,
	NewNewsletters,
	NewProjects,
)

var Module = fx.Module(
	"controllers",
	constructors,
	fx.Invoke(func(r *router.Router, c Pages) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Assets) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c API) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c api.Articles) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Sessions) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Registrations) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Confirmations) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c SubscriberConfirmations) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c ResetPasswords) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Articles) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Tags) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Newsletters) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Subscribers) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Projects) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Dashboards) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Jobs) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Posts) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Newsletters) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Projects) error {
		return c.RegisterRoutes(r)
	}),
)
