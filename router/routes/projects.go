package routes

import (
	"mortenvistisen/internal/routing"
)

const ProjectPrefix = "/projects"

var ProjectOverview = routing.NewSimpleRoute(
	"",
	"projects.overview",
	ProjectPrefix,
)

var ProjectIndex = routing.NewSimpleRoute(
	"",
	"projects.index",
	AdminPrefix+ProjectPrefix,
)

var Project = routing.NewRouteWithSlug(
	"/:slug",
	"projects.show.slug",
	ProjectPrefix,
)

var ProjectShow = routing.NewRouteWithSerialID(
	"/:id",
	"projects.show",
	AdminPrefix+ProjectPrefix,
)

var ProjectNew = routing.NewSimpleRoute(
	"/new",
	"projects.new",
	AdminPrefix+ProjectPrefix,
)

var ProjectCreate = routing.NewSimpleRoute(
	"",
	"projects.create",
	AdminPrefix+ProjectPrefix,
)

var ProjectEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"projects.edit",
	AdminPrefix+ProjectPrefix,
)

var ProjectUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"projects.update",
	AdminPrefix+ProjectPrefix,
)

var ProjectDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"projects.destroy",
	AdminPrefix+ProjectPrefix,
)
