package web

import (
	"os"

	"github.com/MBvisti/mortenvistisen/routes/middleware"
	"github.com/labstack/echo/v4"
)

func (w *Web) DashboardRoutes() {
	adminGroup := w.router.Group("/dashboard")
	if os.Getenv("ENVIRONMENT") != "development" {
		adminGroup.Use(middleware.AuthOnly)
	}

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
	adminGroup.GET("/article/:slug/details", func(c echo.Context) error {
		return w.controllers.DashboardArticleDetails(c)
	})
	adminGroup.POST("/article/:slug/notify-subscribers", func(c echo.Context) error {
		return w.controllers.DashboardNotifySubscribers(c)
	})
}
