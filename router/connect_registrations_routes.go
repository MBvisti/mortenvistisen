package router

import (
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v4"
)

func registerRegistrationsRoutes(handler *echo.Echo, registrationsController controllers.Registrations) {
	handler.Add(
		http.MethodGet, routes.RegistrationNew.Path(), registrationsController.New,
	).Name = routes.RegistrationNew.Name()

	handler.Add(
		http.MethodPost, routes.RegistrationCreate.Path(), registrationsController.Create,
	).Name = routes.RegistrationCreate.Name()
}
