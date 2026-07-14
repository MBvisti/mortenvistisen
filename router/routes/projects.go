package routes

import (
	"mortenvistisen/internal/routing"
)

const ProjectPrefix = "/projects"

var ProjectIndex = routing.NewSimpleRoute(
	"",
	"projects.index",
	ProjectPrefix,
)
var ProjectShow = routing.NewRouteWithSlug(
	"/:slug",
	"projects.show",
	ProjectPrefix,
)
