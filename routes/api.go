package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/MBvisti/mortenvistisen/controllers"
)

func apiRoutes(router *echo.Group, controllers controllers.Controller) {
	router.GET("/health", func(c echo.Context) error {
		return controllers.AppHealth(c)
	})

}
