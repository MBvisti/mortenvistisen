package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/MBvisti/mortenvistisen/config"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/MBvisti/mortenvistisen/routes/paths"
	"github.com/dromara/carbon/v2"
	"github.com/labstack/echo/v4"
)

type Resource struct {
	db psql.Postgres
}

func newResource(db psql.Postgres) Resource {
	return Resource{db}
}

func (r Resource) Sitemap(c echo.Context) error {
	sitemap, err := createSitemap(c, r.db)
	if err != nil {
		return err
	}

	return c.XML(http.StatusOK, sitemap)
}

type URL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	ChangeFreq string   `xml:"changefreq"`
	LastMod    string   `xml:"lastmod"`
	Priority   string   `xml:"priority"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URL     []URL    `xml:"url"`
}

func createSitemap(c echo.Context, db psql.Postgres) (Sitemap, error) {
	baseUrl := config.Cfg.GetFullDomain()

	articles, err := models.GetAllArticles(context.Background(), db.Pool)
	if err != nil {
		return Sitemap{}, err
	}

	newsletters, err := models.GetAllNewsletters(context.Background(), db.Pool)
	if err != nil {
		return Sitemap{}, err
	}

	var urls []URL

	urls = append(urls, URL{
		Loc:        baseUrl,
		ChangeFreq: "monthly",
		LastMod:    "2024-10-22T09:43:09+00:00",
		Priority:   "1",
	})

	routes := c.Echo().Routes()
	for _, r := range routes {
		switch r.Name {
		case paths.About.String(),
			paths.Articles.String(),
			paths.Projects.String(),
			paths.Newsletters.String():
			urls = append(urls, URL{
				Loc: fmt.Sprintf(
					"%s%s",
					baseUrl,
					r.Path,
				),
				ChangeFreq: "monthly",
			})
		}
	}

	for _, article := range articles {
		urls = append(urls, URL{
			Loc: fmt.Sprintf(
				"%v/posts/%v",
				baseUrl,
				article.Slug,
			),
			ChangeFreq: "monthly",
			LastMod: carbon.CreateFromStdTime(article.UpdatedAt).
				ToDateString(),
			Priority: "0.9",
		})
	}

	for _, newsletter := range newsletters {
		urls = append(urls, URL{
			Loc: fmt.Sprintf(
				"%v/newsletters/%v",
				baseUrl,
				newsletter.Slug,
			),
			ChangeFreq: "monthly",
			LastMod: carbon.CreateFromStdTime(newsletter.UpdatedAt).
				ToDateString(),
			Priority: "0.9",
		})
	}

	sitemap := Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URL:   urls,
	}

	return sitemap, nil
}
