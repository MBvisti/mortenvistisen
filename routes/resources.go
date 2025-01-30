package routes

import (
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/labstack/echo/v4"
)

func resourceRoutes(router *echo.Echo, handlers handlers.Resource) {
	router.GET("/robots.txt", func(c echo.Context) error {
		return c.File("./resources/seo/robots.txt")
	})
	router.GET("/sitemap.xml", func(c echo.Context) error {
		return handlers.Sitemap(c)
	})

	router.GET("/favicon.ico", func(c echo.Context) error {
		return c.File("./static/images/favicon.ico")
	})
}
