package router

import (
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v4"
)

func registerSubscriberRoutes(handler *echo.Echo, subscriber controllers.Subscribers) {
	handler.Add(
		http.MethodGet, routes.SubscriberIndex.Path(), subscriber.Index,
	).Name = routes.SubscriberIndex.Name()

	handler.Add(
		http.MethodGet, routes.SubscriberShow.Path(), subscriber.Show,
	).Name = routes.SubscriberShow.Name()

	handler.Add(
		http.MethodGet, routes.SubscriberNew.Path(), subscriber.New,
	).Name = routes.SubscriberNew.Name()

	handler.Add(
		http.MethodPost, routes.SubscriberCreate.Path(), subscriber.Create,
	).Name = routes.SubscriberCreate.Name()

	handler.Add(
		http.MethodGet, routes.SubscriberEdit.Path(), subscriber.Edit,
	).Name = routes.SubscriberEdit.Name()

	handler.Add(
		http.MethodPut, routes.SubscriberUpdate.Path(), subscriber.Update,
	).Name = routes.SubscriberUpdate.Name()

	handler.Add(
		http.MethodDelete, routes.SubscriberDestroy.Path(), subscriber.Destroy,
	).Name = routes.SubscriberDestroy.Name()
}
