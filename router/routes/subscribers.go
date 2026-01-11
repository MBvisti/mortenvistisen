package routes

import (
	"mortenvistisen/internal/routing"
)

const SubscriberPrefix = "/subscribers"

var SubscriberIndex = routing.NewSimpleRoute(
	"/",
	"subscribers.index",
	SubscriberPrefix,
)

var SubscriberShow = routing.NewRouteWithID(
	"/:id",
	"subscribers.show",
	SubscriberPrefix,
)

var SubscriberNew = routing.NewSimpleRoute(
	"/new",
	"subscribers.new",
	SubscriberPrefix,
)

var SubscriberCreate = routing.NewSimpleRoute(
	"/",
	"subscribers.create",
	SubscriberPrefix,
)

var SubscriberEdit = routing.NewRouteWithID(
	"/:id/edit",
	"subscribers.edit",
	SubscriberPrefix,
)

var SubscriberUpdate = routing.NewRouteWithID(
	"/:id",
	"subscribers.update",
	SubscriberPrefix,
)

var SubscriberDestroy = routing.NewRouteWithID(
	"/:id",
	"subscribers.destroy",
	SubscriberPrefix,
)
