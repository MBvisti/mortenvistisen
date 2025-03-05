package routes

import (
	"github.com/MBvisti/mortenvistisen/http"
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/routes/paths"
	"github.com/labstack/echo/v4"
)

func dashboardRoutes(
	router *echo.Echo,
	handlers handlers.Dashboard,
) {
	dashboardRouter := router.Group("/dashboard", http.AuthOnly)

	dashboardRouter.GET("", func(c echo.Context) error {
		return handlers.Home(c)
	}).Name = paths.Dashboard.String()

	dashboardRouter.GET("/subscribers/:id", func(c echo.Context) error {
		return handlers.ShowSubscriber(c)
	}).Name = paths.DashboardShowSubscriber.String()

	dashboardRouter.PUT("/subscribers/:id", func(c echo.Context) error {
		return handlers.UpdateSubscriber(c)
	}).Name = paths.DashboardUpdateSubscriber.String()

	dashboardRouter.DELETE("/subscribers/:id", func(c echo.Context) error {
		return handlers.DeleteSubscriber(c)
	}).Name = paths.DashboardDeleteSubscriber.String()

	dashboardRouter.GET("/newsletters", func(c echo.Context) error {
		return handlers.Newsletters(c)
	}).Name = paths.DashboardNewsletters.String()

	dashboardRouter.GET("/newsletters/new", func(c echo.Context) error {
		return handlers.NewNewsletters(c)
	}).Name = paths.DashboardNewNewsletter.String()

	dashboardRouter.POST("/newsletters", func(c echo.Context) error {
		return handlers.CreateNewsletter(c)
	}).Name = paths.DashboardCreateNewsletter.String()
}
