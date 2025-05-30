package routes

import (
	"net/http"
)

const (
	dashboardRoutePrefix = "/dashboard"
	dashboardNamePrefix  = "dashboard"
)

var Dashboard = []Route{
	DashboardHome,
	// DashboardArticles,
	// DashboardAnalytics,
	// DashboardSettings,
	DashboardNewArticle,
	DashboardStoreArticle,
	DashboardEditArticle,
	DashboardUpdateArticle,
	DashboardDeleteArticle,
	DashboardTags,
	DashboardCreateTag,
	DashboardUpdateTag,
	DashboardDeleteTag,
	DashboardSubscribers,
	DashboardEditSubscriber,
	DashboardUpdateSubscriber,
	DashboardDeleteSubscriber,
}

var DashboardHome = Route{
	Name:        dashboardNamePrefix + ".home",
	Path:        dashboardRoutePrefix,
	Method:      http.MethodGet,
	HandlerName: "Index",
	Middleware: []string{
		"AuthOnly",
	},
}

// var DashboardArticles = Route{
// 	Name:        dashboardNamePrefix + ".articles",
// 	Path:        dashboardRoutePrefix + "/articles",
// 	Method:      http.MethodGet,
// 	HandlerName: "Articles",
// 	Middleware: []string{
// 		"AuthOnly",
// 	},
// }
//
// var DashboardAnalytics = Route{
// 	Name:        dashboardNamePrefix + ".analytics",
// 	Path:        dashboardRoutePrefix + "/analytics",
// 	Method:      http.MethodGet,
// 	HandlerName: "Analytics",
// 	Middleware: []string{
// 		"AuthOnly",
// 	},
// }
//
// var DashboardSettings = Route{
// 	Name:        dashboardNamePrefix + ".settings",
// 	Path:        dashboardRoutePrefix + "/settings",
// 	Method:      http.MethodGet,
// 	HandlerName: "Settings",
// 	Middleware: []string{
// 		"AuthOnly",
// 	},
// }

var DashboardNewArticle = Route{
	Name:        dashboardNamePrefix + ".articles.new",
	Path:        dashboardRoutePrefix + "/articles/new",
	Method:      http.MethodGet,
	HandlerName: "NewArticle",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardStoreArticle = Route{
	Name:        dashboardNamePrefix + ".articles.create",
	Path:        dashboardRoutePrefix + "/articles/new",
	Method:      http.MethodPost,
	HandlerName: "StoreArticle",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardEditArticle = Route{
	Name:        dashboardNamePrefix + ".articles.edit",
	Path:        dashboardRoutePrefix + "/articles/:id/edit",
	Method:      http.MethodGet,
	HandlerName: "EditArticle",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardUpdateArticle = Route{
	Name:        dashboardNamePrefix + ".articles.update",
	Path:        dashboardRoutePrefix + "/articles/:id/edit",
	Method:      http.MethodPost,
	HandlerName: "UpdateArticle",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardDeleteArticle = Route{
	Name:        dashboardNamePrefix + ".articles.delete",
	Path:        dashboardRoutePrefix + "/articles/:id/delete",
	Method:      http.MethodPost,
	HandlerName: "DeleteArticle",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardTags = Route{
	Name:        dashboardNamePrefix + ".tags",
	Path:        dashboardRoutePrefix + "/tags",
	Method:      http.MethodGet,
	HandlerName: "Tags",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardCreateTag = Route{
	Name:        dashboardNamePrefix + ".tags.create",
	Path:        dashboardRoutePrefix + "/tags",
	Method:      http.MethodPost,
	HandlerName: "CreateTag",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardUpdateTag = Route{
	Name:        dashboardNamePrefix + ".tags.update",
	Path:        dashboardRoutePrefix + "/tags/:id/edit",
	Method:      http.MethodPost,
	HandlerName: "UpdateTag",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardDeleteTag = Route{
	Name:        dashboardNamePrefix + ".tags.delete",
	Path:        dashboardRoutePrefix + "/tags/:id/delete",
	Method:      http.MethodPost,
	HandlerName: "DeleteTag",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardSubscribers = Route{
	Name:        dashboardNamePrefix + ".subscribers",
	Path:        dashboardRoutePrefix + "/subscribers",
	Method:      http.MethodGet,
	HandlerName: "Subscribers",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardEditSubscriber = Route{
	Name:        dashboardNamePrefix + ".subscribers.edit",
	Path:        dashboardRoutePrefix + "/subscribers/:id/edit",
	Method:      http.MethodGet,
	HandlerName: "EditSubscriber",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardUpdateSubscriber = Route{
	Name:        dashboardNamePrefix + ".subscribers.update",
	Path:        dashboardRoutePrefix + "/subscribers/:id/edit",
	Method:      http.MethodPost,
	HandlerName: "UpdateSubscriber",
	Middleware: []string{
		"AuthOnly",
	},
}

var DashboardDeleteSubscriber = Route{
	Name:        dashboardNamePrefix + ".subscribers.delete",
	Path:        dashboardRoutePrefix + "/subscribers/:id/delete",
	Method:      http.MethodPost,
	HandlerName: "DeleteSubscriber",
	Middleware: []string{
		"AuthOnly",
	},
}
