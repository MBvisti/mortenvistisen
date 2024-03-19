package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/server/middleware"
)

func apiRoutes(router *echo.Group, controllers controllers.Controller, middleware middleware.Middleware) {
	router.GET("/health", func(c echo.Context) error {
		return controllers.AppHealth(c)
	})

}
