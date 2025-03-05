package routes

import (
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/routes/paths"
	"github.com/labstack/echo/v4"
)

func apiV1Routes(
	router *echo.Group,
	handlers handlers.Api,
) {
	router.GET("/health", func(c echo.Context) error {
		return handlers.AppHealth(c)
	}).Name = paths.APIHealth.String()

	router.POST("/collect", func(c echo.Context) error {
		return handlers.Collect(c)
	}).Name = paths.APICollect.String()
}
