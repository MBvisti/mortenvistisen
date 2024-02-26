package web

import (
	"github.com/labstack/echo/v4"
)

func (w *Web) ArticleRoutes() {
	w.router.GET("/posts/:postSlug", func(c echo.Context) error {
		return w.controllers.Article(c)
	})
	w.router.GET("/modal", func(c echo.Context) error {
		return w.controllers.RenderModal(c)
	})
}
