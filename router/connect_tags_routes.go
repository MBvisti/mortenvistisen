package router

import (
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v4"
)

func registerTagRoutes(handler *echo.Echo, tag controllers.Tags) {
	handler.Add(
		http.MethodGet, routes.TagIndex.Path(), tag.Index,
	).Name = routes.TagIndex.Name()

	handler.Add(
		http.MethodGet, routes.TagShow.Path(), tag.Show,
	).Name = routes.TagShow.Name()

	handler.Add(
		http.MethodGet, routes.TagNew.Path(), tag.New,
	).Name = routes.TagNew.Name()

	handler.Add(
		http.MethodPost, routes.TagCreate.Path(), tag.Create,
	).Name = routes.TagCreate.Name()

	handler.Add(
		http.MethodGet, routes.TagEdit.Path(), tag.Edit,
	).Name = routes.TagEdit.Name()

	handler.Add(
		http.MethodPut, routes.TagUpdate.Path(), tag.Update,
	).Name = routes.TagUpdate.Name()

	handler.Add(
		http.MethodDelete, routes.TagDestroy.Path(), tag.Destroy,
	).Name = routes.TagDestroy.Name()
}
