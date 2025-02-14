package routes

import (
	"log/slog"

	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/static"
	"github.com/labstack/echo/v4"
)

func resourceRoutes(router *echo.Echo, handlers handlers.Resource) {
	router.GET("/robots.txt", func(c echo.Context) error {
		return c.File("./resources/seo/robots.txt")
	})

	router.GET(
		"/4zd8j69sf3ju2hnfxmebr3czub8uu63m.txt",
		func(c echo.Context) error {
			indexNowTxt := []byte(`4zd8j69sf3ju2hnfxmebr3czub8uu63m`)
			return c.Blob(200, "text/plain", indexNowTxt)
		},
	)

	router.GET("/sitemap.xml", func(c echo.Context) error {
		return handlers.Sitemap(c)
	})

	router.GET("/favicon.ico", func(c echo.Context) error {
		return c.File("./static/images/favicon.ico")
	})

	router.GET("/script.js", func(c echo.Context) error {
		bytes, err := static.Files.ReadFile("js/analytics.js")
		if err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"ANALYTICS SCRIPT",
				"error",
				err,
			)
			return err
		}

		return c.Blob(200, "application/javascript", bytes)
	})
}
