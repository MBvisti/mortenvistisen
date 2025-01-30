package routes

import (
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/views/paths"
	"github.com/labstack/echo/v4"
)

func appRoutes(router *echo.Echo, handlers handlers.App) {
	router.GET("/", func(c echo.Context) error {
		return handlers.LandingPage(c)
	}).Name = paths.HomePage

	router.GET("/about", func(c echo.Context) error {
		return handlers.AboutPage(c)
	}).Name = paths.AboutPage

	router.GET("/posts/:postSlug", func(c echo.Context) error {
		return handlers.ArticlePage(c)
	}).Name = paths.ArticlePage
	router.GET("/articles", func(c echo.Context) error {
		return handlers.ArticlesPage(c)
	}).Name = paths.ArticlesPage

	router.GET("/projects", func(c echo.Context) error {
		return handlers.ProjectsPage(c)
	}).Name = paths.ProjectsPage

	router.GET("/newsletters", func(c echo.Context) error {
		return handlers.NewslettersPage(c)
	}).Name = paths.NewslettersPage

	router.POST("/subscribe", func(c echo.Context) error {
		return handlers.SubscriptionEvent(c)
	}).Name = paths.SubscribeEvent

	router.GET("/verify-subscriber", func(c echo.Context) error {
		return handlers.SubscriberEmailVerification(c)
	}).Name = paths.VerifySubEvent

	router.GET("/unsubscribe", func(c echo.Context) error {
		return handlers.UnsubscriptionEvent(c)
	}).Name = paths.UnsubscribeEvent
}
