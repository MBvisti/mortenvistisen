package admin

import (
	"errors"
	"mortenvistisen/internal/inertia"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/router"
	"mortenvistisen/router/routes"
	"net/http"

	"github.com/labstack/echo/v5"
)

const dashboardRecentLimit int64 = 5

type Dashboards struct {
	db storage.Pool
}

func NewDashboards(db storage.Pool) Dashboards {
	return Dashboards{db}
}

func (d Dashboards) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminDashboardIndex.Path(),
		Name:    routes.AdminDashboardIndex.Name(),
		Handler: d.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type DashboardStats struct {
	Subscribers int64
	Articles    int64
	Newsletters int64
}

func (d Dashboards) Index(etx *echo.Context) error {
	ctx := etx.Request().Context()
	db := d.db.Executor()

	subscribersList, err := models.Subscriber.Paginate(
		ctx,
		db,
		1,
		dashboardRecentLimit,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	articlesList, err := models.Article.Paginate(
		ctx,
		db,
		models.ArticleFilter{},
		1,
		dashboardRecentLimit,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	newslettersList, err := models.Newsletter.Paginate(
		ctx,
		db,
		1,
		dashboardRecentLimit,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Dashboard/Index", inertia.Props{
		"stats": DashboardStats{
			Subscribers: subscribersList.TotalCount,
			Articles:    articlesList.TotalCount,
			Newsletters: newslettersList.TotalCount,
		},
		"subscribers": newSubscriberDataList(subscribersList.Subscribers),
		"articles":    newArticleDataList(articlesList.Articles),
		"newsletters": newNewsletterDataList(newslettersList.Newsletters),
	})
}
