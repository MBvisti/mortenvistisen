package routes

import "mortenvistisen/internal/routing"

const SubscriberPrefix = "/subscriptions"

var SubscriberCreate = routing.NewSimpleRoute(
	"",
	"subscriptions.create",
	SubscriberPrefix,
)

var SubscriberConfirmationNew = routing.NewSimpleRoute(
	"/confirmation",
	"subscriptions.new_confirmation",
	SubscriberPrefix,
)

var SubscriberConfirmationCreate = routing.NewSimpleRoute(
	"/confirmation",
	"subscriptions.confirmation",
	SubscriberPrefix,
)
