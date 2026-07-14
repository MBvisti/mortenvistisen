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
	"strconv"

	"github.com/labstack/echo/v5"
)

type Posts struct {
	db storage.Pool
}

func NewPosts(db storage.Pool) Posts {
	return Posts{db}
}

func (p Posts) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.PostIndex.Path(),
		Name:    routes.PostIndex.Name(),
		Handler: p.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.PostShow.Path(),
		Name:    routes.PostShow.Name(),
		Handler: p.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func parsePublicPage(value string) (int64, bool) {
	if value == "" {
		return 1, true
	}
	page, err := strconv.ParseInt(value, 10, 64)
	return page, err == nil && page > 0 && value == strconv.FormatInt(page, 10)
}

func publicPageExists(page, totalPages int64) bool {
	return page == 1 || page <= totalPages
}

func (p Posts) Index(etx *echo.Context) error {
	page, ok := parsePublicPage(etx.QueryParam("page"))
	if !ok {
		return publicNotFound(etx)
	}
	articles, err := models.Article.PaginatePublished(
		etx.Request().Context(),
		p.db.Executor(),
		page,
		12,
	)
	if err != nil {
		return fmt.Errorf("paginate published posts: %w", err)
	}
	if !publicPageExists(page, articles.TotalPages) {
		return publicNotFound(etx)
	}

	return hypermedia.RenderPage(etx, views.PostIndex{
		Items:       articles.Articles,
		CurrentPage: articles.Page,
		TotalPages:  articles.TotalPages,
	}.Page())
}

func (p Posts) Show(etx *echo.Context) error {
	post, err := models.Article.FindPublishedBySlug(
		etx.Request().Context(),
		p.db.Executor(),
		etx.Param("slug"),
	)
	if errors.Is(err, models.ErrNotFound) {
		return publicNotFound(etx)
	}
	if err != nil {
		return fmt.Errorf("find published post: %w", err)
	}

	tags, err := models.ArticleTagConnection.TagsForArticle(
		etx.Request().Context(),
		p.db.Executor(),
		post.ID,
	)
	if err != nil {
		return fmt.Errorf("find post tags: %w", err)
	}

	content := ""
	if post.Content.Valid {
		content = post.Content.String
	}
	rendered, headings, err := views.RenderMarkdown(content)
	if err != nil {
		return err
	}

	related, err := models.Article.LatestPublished(etx.Request().Context(), p.db.Executor(), 4)
	if err != nil {
		return fmt.Errorf("load related posts: %w", err)
	}
	related = slices.DeleteFunc(related, func(item models.ArticleEntity) bool { return item.ID == post.ID })
	related = related[:min(3, len(related))]

	return hypermedia.RenderPage(etx, views.PostShow{
		Item: post, Content: rendered, Headings: headings, Tags: tags, Related: related,
	}.Page())
}
