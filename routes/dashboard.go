package routes

import (
	"github.com/MBvisti/mortenvistisen/http"
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/views/paths"
	"github.com/labstack/echo/v4"
)

func dashboardRoutes(
	router *echo.Echo,
	handlers handlers.Dashboard,
) {
	dashboardRouter := router.Group("/dashboard", http.AuthOnly)

	dashboardRouter.GET("", func(c echo.Context) error {
		return handlers.Home(c)
	}).Name = paths.DashboardHomePage.ToString()

	dashboardRouter.GET("/subscribers/:id", func(c echo.Context) error {
		return handlers.ShowSubscriber(c)
	}).Name = paths.DashboardSubscriberPage.ToString()

	dashboardRouter.GET("/newsletters", func(c echo.Context) error {
		return handlers.Newsletters(c)
	}).Name = paths.DashboardNewsletter.ToString()
	dashboardRouter.GET("/newsletters/new", func(c echo.Context) error {
		return handlers.CreateNewsletters(c)
	}).Name = paths.DashboardNewsletterNew.ToString()
	dashboardRouter.POST("/newsletters/new", func(c echo.Context) error {
		return handlers.StoreNewsletter(c)
	}).Name = paths.DashboardNewsletterStore.ToString()
}
