package routes

import (
	"mortenvistisen/internal/routing"
)

var HomePage = routing.NewSimpleRoute(
	"/",
	"pages.home",
	"",
)

var AboutPage = routing.NewSimpleRoute(
	"/about",
	"pages.about",
	"",
)
