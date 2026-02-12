package routes

import (
	"mortenvistisen/internal/routing"
)

const AdminPrefix = "/admin"

var AdminHome = routing.NewSimpleRoute(
	"",
	"admin.home",
	AdminPrefix,
)
