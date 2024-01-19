package web

import "github.com/labstack/echo/v4"

func (w *Web) SubscribeRoutes() {
	w.router.POST("/subscribe", func(c echo.Context) error {
		return w.controllers.SubscriptionEvent(c)
	})
}
