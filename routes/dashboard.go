package routes

import (
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/views/paths"
	"github.com/labstack/echo/v4"
)

func dashboardRoutes(
	router *echo.Echo,
	ctrl handlers.Dashboard,
) {
	dashboardRouter := router.Group("/dashboard")

	dashboardRouter.GET("", func(c echo.Context) error {
		return ctrl.Index(c)
	}).Name = paths.DashboardHomePage
}
