package web

import (
	"github.com/labstack/echo/v4"
)

func (w *Web) DashboardRoutes() {
	adminGroup := w.router.Group("/dashboard")

	adminGroup.GET("", func(c echo.Context) error {
		return w.controllers.DashboardIndex(c)
	})
	adminGroup.GET("/subscribers", func(c echo.Context) error {
		return w.controllers.DashboardSubscribers(c)
	})
	adminGroup.DELETE("/subscriber/:ID", func(c echo.Context) error {
		return w.controllers.DeleteSubscriber(c)
	})
	adminGroup.GET("/articles", func(c echo.Context) error {
		return w.controllers.DashboardArticles(c)
	})
}
