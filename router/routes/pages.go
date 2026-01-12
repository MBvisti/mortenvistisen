package routes

import (
	"mortenvistisen/internal/routing"
)

var HomePage = routing.NewSimpleRoute(
	"/",
	"pages.home",
	"",
)

var ProjectsPage = routing.NewSimpleRoute(
	"/projects",
	"pages.projects",
	"",
)
