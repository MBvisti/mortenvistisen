package routes

import (
	"mortenvistisen/internal/routing"
)

const AdminTagPrefix = "/admin/tags"

var AdminTagIndex = routing.NewSimpleRoute(
	"",
	"admin.tags.index",
	AdminTagPrefix,
)
var AdminTagCreate = routing.NewSimpleRoute(
	"",
	"admin.tags.create",
	AdminTagPrefix,
)
var AdminTagUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"admin.tags.update",
	AdminTagPrefix,
)
var AdminTagDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"admin.tags.destroy",
	AdminTagPrefix,
)
