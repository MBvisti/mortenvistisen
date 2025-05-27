package routes

import (
	"net/http"
)

const (
	dashboardRoutePrefix = "/dashboard"
	dashboardNamePrefix  = "dashboard"
)

var Dashboard = []Route{
	DashboardHome,
	DashboardNewArticle,
	DashboardStoreArticle,
}

var DashboardHome = Route{
	Name:        dashboardNamePrefix + ".home",
	Path:        dashboardRoutePrefix,
	Method:      http.MethodGet,
	HandlerName: "Index",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardNewArticle = Route{
	Name:        dashboardNamePrefix + ".articles.new",
	Path:        dashboardRoutePrefix + "/articles/new",
	Method:      http.MethodGet,
	HandlerName: "NewArticle",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardStoreArticle = Route{
	Name:        dashboardNamePrefix + ".articles.create",
	Path:        dashboardRoutePrefix + "/articles/new",
	Method:      http.MethodPost,
	HandlerName: "StoreArticle",
	Middleware: []string{
		"AuthOnly",
	},
}
