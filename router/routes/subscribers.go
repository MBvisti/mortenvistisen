package routes

import (
	"mortenvistisen/internal/routing"
)

const SubscriberPrefix = "/subscribers"

var SubscriberIndex = routing.NewSimpleRoute(
	"",
	"subscribers.index",
	AdminPrefix+SubscriberPrefix,
)

var SubscriberShow = routing.NewRouteWithSerialID(
	"/:id",
	"subscribers.show",
	AdminPrefix+SubscriberPrefix,
)

var SubscriberNew = routing.NewSimpleRoute(
	"/new",
	"subscribers.new",
	AdminPrefix+SubscriberPrefix,
)

var SubscriberCreate = routing.NewSimpleRoute(
	"",
	"subscribers.create",
	AdminPrefix+SubscriberPrefix,
)

var SubscriberSignup = routing.NewSimpleRoute(
	"/signup",
	"subscribers.signup",
	SubscriberPrefix,
)

var SubscriberVerificationNew = routing.NewSimpleRoute(
	"/verify",
	"subscribers.verify.new",
	SubscriberPrefix,
)

var SubscriberVerificationCreate = routing.NewSimpleRoute(
	"/verify",
	"subscribers.verify.create",
	SubscriberPrefix,
)

var SubscriberUnsubscribe = routing.NewSimpleRoute(
	"/unsubscribe",
	"subscribers.unsubscribe",
	"",
)

var SubscriberEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"subscribers.edit",
	AdminPrefix+SubscriberPrefix,
)

var SubscriberUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"subscribers.update",
	AdminPrefix+SubscriberPrefix,
)

var SubscriberDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"subscribers.destroy",
	AdminPrefix+SubscriberPrefix,
)
