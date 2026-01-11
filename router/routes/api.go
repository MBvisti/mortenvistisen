package routes

import (
	"mortenvistisen/internal/routing"
)

const APIPrefix = "/api"

var Health = routing.NewSimpleRoute(
	"/health",
	"api.health",
	APIPrefix,
)
