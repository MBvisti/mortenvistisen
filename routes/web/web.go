package web

import (
	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/labstack/echo/v4"
)

type Web struct {
	controllers controllers.Controller
	router      *echo.Echo
}

func NewWeb(router *echo.Echo, controllers controllers.Controller) Web {
	return Web{
		controllers,
		router,
	}
}

func (w *Web) miscRoutes() {
	w.router.GET("/sitemap.xml", func(c echo.Context) error {
		return nil
	})
}

func (w *Web) SetupWebRoutes() {
	w.UtilityRoutes()

	w.ArticleRoutes()
	w.SiteRoutes()
	// w.UserRoutes()
	w.DashboardRoutes()
}
