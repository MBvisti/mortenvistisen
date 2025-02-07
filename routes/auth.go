package routes

import (
	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/views/paths"
	"github.com/labstack/echo/v4"
)

func authRoutes(
	router *echo.Echo,
	handlers handlers.Authentication,
) {
	router.GET("/login", func(c echo.Context) error {
		return handlers.CreateAuthenticatedSession(c)
	}).Name = paths.LoginPage.ToString()
	router.POST("/login", func(c echo.Context) error {
		return handlers.StoreAuthenticatedSession(c)
	}).Name = paths.Login.ToString()

	router.GET("/forgot-password", func(c echo.Context) error {
		return handlers.CreatePasswordReset(c)
	}).Name = paths.ForgotPasswordPage.ToString()
	router.POST("/forgot-password", func(c echo.Context) error {
		return handlers.StorePasswordReset(c)
	}).Name = paths.ForgotPassword.ToString()

	router.GET("/reset-password", func(c echo.Context) error {
		return handlers.CreateResetPassword(c)
	}).Name = paths.ResetPasswordPage.ToString()
	router.POST("/reset-password", func(c echo.Context) error {
		return handlers.StoreResetPassword(c)
	}).Name = paths.ResetPassword.ToString()
}
