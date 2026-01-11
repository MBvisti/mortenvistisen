package router

import (
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v4"
)

func registerNewsletterRoutes(handler *echo.Echo, newsletter controllers.Newsletters) {
	handler.Add(
		http.MethodGet, routes.NewsletterIndex.Path(), newsletter.Index,
	).Name = routes.NewsletterIndex.Name()

	handler.Add(
		http.MethodGet, routes.NewsletterShow.Path(), newsletter.Show,
	).Name = routes.NewsletterShow.Name()

	handler.Add(
		http.MethodGet, routes.NewsletterNew.Path(), newsletter.New,
	).Name = routes.NewsletterNew.Name()

	handler.Add(
		http.MethodPost, routes.NewsletterCreate.Path(), newsletter.Create,
	).Name = routes.NewsletterCreate.Name()

	handler.Add(
		http.MethodGet, routes.NewsletterEdit.Path(), newsletter.Edit,
	).Name = routes.NewsletterEdit.Name()

	handler.Add(
		http.MethodPut, routes.NewsletterUpdate.Path(), newsletter.Update,
	).Name = routes.NewsletterUpdate.Name()

	handler.Add(
		http.MethodDelete, routes.NewsletterDestroy.Path(), newsletter.Destroy,
	).Name = routes.NewsletterDestroy.Name()
}
