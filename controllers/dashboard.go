package controllers

import (
	"net/http"

	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func (c *Controller) DashboardIndex(ctx echo.Context) error {
	return dashboard.Index().Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardSubscribers(ctx echo.Context) error {
	subs, err := c.db.QueryAllSubscribers(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.SubscriberViewData, 0, len(subs))
	for _, sub := range subs {
		viewData = append(viewData, dashboard.SubscriberViewData{
			Email:        sub.Email.String,
			ID:           sub.ID.String(),
			Verified:     sub.IsVerified.Bool,
			SubscribedAt: sub.SubscribedAt.Time.String(),
			Refererer:    sub.Referer.String,
		})
	}

	return dashboard.Subscribers(viewData, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardArticles(ctx echo.Context) error {
	articles, err := c.db.QueryAllPosts(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.ArticleViewData, 0, len(articles))
	for _, article := range articles {
		viewData = append(viewData, dashboard.ArticleViewData{
			ID:         article.ID.String(),
			Title:      article.Title,
			Draft:      article.Draft,
			ReleasedAt: article.ReleasedAt.Time.String(),
			Slug:       article.Slug,
		})
	}

	return dashboard.Articles(viewData, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardArticleEdit(ctx echo.Context) error {
	slug := ctx.Param("slug")

	post, err := c.db.GetPostBySlug(ctx.Request().Context(), slug)
	if err != nil {
		return err
	}

	postContent, err := c.postManager.Parse(post.Filename)
	if err != nil {
		return err
	}

	tags, err := c.db.GetTagsForPost(ctx.Request().Context(), post.ID)
	if err != nil {
		return err
	}

	var keywords string
	for i, kw := range tags {
		if i == len(tags)-1 {
			keywords = keywords + kw.Name
		} else {
			keywords = keywords + kw.Name + ", "
		}
	}

	return dashboard.ArticleEdit(dashboard.ArticleEditViewData{
		ID:          post.ID.String(),
		CreatedAt:   post.CreatedAt.Time,
		UpdatedAt:   post.UpdatedAt.Time,
		Title:       post.Title,
		HeaderTitle: post.HeaderTitle.String,
		Filename:    post.Filename,
		Slug:        post.Slug,
		Excerpt:     post.Excerpt,
		Draft:       post.Draft,
		ReleasedAt:  post.ReleasedAt.Time,
		ReadTime:    post.ReadTime.Int32,
		Content:     postContent,
		Keywords:    keywords,
	}, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DeleteSubscriber(ctx echo.Context) error {
	id := ctx.Param("ID")

	parsedID := uuid.MustParse(id)

	if err := c.db.DeleteSubscriber(ctx.Request().Context(), parsedID); err != nil {
		return err
	}

	ctx.Response().Header().Add("HX-Refresh", "true")

	return ctx.Redirect(http.StatusSeeOther, "/dashboard/subscribers")
}
