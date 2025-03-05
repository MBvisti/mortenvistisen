package routes

import (
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/routes/paths"
	"github.com/labstack/echo/v4"
)

func appRoutes(router *echo.Echo, handlers handlers.App) {
	router.GET("/", func(c echo.Context) error {
		return handlers.LandingPage(c)
	}).Name = paths.Home.String()

	router.GET("/about", func(c echo.Context) error {
		return handlers.AboutPage(c)
	}).Name = paths.About.String()

	router.GET("/posts/:postSlug", func(c echo.Context) error {
		return handlers.ArticlePage(c)
	}).Name = paths.Article.String()
	router.GET("/articles", func(c echo.Context) error {
		return handlers.ArticlesPage(c)
	}).Name = paths.Articles.String()

	router.GET("/projects", func(c echo.Context) error {
		return handlers.ProjectsPage(c)
	}).Name = paths.Projects.String()

	router.GET("/newsletters", func(c echo.Context) error {
		return handlers.NewslettersPage(c)
	}).Name = paths.Newsletters.String()

	router.GET("/newsletters/:slug", func(c echo.Context) error {
		return handlers.NewsletterPage(c)
	}).Name = paths.Newsletter.String()

	router.POST("/subscribers", func(c echo.Context) error {
		return handlers.CreateSubscription(c)
	}).Name = paths.CreateSubscription.String()

	router.GET("/verify-subscriber", func(c echo.Context) error {
		return handlers.SubscriberEmailVerification(c)
	}).Name = paths.VerifySubscriber.String()

	router.GET("/unsubscribe", func(c echo.Context) error {
		return handlers.DeleteSubscription(c)
	}).Name = paths.UnSubscribe.String()
}
