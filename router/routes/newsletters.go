package routes

import (
	"mortenvistisen/internal/routing"
)

const NewsletterPrefix = "/newsletters"

var NewsletterIndex = routing.NewSimpleRoute(
	"",
	"newsletters.index",
	NewsletterPrefix,
)

var NewsletterShow = routing.NewRouteWithID(
	"/:id",
	"newsletters.show",
	NewsletterPrefix,
)

var NewsletterShowSlug = routing.NewRouteWithSlug(
	"/:slug",
	"newsletters.show",
	NewsletterPrefix,
)

var NewsletterNew = routing.NewSimpleRoute(
	"/new",
	"newsletters.new",
	NewsletterPrefix,
)

var NewsletterCreate = routing.NewSimpleRoute(
	"",
	"newsletters.create",
	NewsletterPrefix,
)

var NewsletterEdit = routing.NewRouteWithID(
	"/:id/edit",
	"newsletters.edit",
	NewsletterPrefix,
)

var NewsletterUpdate = routing.NewRouteWithID(
	"/:id",
	"newsletters.update",
	NewsletterPrefix,
)

var NewsletterDestroy = routing.NewRouteWithID(
	"/:id",
	"newsletters.destroy",
	NewsletterPrefix,
)
