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

type Newsletters struct {
	db storage.Pool
}

func NewNewsletters(db storage.Pool) Newsletters {
	return Newsletters{db}
}

func (n Newsletters) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.NewsletterIndex.Path(),
		Name:    routes.NewsletterIndex.Name(),
		Handler: n.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.NewsletterShow.Path(),
		Name:    routes.NewsletterShow.Name(),
		Handler: n.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (n Newsletters) Index(etx *echo.Context) error {
	page, ok := parsePublicPage(etx.QueryParam("page"))
	if !ok {
		return publicNotFound(etx)
	}
	newsletters, err := models.Newsletter.PaginatePublished(
		etx.Request().Context(),
		n.db.Executor(),
		page,
		12,
	)
	if err != nil {
		return fmt.Errorf("paginate published newsletters: %w", err)
	}
	if !publicPageExists(page, newsletters.TotalPages) {
		return publicNotFound(etx)
	}

	return hypermedia.RenderPage(etx, views.NewsletterIndex{
		Items:       newsletters.Newsletters,
		CurrentPage: newsletters.Page,
		TotalPages:  newsletters.TotalPages,
	}.Page())
}

func (n Newsletters) Show(etx *echo.Context) error {
	newsletter, err := models.Newsletter.FindPublishedBySlug(
		etx.Request().Context(),
		n.db.Executor(),
		etx.Param("slug"),
	)
	if errors.Is(err, models.ErrNotFound) {
		return publicNotFound(etx)
	}
	if err != nil {
		return fmt.Errorf("find published newsletter: %w", err)
	}

	content := ""
	if newsletter.Content.Valid {
		content = newsletter.Content.String
	}
	rendered, headings, err := views.RenderMarkdown(content)
	if err != nil {
		return err
	}

	related, err := models.Newsletter.LatestPublished(etx.Request().Context(), n.db.Executor(), 4)
	if err != nil {
		return fmt.Errorf("load related newsletters: %w", err)
	}
	related = slices.DeleteFunc(
		related,
		func(item models.NewsletterEntity) bool { return item.ID == newsletter.ID },
	)
	related = related[:min(3, len(related))]

	return hypermedia.RenderPage(etx, views.NewsletterShow{
		Item: newsletter, Content: rendered, Headings: headings, Related: related,
	}.Page())
}
