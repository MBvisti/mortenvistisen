package routes

import (
	"mortenvistisen/internal/routing"
)

const AdminArticlePrefix = "/admin/articles"

var AdminArticleIndex = routing.NewSimpleRoute(
	"",
	"admin.articles.index",
	AdminArticlePrefix,
)
var AdminArticleShow = routing.NewRouteWithSerialID(
	"/:id",
	"admin.articles.show",
	AdminArticlePrefix,
)
var AdminArticleNew = routing.NewSimpleRoute(
	"/new",
	"admin.articles.new",
	AdminArticlePrefix,
)
var AdminArticleCreate = routing.NewSimpleRoute(
	"",
	"admin.articles.create",
	AdminArticlePrefix,
)
var AdminArticleEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"admin.articles.edit",
	AdminArticlePrefix,
)
var AdminArticleUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"admin.articles.update",
	AdminArticlePrefix,
)
var AdminArticleDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"admin.articles.destroy",
	AdminArticlePrefix,
)
