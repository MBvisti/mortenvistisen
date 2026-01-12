package routes

import (
	"mortenvistisen/internal/routing"
)

const ArticlePrefix = "/articles"

var ArticleIndex = routing.NewSimpleRoute(
	"",
	"articles.index",
	"/posts",
)

var ArticleShow = routing.NewRouteWithID(
	"/:id",
	"articles.show",
	"/posts",
)

var ArticleShowSlug = routing.NewRouteWithSlug(
	"/:slug",
	"articles.slug",
	"/posts",
)

var ArticleNew = routing.NewSimpleRoute(
	"/new",
	"articles.new",
	ArticlePrefix,
)

var ArticleCreate = routing.NewSimpleRoute(
	"",
	"articles.create",
	ArticlePrefix,
)

var ArticleEdit = routing.NewRouteWithID(
	"/:id/edit",
	"articles.edit",
	ArticlePrefix,
)

var ArticleUpdate = routing.NewRouteWithID(
	"/:id",
	"articles.update",
	ArticlePrefix,
)

var ArticleDestroy = routing.NewRouteWithID(
	"/:id",
	"articles.destroy",
	ArticlePrefix,
)
