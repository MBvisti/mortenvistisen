package controllers

import (
	"context"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/router/routes"
	"mortenvistisen/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Pages struct {
	db         storage.Pool
	insertOnly queue.InsertOnly
	cache      *Cache[templ.Component]
}

func NewPages(
	db storage.Pool,
	insertOnly queue.InsertOnly,
	cache *Cache[templ.Component],
) Pages {
	return Pages{db, insertOnly, cache}
}

func (p Pages) Home(etx echo.Context) error {
	cacheKey := "home"

	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		articles, err := models.FindPublishedArticles(context.Background(), p.db.Conn())
		if err != nil {
			return nil, err
		}

		articleViews := make([]views.ArticleViewData, len(articles))
		for i, article := range articles {
			articleViews[i] = views.ArticleViewData{
				PublishedAt: article.FirstPublishedAt,
				Title:       article.Title,
				Excerpt:     article.Excerpt,
				URL:         routes.ArticleShowSlug.URL(article.Slug),
				Tags:        article.Tags,
			}
		}

		return views.Home(articleViews), nil
	})
	if err != nil {
		return err
	}

	return render(etx, component)
}

func (p Pages) NotFound(etx echo.Context) error {
	cacheKey := "not_found"

	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		return views.NotFound(), nil
	})
	if err != nil {
		return err
	}

	return render(etx, component)
}

func (p Pages) Projects(etx echo.Context) error {
	cacheKey := "projects"

	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		return views.Projects(), nil
	})
	if err != nil {
		return err
	}

	return render(etx, component)
}
