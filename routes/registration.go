package routes

import (
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/routes/paths"
	"github.com/labstack/echo/v4"
)

func registrationRoutes(
	router *echo.Echo,
	handlers handlers.Registration,
) {
	router.GET("/register", func(c echo.Context) error {
		return handlers.NewUser(c)
	}).Name = paths.NewUser.String()
	router.POST("/register", func(c echo.Context) error {
		return handlers.CreateUser(c)
	}).Name = paths.CreateUser.String()

	router.GET("/verify-email", func(c echo.Context) error {
		return handlers.VerifyUserEmail(c)
	}).Name = paths.VerifyEmail.String()
}
