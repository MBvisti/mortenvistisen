package routes

import (
	"mortenvistisen/internal/routing"
)

const AdminSubscriberPrefix = "/admin/subscribers"

var AdminSubscriberIndex = routing.NewSimpleRoute(
	"",
	"admin.subscribers.index",
	AdminSubscriberPrefix,
)
var AdminSubscriberBulkDestroy = routing.NewSimpleRoute(
	"/bulk-delete",
	"admin.subscribers.bulk-destroy",
	AdminSubscriberPrefix,
)
var AdminSubscriberConfirmationCreate = routing.NewRouteWithSerialID(
	"/:id/confirmation",
	"admin.subscribers.confirmation",
	AdminSubscriberPrefix,
)
var AdminSubscriberShow = routing.NewRouteWithSerialID(
	"/:id",
	"admin.subscribers.show",
	AdminSubscriberPrefix,
)
var AdminSubscriberNew = routing.NewSimpleRoute(
	"/new",
	"admin.subscribers.new",
	AdminSubscriberPrefix,
)
var AdminSubscriberCreate = routing.NewSimpleRoute(
	"",
	"admin.subscribers.create",
	AdminSubscriberPrefix,
)
var AdminSubscriberEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"admin.subscribers.edit",
	AdminSubscriberPrefix,
)
var AdminSubscriberUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"admin.subscribers.update",
	AdminSubscriberPrefix,
)
var AdminSubscriberDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"admin.subscribers.destroy",
	AdminSubscriberPrefix,
)
