package routes

import (
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/views/paths"
	"github.com/labstack/echo/v4"
)

func registrationRoutes(
	router *echo.Echo,
	handlers handlers.Registration,
) {
	router.GET("/register", func(c echo.Context) error {
		return handlers.CreateUser(c)
	}).Name = paths.RegisterPage.ToString()
	router.POST("/register", func(c echo.Context) error {
		return handlers.StoreUser(c)
	}).Name = paths.RegisterUser.ToString()

	router.GET("/verify-email", func(c echo.Context) error {
		return handlers.VerifyUserEmail(c)
	}).Name = paths.VerifyEmailPage.ToString()
}
