package handlers

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func (d Dashboard) ArticlesIndex(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	pageLimit := 7

	offset, currentPage, err := d.base.GetOffsetAndCurrPage(page, pageLimit)
	if err != nil {
		return err
	}

	articles, err := d.articleSvc.List(ctx.Request().Context(), int32(offset), int32(pageLimit))
	if err != nil {
		return err
	}

	totalPostsCount, err := d.articleSvc.Count(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.ArticleViewData, 0, len(articles))
	for _, article := range articles {
		viewData = append(viewData, dashboard.ArticleViewData{
			ID:         article.ID.String(),
			Title:      article.Title,
			Draft:      article.Draft,
			ReleasedAt: article.ReleaseDate.String(),
			Slug:       article.Slug,
		})
	}

	pagination := components.PaginationProps{
		CurrentPage: currentPage,
		TotalPages:  d.base.CalculateNumberOfPages(int(totalPostsCount), 7),
	}

	return dashboard.Articles(viewData, pagination, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (d Dashboard) ArticleCreate(ctx echo.Context) error {
	tags, err := d.tagSvc.All(ctx.Request().Context())
	if err != nil {
		return err
	}

	var keywords []dashboard.Keyword
	for _, tag := range tags {
		keywords = append(keywords, dashboard.Keyword{
			ID:    tag.ID.String(),
			Value: tag.Name,
		})
	}

	filenames, err := posts.GetAllFiles()
	if err != nil {
		return err
	}

	articles, err := d.articleSvc.All(ctx.Request().Context())
	if err != nil {
		return err
	}

	var usedFilenames []string
	for _, article := range articles {
		usedFilenames = append(usedFilenames, article.Filename)
	}

	var unusedFileNames []string
	for _, filename := range filenames {
		if !slices.Contains(usedFilenames, filename) {
			unusedFileNames = append(unusedFileNames, filename)
		}
	}

	return dashboard.CreateArticle(
		dashboard.CreateArticleViewData{
			Keywords:  keywords,
			Filenames: unusedFileNames,
		},
		csrf.Token(ctx.Request())).Render(views.ExtractRenderDeps(ctx))
}

func (d Dashboard) ArticleEdit(ctx echo.Context) error {
	slug := ctx.Param("slug")

	post, err := d.articleSvc.BySlug(ctx.Request().Context(), slug)
	if err != nil {
		return err
	}

	postContent, err := d.postManager.Parse(post.Filename)
	if err != nil {
		return err
	}

	var keywords string
	for i, tag := range post.Tags {
		if i == len(post.Tags)-1 {
			keywords = keywords + tag.Name
		} else {
			keywords = keywords + tag.Name + ", "
		}
	}

	return dashboard.ArticleEdit(dashboard.ArticleEditViewData{
		ID:          post.ID.String(),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		HeaderTitle: post.HeaderTitle,
		Filename:    post.Filename,
		Slug:        post.Slug,
		Excerpt:     post.Excerpt,
		Draft:       post.Draft,
		ReleasedAt:  post.ReleaseDate,
		ReadTime:    post.ReadTime,
		Content:     postContent,
		Keywords:    keywords,
	}, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

type articleUpdatePayload struct {
	Title       string `form:"title"`
	HeaderTitle string `form:"header-title"`
	Slug        string `form:"slug"`
	Excerpt     string `form:"excerpt"`
	EstReadTime int32  `form:"est-read-time"`
	UpdatedAt   string `form:"updated-at"`
	ReleasedAt  string `form:"released-at"`
	IsLive      string `form:"is-live"`
	Filename    string `form:"filename"`
}

func (d Dashboard) ArticleUpdate(ctx echo.Context) error {
	id := ctx.Param("id")

	var updateArticlePayload articleUpdatePayload
	if err := ctx.Bind(&updateArticlePayload); err != nil {
		return err
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	article, err := d.articleSvc.Update(ctx.Request().Context(), models.UpdateArticlePayload{
		ID:          parsedID,
		Title:       updateArticlePayload.Title,
		HeaderTitle: updateArticlePayload.HeaderTitle,
		Filename:    updateArticlePayload.Filename,
		Excerpt:     updateArticlePayload.Excerpt,
		Readtime:    updateArticlePayload.EstReadTime,
		TagIDs:      []uuid.UUID{},
	})
	if err != nil {
		return err
	}

	// parsedReleasedAt := carbon.Parse(updateArticlePayload.ReleasedAt).ToStdTime()
	//
	// var release bool
	// if updateArticlePayload.IsLive == "on" {
	// 	release = true
	// }

	return dashboard.ArticleEditForm(dashboard.ArticleEditViewData{
		ID:          article.ID.String(),
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
		Title:       article.Title,
		HeaderTitle: article.HeaderTitle,
		Filename:    article.Filename,
		Slug:        article.Slug,
		Excerpt:     article.Excerpt,
		Draft:       article.Draft,
		ReleasedAt:  article.ReleaseDate,
		ReadTime:    article.ReadTime,
	}, true).
		Render(views.ExtractRenderDeps(ctx))
}

func (d Dashboard) ArticleStore(ctx echo.Context) error {
	type newArticleFormPayload struct {
		Title             string   `form:"title"`
		HeaderTitle       string   `form:"header-title"`
		Excerpt           string   `form:"excerpt"`
		EstimatedReadTime string   `form:"estimated-read-time"`
		Filename          string   `form:"filename"`
		SelectedKeywords  []string `form:"selected-keyword"`
		Release           string   `form:"release"`
	}
	var postPayload newArticleFormPayload
	if err := ctx.Bind(&postPayload); err != nil {
		return err
	}

	var releaseNow bool
	if postPayload.Release == "on" {
		releaseNow = true
	}

	estimatedReadTime, err := strconv.Atoi(postPayload.EstimatedReadTime)
	if err != nil {
		return err
	}

	var tagIDs []uuid.UUID
	for _, seletedKeyword := range postPayload.SelectedKeywords {
		id, err := uuid.Parse(seletedKeyword)
		if err != nil {
			return err
		}

		tagIDs = append(tagIDs, id)
	}

	_, err = d.articleSvc.New(ctx.Request().Context(), models.NewArticlePayload{
		ReleaseNow:  releaseNow,
		Title:       postPayload.Title,
		HeaderTitle: postPayload.HeaderTitle,
		Filename:    postPayload.Filename,
		Excerpt:     postPayload.Excerpt,
		Readtime:    int32(estimatedReadTime),
		TagIDs:      tagIDs,
	})
	if err != nil {
		if errors.Is(err, models.ErrFailValidation) {
			var validationErrors validation.ValidationErrs
			if ok := errors.As(err, &validationErrors); !ok {
				return err
			}

			mappedErrors := make(map[string]components.InputError, len(validationErrors))
			for _, ve := range validationErrors {
				var msg string
				for i, cause := range ve.Causes() {
					if i == 0 {
						msg = cause.Error()
					} else {
						msg = fmt.Sprintf("%v, %v", msg, cause)
					}
				}

				mappedErrors[ve.Field()] = components.InputError{
					Msg:      msg,
					OldValue: ve.Value(),
				}
			}

			tags, err := d.tagSvc.All(ctx.Request().Context())
			if err != nil {
				return err
			}

			var keywords []dashboard.Keyword
			for _, tag := range tags {
				keywords = append(keywords, dashboard.Keyword{
					ID:    tag.ID.String(),
					Value: tag.Name,
				})
			}

			filenames, err := posts.GetAllFiles()
			if err != nil {
				return err
			}

			articles, err := d.articleSvc.All(ctx.Request().Context())
			if err != nil {
				return err
			}

			var usedFilenames []string
			for _, article := range articles {
				usedFilenames = append(usedFilenames, article.Filename)
			}

			var unusedFileNames []string
			for _, filename := range filenames {
				if !slices.Contains(usedFilenames, filename) {
					unusedFileNames = append(unusedFileNames, filename)
				}
			}

			return dashboard.CreateArticleFormContent(mappedErrors, keywords, unusedFileNames).
				Render(views.ExtractRenderDeps(ctx))
		}

		return err
	}

	ctx.Response().Writer.Header().Add("HX-Redirect", "/dashboard/articles")
	return nil
}

func (d Dashboard) TagStore(ctx echo.Context) error {
	type newTagFormPayload struct {
		Name             string   `form:"tag-name"`
		SelectedKeywords []string `form:"selected-keyword"`
	}
	var tagPayload newTagFormPayload
	if err := ctx.Bind(&tagPayload); err != nil {
		return err
	}

	if _, err := d.tagSvc.New(ctx.Request().Context(), tagPayload.Name); err != nil {
		return err
	}

	tags, err := d.tagSvc.All(ctx.Request().Context())
	if err != nil {
		return err
	}

	var keywords []dashboard.Keyword
	for _, tag := range tags {
		selected := false
		for _, selectedKW := range tagPayload.SelectedKeywords {
			if selectedKW == tag.ID.String() {
				selected = true
			}
		}

		kw := dashboard.Keyword{
			ID:       tag.ID.String(),
			Value:    tag.Name,
			Selected: selected,
		}
		keywords = append(keywords, kw)
	}

	return dashboard.KeywordsGrid(keywords).Render(views.ExtractRenderDeps(ctx))
}
