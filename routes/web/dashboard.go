package web

import (
	"github.com/labstack/echo/v4"
)

func (w *Web) DashboardRoutes() {
	w.router.GET("/dashboard", func(c echo.Context) error {
		return w.controllers.DashboardIndex(c)
	})
	w.router.GET("/dashboard/subscribers", func(c echo.Context) error {
		return w.controllers.DashboardSubscribers(c)
	})
	w.router.DELETE("/dashboard/subscriber/:ID", func(c echo.Context) error {
		return w.controllers.DeleteSubscriber(c)
	})
}
