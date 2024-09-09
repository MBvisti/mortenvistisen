package routes

import (
	"github.com/labstack/echo/v4"
)

func (r *Router) loadDashboardRoutes() {
	router := r.router.Group("/dashboard", r.middleware.AuthOnly)

	router.GET("", func(c echo.Context) error {
		return r.dashboardHandlers.Index(c)
	})

	router.GET("/subscribers", func(c echo.Context) error {
		return r.dashboardHandlers.SubscribersIndex(c)
	})

	router.POST(
		"/subscribers/:id/send-verification-mail",
		func(c echo.Context) error {
			return r.dashboardHandlers.ResendVerificationMail(c)
		},
	)

	router.GET("/articles", func(c echo.Context) error {
		return r.dashboardHandlers.ArticlesIndex(c)
	})
	router.GET("/articles/:slug/edit", func(c echo.Context) error {
		return r.dashboardHandlers.ArticleEdit(c)
	})
	router.GET("/articles/create", func(c echo.Context) error {
		return r.dashboardHandlers.ArticleCreate(c)
	})
	router.POST("/articles/store", func(c echo.Context) error {
		return r.dashboardHandlers.ArticleStore(c)
	})
	router.PUT("/articles/:id/update", func(c echo.Context) error {
		return r.dashboardHandlers.ArticleUpdate(c)
	})
	router.POST("/tags/store", func(c echo.Context) error {
		return r.dashboardHandlers.TagStore(c)
	})

	router.GET("/newsletters", func(c echo.Context) error {
		return r.dashboardHandlers.NewslettersIndex(c)
	})
	router.GET("/newsletters/create", func(c echo.Context) error {
		return r.dashboardHandlers.NewsletterCreate(c)
	})
	router.POST("/newsletters/store", func(c echo.Context) error {
		return r.dashboardHandlers.NewsletterStore(c)
	})
	// router.GET("/newsletters/:id/edit", func(c echo.Context) error {
	// 	return r.dashboardHandlers.NewslettersEdit(c)
	// })
	// router.PUT("/newsletters/:id/update", func(c echo.Context) error {
	// 	return r.dashboardHandlers.NewsletterUpdate(c)
	// })

	router.DELETE("/subscribers/:ID", func(c echo.Context) error {
		return r.dashboardHandlers.DeleteSubscriber(c)
	})
}
