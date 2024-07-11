package routes

import "github.com/labstack/echo/v4"

func (r *Router) loadApiV1Routes() {
	router := r.router.Group("/api/v1")

	router.GET("/health", func(c echo.Context) error {
		return r.apiHandlers.Health(c)
	})
}
