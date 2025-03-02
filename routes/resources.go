package routes

import (
	"log/slog"
	"net/http"

	"github.com/MBvisti/mortenvistisen/config"
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/static"
	"github.com/labstack/echo/v4"
)

func resourceRoutes(router *echo.Echo, handlers handlers.Resource) {
	router.GET("/robots.txt", func(c echo.Context) error {
		return c.File("./resources/seo/robots.txt")
	})

	router.GET("/css/trix.css", func(c echo.Context) error {
		stylesheet, err := static.Files.ReadFile(
			"css/trix.css",
		)
		if err != nil {
			return err
		}

		if config.Cfg.Environment == config.PROD_ENVIRONMENT {
			c.Response().
				Header().
				Set("Cache-Control", "public, max-age=31536000, immutable")
			c.Response().
				Header().
				Set("Vary", "Accept-Encoding")
			c.Response().
				Header().
				Set("ETag", "\"bootstrap-v5_3_0\"")
		}

		return c.Blob(http.StatusOK, "text/css", stylesheet)
	})

	router.GET("/js/theme-switcher.js", func(c echo.Context) error {
		stylesheet, err := static.Files.ReadFile(
			"js/themeSwitcher.js",
		)
		if err != nil {
			return err
		}

		if config.Cfg.Environment == config.PROD_ENVIRONMENT {
			c.Response().
				Header().
				Set("Cache-Control", "public, max-age=31536000, immutable")
			c.Response().
				Header().
				Set("Vary", "Accept-Encoding")
			c.Response()
		}

		return c.Blob(http.StatusOK, "text/javascript", stylesheet)
	})

	router.GET("/css/bootstrap.css", func(c echo.Context) error {
		stylesheet, err := static.Files.ReadFile(
			"css/bootstrap-v5_3_0.css",
		)
		if err != nil {
			return err
		}

		if config.Cfg.Environment == config.PROD_ENVIRONMENT {
			c.Response().
				Header().
				Set("Cache-Control", "public, max-age=31536000, immutable")
			c.Response().
				Header().
				Set("Vary", "Accept-Encoding")
			c.Response().
				Header().
				Set("ETag", "\"bootstrap-v5_3_0\"")
		}

		return c.Blob(http.StatusOK, "text/css", stylesheet)
	})
	router.GET("/css/bootstrap-overrides.css", func(c echo.Context) error {
		stylesheet, err := static.Files.ReadFile(
			"css/bs-color-overrides.css",
		)
		if err != nil {
			return err
		}

		if config.Cfg.Environment == config.PROD_ENVIRONMENT {
			c.Response().
				Header().
				Set("Cache-Control", "public, max-age=2592000, immutable")
			c.Response().
				Header().
				Set("Vary", "Accept-Encoding")
		}

		return c.Blob(http.StatusOK, "text/css", stylesheet)
	})

	router.GET("/css/tailwind", func(c echo.Context) error {
		stylesheet, err := static.Files.ReadFile(
			"css/main-dev.css",
		)
		if err != nil {
			return err
		}

		if config.Cfg.Environment == config.PROD_ENVIRONMENT {
			c.Response().
				Header().
				Set("Cache-Control", "public, max-age=2592000, immutable")
			c.Response().
				Header().
				Set("Vary", "Accept-Encoding")
		}

		return c.Blob(http.StatusOK, "text/css", stylesheet)
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

		if config.Cfg.Environment == config.PROD_ENVIRONMENT {
			c.Response().
				Header().
				Set("Cache-Control", "public, max-age=2592000, immutable")
			c.Response().
				Header().
				Set("Vary", "Accept-Encoding")
		}

		return c.Blob(200, "application/javascript", bytes)
	})
}
