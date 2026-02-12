package routes

import (
	"mortenvistisen/internal/routing"
)

const ArticlePrefix = "/articles"

var ArticleOverview = routing.NewSimpleRoute(
	"",
	"articles.overview",
	"/posts",
)

var ArticleIndex = routing.NewSimpleRoute(
	"",
	"articles.index",
	AdminPrefix+ArticlePrefix,
)

// Article TODO: naming
var Article = routing.NewRouteWithSlug(
	"/:slug",
	"articles.show.slug",
	"/posts",
)

var ArticleShow = routing.NewRouteWithSerialID(
	"/:id",
	"articles.show",
	AdminPrefix+ArticlePrefix,
)

var ArticleNew = routing.NewSimpleRoute(
	"/new",
	"articles.new",
	AdminPrefix+ArticlePrefix,
)

var ArticleCreate = routing.NewSimpleRoute(
	"",
	"articles.create",
	AdminPrefix+ArticlePrefix,
)

var ArticleEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"articles.edit",
	AdminPrefix+ArticlePrefix,
)

var ArticleUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"articles.update",
	AdminPrefix+ArticlePrefix,
)

var ArticleDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"articles.destroy",
	AdminPrefix+ArticlePrefix,
)
