package controllers

import (
	"net/http"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
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

func (p Pages) Home(etx *echo.Context) error {
	cacheKey := "pages:home"

	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		articles, err := models.AllPublishedArticles(etx.Request().Context(), p.db.Conn())
		if err != nil {
			return views.Home(nil), err
		}

		return views.Home(articles), nil
	})
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, component)
}

func (p Pages) About(etx *echo.Context) error {
	return render(etx, views.About())
}

func (p Pages) Article(etx *echo.Context) error {
	slug := etx.Param("slug")
	cacheKey := "pages:article:" + slug

	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		article, err := models.FindArticleBySlug(etx.Request().Context(), p.db.Conn(), slug)
		if err != nil {
			return views.Article(models.Article{}), err
		}

		return views.Article(article), nil
	})
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, component)
}

func (p Pages) ArticlesOverview(etx *echo.Context) error {
	cacheKey := "pages:articles_overview"

	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		articles, err := models.AllPublishedArticles(etx.Request().Context(), p.db.Conn())
		if err != nil {
			return views.ArticlesOverview(nil), err
		}

		published := make([]models.Article, 0, len(articles))
		for _, article := range articles {
			if article.Published {
				published = append(published, article)
			}
		}

		return views.ArticlesOverview(published), nil
	})
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, component)
}

func (p Pages) Project(etx *echo.Context) error {
	slug := etx.Param("slug")
	cacheKey := "pages:project:" + slug

	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		project, err := models.FindProjectBySlug(etx.Request().Context(), p.db.Conn(), slug)
		if err != nil {
			return views.Project(models.Project{}), err
		}

		return views.Project(project), nil
	})
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, component)
}

func (p Pages) ProjectsOverview(etx *echo.Context) error {
	cacheKey := "pages:projects_overview"

	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		projects, err := models.AllPublishedProjects(etx.Request().Context(), p.db.Conn())
		if err != nil {
			return views.ProjectsOverview(nil), err
		}

		published := make([]models.Project, 0, len(projects))
		for _, project := range projects {
			if project.Published {
				published = append(published, project)
			}
		}

		return views.ProjectsOverview(published), nil
	})
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, component)
}

func (p Pages) Newsletter(etx *echo.Context) error {
	slug := etx.Param("slug")
	cacheKey := "pages:newsletter:" + slug

	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		newsletter, err := models.FindNewsletterBySlug(etx.Request().Context(), p.db.Conn(), slug)
		if err != nil {
			return views.Newsletter(models.Newsletter{}), err
		}

		return views.Newsletter(newsletter), nil
	})
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, component)
}

func (p Pages) NewslettersOverview(etx *echo.Context) error {
	cacheKey := "pages:newsletters_overview"
	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		newsletters, err := models.AllPublishedNewsletters(etx.Request().Context(), p.db.Conn())
		if err != nil {
			return views.NewslettersOverview(nil), err
		}

		published := make([]models.Newsletter, 0, len(newsletters))
		for _, newsletter := range newsletters {
			if newsletter.IsPublished {
				published = append(published, newsletter)
			}
		}

		return views.NewslettersOverview(published), nil
	})
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, component)
}

func (p Pages) AdminHome(etx *echo.Context) error {
	ctx := etx.Request().Context()
	conn := p.db.Conn()

	articleCount, err := models.CountArticles(ctx, conn)
	if err != nil {
		return err
	}

	newsletterCount, err := models.CountNewsletters(ctx, conn)
	if err != nil {
		return err
	}

	projectCount, err := models.CountProjects(ctx, conn)
	if err != nil {
		return err
	}

	subscriberCount, err := models.CountSubscribers(ctx, conn)
	if err != nil {
		return err
	}

	return render(
		etx,
		views.AdminHome(articleCount, newsletterCount, projectCount, subscriberCount),
	)
}

func (p Pages) NotFound(etx *echo.Context) error {
	cacheKey := "not_found"

	component, err := p.cache.Get(cacheKey, func() (templ.Component, error) {
		return views.NotFound(), nil
	})
	if err != nil {
		return err
	}

	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := component.Render(etx.Request().Context(), buf); err != nil {
		return err
	}

	return etx.HTML(http.StatusNotFound, buf.String())
}
