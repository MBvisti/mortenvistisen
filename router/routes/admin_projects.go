package routes

import (
	"mortenvistisen/internal/routing"
)

const AdminProjectPrefix = "/admin/projects"

var AdminProjectIndex = routing.NewSimpleRoute(
	"",
	"admin.projects.index",
	AdminProjectPrefix,
)
var AdminProjectShow = routing.NewRouteWithSerialID(
	"/:id",
	"admin.projects.show",
	AdminProjectPrefix,
)
var AdminProjectNew = routing.NewSimpleRoute(
	"/new",
	"admin.projects.new",
	AdminProjectPrefix,
)
var AdminProjectCreate = routing.NewSimpleRoute(
	"",
	"admin.projects.create",
	AdminProjectPrefix,
)
var AdminProjectEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"admin.projects.edit",
	AdminProjectPrefix,
)
var AdminProjectUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"admin.projects.update",
	AdminProjectPrefix,
)
var AdminProjectDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"admin.projects.destroy",
	AdminProjectPrefix,
)
