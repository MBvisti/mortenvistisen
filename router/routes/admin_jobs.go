package routes

import "mortenvistisen/internal/routing"

const AdminJobPrefix = "/admin/jobs"

var AdminJobIndex = routing.NewSimpleRoute(
	"",
	"admin.jobs.index",
	AdminJobPrefix,
)

var AdminJobShow = routing.NewRouteWithBigSerialID(
	"/:id",
	"admin.jobs.show",
	AdminJobPrefix,
)

var AdminJobRun = routing.NewRouteWithBigSerialID(
	"/:id/run",
	"admin.jobs.run",
	AdminJobPrefix,
)

var AdminJobRetry = routing.NewRouteWithBigSerialID(
	"/:id/retry",
	"admin.jobs.retry",
	AdminJobPrefix,
)

var AdminJobCancel = routing.NewRouteWithBigSerialID(
	"/:id/cancel",
	"admin.jobs.cancel",
	AdminJobPrefix,
)

var AdminJobDestroy = routing.NewRouteWithBigSerialID(
	"/:id",
	"admin.jobs.destroy",
	AdminJobPrefix,
)
