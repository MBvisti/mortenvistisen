package routes

import (
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func fragmentRoutes(router *echo.Echo) {
	fragments := router.Group("/fragments")
	fragments.GET("/load-csrf", func(c echo.Context) error {
		return views.CsrfToken(csrf.Token(c.Request())).
			Render(c.Request().Context(), c.Response())
	})
}
