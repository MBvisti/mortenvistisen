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
	"admin.articles.new",
	AdminPrefix+ArticlePrefix,
)

var ArticleCreate = routing.NewSimpleRoute(
	"",
	"admin.articles.create",
	AdminPrefix+ArticlePrefix,
)

var ArticleEdit = routing.NewRouteWithID(
	"/:id/edit",
	"admin.articles.edit",
	AdminPrefix+ArticlePrefix,
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
