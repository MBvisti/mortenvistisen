package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (r *Router) loadAppRoutes() {
	router := r.router.Group("")

	router.GET("/", func(c echo.Context) error {
		return r.appHandlers.Index(c)
	})
	router.GET("", func(c echo.Context) error {
		return r.appHandlers.Index(c)
	})

	router.GET("/about", func(c echo.Context) error {
		return r.appHandlers.About(c)
	})

	router.GET("/newsletter", func(c echo.Context) error {
		return r.appHandlers.Newsletter(c)
	})

	router.GET("/projects", func(c echo.Context) error {
		return r.appHandlers.Projects(c)
	})

	router.GET("/posts/:postSlug", func(c echo.Context) error {
		return r.appHandlers.Article(c)
	})

	router.GET("/modal", func(c echo.Context) error {
		return r.appHandlers.RenderModal(c)
	})

	router.GET("/verify-subscriber", func(c echo.Context) error {
		return r.appHandlers.SubscriberEmailVerification(c)
	})

	router.POST("/subscribe", func(c echo.Context) error {
		return r.appHandlers.SubscriptionEvent(c)
	})

	router.GET("/redirect", func(c echo.Context) error {
		to := c.QueryParam("to")
		return r.baseHandlers.RedirectHx(c.Response(), fmt.Sprintf("/%s", to))
	})
}
