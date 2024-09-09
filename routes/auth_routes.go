package routes

import "github.com/labstack/echo/v4"

func (r *Router) loadAuthRoutes() {
	router := r.router.Group("")
	router.GET("/register", func(c echo.Context) error {
		return r.registrationHandlers.CreateUser(c)
	})
	router.POST("/register", func(c echo.Context) error {
		return r.registrationHandlers.StoreUser(c)
	}, r.middleware.AdminOnly)

	router.GET("/login", func(c echo.Context) error {
		return r.authenticationHandlers.CreateSession(c)
	})
	router.POST("/login", func(c echo.Context) error {
		return r.authenticationHandlers.StoreAuthenticatedSession(c)
	})

	router.GET("/verify-email", func(c echo.Context) error {
		return r.registrationHandlers.UserEmailVerification(c)
	})

	router.GET("/forgot-password", func(c echo.Context) error {
		return r.authenticationHandlers.CreatePasswordReset(c)
	})
	router.POST("/forgot-password", func(c echo.Context) error {
		return r.authenticationHandlers.StorePasswordReset(c)
	})
	router.GET("/reset-password", func(c echo.Context) error {
		return r.authenticationHandlers.CreateResetPassword(c)
	})
	router.POST("/reset-password", func(c echo.Context) error {
		return r.authenticationHandlers.StoreResetPassword(c)
	})
}
