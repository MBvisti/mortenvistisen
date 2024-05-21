package api

import "github.com/labstack/echo/v4"

func LoadRoutes(router *echo.Echo) {
	r := router.Group("/v1")

	r.GET("/health", func(c echo.Context) error {
		return appHealth(c)
	})
}
