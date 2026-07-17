package controllers

import (
	"errors"
	"fmt"
	"mortenvistisen/internal/hypermedia"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/router"
	"mortenvistisen/router/routes"
	"mortenvistisen/views"
	"net/http"
	"slices"

	"github.com/labstack/echo/v5"
)

type Projects struct {
	db storage.Pool
}

func NewProjects(db storage.Pool) Projects {
	return Projects{db}
}

func (p Projects) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.ProjectIndex.Path(),
		Name:    routes.ProjectIndex.Name(),
		Handler: p.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.ProjectShow.Path(),
		Name:    routes.ProjectShow.Name(),
		Handler: p.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (p Projects) Index(etx *echo.Context) error {
	page, ok := parsePublicPage(etx.QueryParam("page"))
	if !ok {
		return publicNotFound(etx)
	}
	projects, err := models.Project.PaginatePublished(
		etx.Request().Context(),
		p.db.Executor(),
		page,
		12,
	)
	if err != nil {
		return fmt.Errorf("paginate published projects: %w", err)
	}
	if !publicPageExists(page, projects.TotalPages) {
		return publicNotFound(etx)
	}

	return hypermedia.RenderPage(etx, views.ProjectIndex{
		Items:       projects.Projects,
		CurrentPage: projects.Page,
		TotalPages:  projects.TotalPages,
	}.Page())
}

func (p Projects) Show(etx *echo.Context) error {
	project, err := models.Project.FindPublishedBySlug(
		etx.Request().Context(),
		p.db.Executor(),
		etx.Param("slug"),
	)
	if errors.Is(err, models.ErrNotFound) {
		return publicNotFound(etx)
	}
	if err != nil {
		return fmt.Errorf("find published project: %w", err)
	}

	rendered, headings, err := views.RenderMarkdown(project.Content)
	if err != nil {
		return err
	}

	related, err := models.Project.LatestPublished(etx.Request().Context(), p.db.Executor(), 4)
	if err != nil {
		return fmt.Errorf("load related projects: %w", err)
	}
	related = slices.DeleteFunc(
		related,
		func(item models.ProjectEntity) bool { return item.ID == project.ID },
	)
	related = related[:min(3, len(related))]

	return hypermedia.RenderPage(etx, views.ProjectShow{
		Item: project, Content: rendered, Headings: headings, Related: related,
	}.Page())
}
