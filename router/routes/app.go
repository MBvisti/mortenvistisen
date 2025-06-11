package routes

import (
	"fmt"
	"net/http"
)

const appNamePrefix = "app"

var App = []Route{
	LandingPage,
	AboutPage,
	ProjectsPage,
	ArticlePage,
	NewslettersPage,
	NewsletterPage,
	Redirect.Route,
}

var LandingPage = Route{
	Name:        appNamePrefix + ".landing_page",
	Path:        "/",
	Method:      http.MethodGet,
	HandlerName: "LandingPage",
}

var AboutPage = Route{
	Name:        appNamePrefix + ".about_page",
	Path:        "/about",
	Method:      http.MethodGet,
	HandlerName: "AboutPage",
}

var ProjectsPage = Route{
	Name:        appNamePrefix + ".projects_page",
	Path:        "/projects",
	Method:      http.MethodGet,
	HandlerName: "ProjectsPage",
}

var NewslettersPage = Route{
	Name:        appNamePrefix + ".newsletters_page",
	Path:        "/newsletters",
	Method:      http.MethodGet,
	HandlerName: "NewslettersPage",
}

var ArticlePage = Route{
	Name:        appNamePrefix + ".article_page",
	Path:        "/posts/:articleSlug",
	Method:      http.MethodGet,
	HandlerName: "ArticlePage",
}

var NewsletterPage = Route{
	Name:        appNamePrefix + ".newsletter_page",
	Path:        "/newsletters/:newsletterSlug",
	Method:      http.MethodGet,
	HandlerName: "NewsletterPage",
}

var Redirect = redirect{
	Route: Route{
		Name:        appNamePrefix + ".redirect",
		Path:        "/redirect",
		HandlerName: "Redirect",
		Method:      http.MethodGet,
	},
}

type redirect struct {
	Route
}

func (r redirect) WithQuery(route Route) string {
	return fmt.Sprintf("%s?to=%s", r.Path, route.Path)
}
