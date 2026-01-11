package router

import (
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v4"
)

func registerAPIRoutes(handler *echo.Echo, apiController controllers.API) {
	handler.Add(
		http.MethodGet, routes.Health.Path(), apiController.Health,
	).Name = routes.Health.Name()
}
