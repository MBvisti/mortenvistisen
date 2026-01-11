package router

import (
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v4"
)

func registerArticleRoutes(handler *echo.Echo, article controllers.Articles) {
	handler.Add(
		http.MethodGet, routes.ArticleIndex.Path(), article.Index,
	).Name = routes.ArticleIndex.Name()

	handler.Add(
		http.MethodGet, routes.ArticleShow.Path(), article.Show,
	).Name = routes.ArticleShow.Name()

	handler.Add(
		http.MethodGet, routes.ArticleNew.Path(), article.New,
	).Name = routes.ArticleNew.Name()

	handler.Add(
		http.MethodPost, routes.ArticleCreate.Path(), article.Create,
	).Name = routes.ArticleCreate.Name()

	handler.Add(
		http.MethodGet, routes.ArticleEdit.Path(), article.Edit,
	).Name = routes.ArticleEdit.Name()

	handler.Add(
		http.MethodPut, routes.ArticleUpdate.Path(), article.Update,
	).Name = routes.ArticleUpdate.Name()

	handler.Add(
		http.MethodDelete, routes.ArticleDestroy.Path(), article.Destroy,
	).Name = routes.ArticleDestroy.Name()
}
