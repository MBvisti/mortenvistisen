package routes

import (
	"mortenvistisen/internal/routing"
)

const AdminDashboardPrefix = "/admin"

var AdminDashboardIndex = routing.NewSimpleRoute(
	"",
	"admin.dashboards.index",
	AdminDashboardPrefix,
)
