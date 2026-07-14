package routes

import (
	"mortenvistisen/internal/routing"
)

const PostPrefix = "/posts"

var PostIndex = routing.NewSimpleRoute(
	"",
	"posts.index",
	PostPrefix,
)
var PostShow = routing.NewRouteWithSlug(
	"/:slug",
	"posts.show",
	PostPrefix,
)
