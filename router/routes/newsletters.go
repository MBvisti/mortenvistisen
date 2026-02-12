package routes

import (
	"mortenvistisen/internal/routing"
)

const NewsletterPrefix = "/newsletters"

var NewsletterOverview = routing.NewSimpleRoute(
	"",
	"newsletters.overview",
	NewsletterPrefix,
)

var NewsletterIndex = routing.NewSimpleRoute(
	"",
	"newsletters.index",
	AdminPrefix+NewsletterPrefix,
)

var Newsletter = routing.NewRouteWithSlug(
	"/:slug",
	"newsletters.show.slug",
	NewsletterPrefix,
)

var NewsletterShow = routing.NewRouteWithSerialID(
	"/:id",
	"newsletters.show",
	AdminPrefix+NewsletterPrefix,
)

var NewsletterNew = routing.NewSimpleRoute(
	"/new",
	"newsletters.new",
	AdminPrefix+NewsletterPrefix,
)

var NewsletterCreate = routing.NewSimpleRoute(
	"",
	"newsletters.create",
	AdminPrefix+NewsletterPrefix,
)

var NewsletterEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"newsletters.edit",
	AdminPrefix+NewsletterPrefix,
)

var NewsletterUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"newsletters.update",
	AdminPrefix+NewsletterPrefix,
)

var NewsletterDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"newsletters.destroy",
	AdminPrefix+NewsletterPrefix,
)
