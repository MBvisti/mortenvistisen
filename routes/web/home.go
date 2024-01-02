package web

import (
	"github.com/labstack/echo/v4"
)

func (w *Web) SiteRoutes() {
	w.router.GET("/", func(c echo.Context) error {
		return w.controllers.HomeIndex(c)
	})

	w.router.GET("/about", func(c echo.Context) error {
		return w.controllers.About(c)
	})

	w.router.GET("/newsletter", func(c echo.Context) error {
		return w.controllers.Newsletter(c)
	})

	w.router.GET("/projects", func(c echo.Context) error {
		return w.controllers.Projects(c)
	})
}
