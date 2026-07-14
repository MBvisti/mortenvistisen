package routes

import (
	"mortenvistisen/internal/routing"
)

const AdminNewsletterPrefix = "/admin/newsletters"

var AdminNewsletterIndex = routing.NewSimpleRoute(
	"",
	"admin.newsletters.index",
	AdminNewsletterPrefix,
)
var AdminNewsletterShow = routing.NewRouteWithSerialID(
	"/:id",
	"admin.newsletters.show",
	AdminNewsletterPrefix,
)
var AdminNewsletterNew = routing.NewSimpleRoute(
	"/new",
	"admin.newsletters.new",
	AdminNewsletterPrefix,
)
var AdminNewsletterCreate = routing.NewSimpleRoute(
	"",
	"admin.newsletters.create",
	AdminNewsletterPrefix,
)
var AdminNewsletterEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"admin.newsletters.edit",
	AdminNewsletterPrefix,
)
var AdminNewsletterUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"admin.newsletters.update",
	AdminNewsletterPrefix,
)
var AdminNewsletterDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"admin.newsletters.destroy",
	AdminNewsletterPrefix,
)
