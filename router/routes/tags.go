package routes

import (
	"mortenvistisen/internal/routing"
)

const TagPrefix = "/tags"

var TagIndex = routing.NewSimpleRoute(
	"",
	"tags.index",
	AdminPrefix+TagPrefix,
)

var TagShow = routing.NewRouteWithSerialID(
	"/:id",
	"tags.show",
	AdminPrefix+TagPrefix,
)

var TagNew = routing.NewSimpleRoute(
	"/new",
	"tags.new",
	AdminPrefix+TagPrefix,
)

var TagCreate = routing.NewSimpleRoute(
	"",
	"tags.create",
	AdminPrefix+TagPrefix,
)

var TagEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"tags.edit",
	AdminPrefix+TagPrefix,
)

var TagUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"tags.update",
	AdminPrefix+TagPrefix,
)

var TagDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"tags.destroy",
	AdminPrefix+TagPrefix,
)
