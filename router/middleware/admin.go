package middleware

import (
	"net/http"

	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		app := cookies.GetApp(c)
		if app.IsAuthenticated && app.IsAdmin {
			return next(c)
		}

		return c.Redirect(http.StatusSeeOther, routes.HomePage.URL())
	}
}
