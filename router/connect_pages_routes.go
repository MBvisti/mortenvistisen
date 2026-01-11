package router

import (
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v4"
)

func registerPagesRoutes(handler *echo.Echo, pagesController controllers.Pages) {
	handler.Add(
		http.MethodGet, routes.HomePage.Path(), pagesController.Home,
	).Name = routes.HomePage.Name()
}
