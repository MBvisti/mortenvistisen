package routes

import "mortenvistisen/internal/routing"

const ApiArticlePrefix = "/api/articles"

var ApiArticleUpdate = routing.NewRouteWithSlug(
	"/:slug",
	"api.articles.update",
	ApiArticlePrefix,
)
