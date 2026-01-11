package routes

import (
	"mortenvistisen/internal/routing"
)

const TagPrefix = "/tags"

var TagIndex = routing.NewSimpleRoute(
	"/",
	"tags.index",
	TagPrefix,
)

var TagShow = routing.NewRouteWithID(
	"/:id",
	"tags.show",
	TagPrefix,
)

var TagNew = routing.NewSimpleRoute(
	"/new",
	"tags.new",
	TagPrefix,
)

var TagCreate = routing.NewSimpleRoute(
	"/",
	"tags.create",
	TagPrefix,
)

var TagEdit = routing.NewRouteWithID(
	"/:id/edit",
	"tags.edit",
	TagPrefix,
)

var TagUpdate = routing.NewRouteWithID(
	"/:id",
	"tags.update",
	TagPrefix,
)

var TagDestroy = routing.NewRouteWithID(
	"/:id",
	"tags.destroy",
	TagPrefix,
)
