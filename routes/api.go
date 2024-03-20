package routes

import (
	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/labstack/echo/v4"
)

func apiRoutes(router *echo.Group, controllers controllers.Controller) {
	router.GET("/health", func(c echo.Context) error {
		return controllers.AppHealth(c)
	})

}
