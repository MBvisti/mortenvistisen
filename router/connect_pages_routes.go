package router

import (
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v4"
)

func registerPagesRoutes(handler *echo.Echo, pages controllers.Pages) {
	handler.Add(
		http.MethodGet, routes.HomePage.Path(), pages.Home,
	).Name = routes.HomePage.Name()

	handler.Add(
		http.MethodGet, routes.ProjectsPage.Path(), pages.Projects,
	).Name = routes.ProjectsPage.Name()

	handler.Add(
		http.MethodGet, routes.AdminHome.Path(), pages.AdminHome,
	).Name = routes.ProjectsPage.Name()
}
