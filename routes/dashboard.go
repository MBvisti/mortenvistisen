package routes

import (
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/views/paths"
	"github.com/labstack/echo/v4"
)

func dashboardRoutes(
	router *echo.Echo,
	handlers handlers.Dashboard,
) {
	dashboardRouter := router.Group("/dashboard")

	dashboardRouter.GET("", func(c echo.Context) error {
		return handlers.Home(c)
	}).Name = paths.DashboardHomePage
	dashboardRouter.GET("/newsletters", func(c echo.Context) error {
		return handlers.Newsletters(c)
	}).Name = paths.DashboardNewsletter
	dashboardRouter.GET("/newsletters/new", func(c echo.Context) error {
		return handlers.CreateNewsletters(c)
	}).Name = paths.DashboardNewsletterNew
}
