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
var NewsletterShow = routing.NewRouteWithSlug(
	"/:slug",
	"newsletters.show",
	NewsletterPrefix,
)
