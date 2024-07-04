package dashboard

import (
	"slices"
	"strconv"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func ArticlesIndex(ctx echo.Context, articleModel models.ArticleService) error {
	page := ctx.QueryParam("page")
	pageLimit := 7

	offset, currentPage, err := controllers.GetOffsetAndCurrPage(page, pageLimit)
	if err != nil {
		return err
	}

	articles, err := articleModel.List(ctx.Request().Context(), int32(offset), int32(pageLimit))
	if err != nil {
		return err
	}

	totalPostsCount, err := articleModel.Count(ctx.Request().Context())
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
		TotalPages:  controllers.CalculateNumberOfPages(int(totalPostsCount), 7),
	}

	return dashboard.Articles(viewData, pagination, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func ArticleCreate(
	ctx echo.Context,
	articleModel models.ArticleService,
	tagModel models.TagService,
) error {
	tags, err := tagModel.All(ctx.Request().Context())
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

	articles, err := articleModel.All(ctx.Request().Context())
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

func ArticleEdit(
	ctx echo.Context,
	articleModel models.ArticleService,
	postManager posts.PostManager,
) error {
	slug := ctx.Param("slug")

	post, err := articleModel.BySlug(ctx.Request().Context(), slug)
	if err != nil {
		return err
	}

	postContent, err := postManager.Parse(post.Filename)
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

func ArticleUpdate(
	ctx echo.Context,
	articleModel models.ArticleService,
) error {
	id := ctx.Param("id")

	var updateArticlePayload articleUpdatePayload
	if err := ctx.Bind(&updateArticlePayload); err != nil {
		return err
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	article, err := articleModel.Update(ctx.Request().Context(), models.UpdateArticlePayload{
		ID:          parsedID,
		Title:       updateArticlePayload.Title,
		HeaderTitle: updateArticlePayload.HeaderTitle,
		Filename:    updateArticlePayload.Filename,
		Slug:        updateArticlePayload.Slug,
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

type newArticleFormPayload struct {
	Title             string   `form:"title"`
	HeaderTitle       string   `form:"header-title"`
	Excerpt           string   `form:"excerpt"`
	EstimatedReadTime string   `form:"estimated-read-time"`
	Filename          string   `form:"filename"`
	SelectedKeywords  []string `form:"selected-keyword"`
	Release           string   `form:"release"`
	Slug              string   `form:"slug"`
}

func ArticleStore(ctx echo.Context, articleModel models.ArticleService) error {
	var postPayload newArticleFormPayload
	if err := ctx.Bind(&postPayload); err != nil {
		return err
	}

	// var releaseNow bool
	// if postPayload.Release == "on" {
	// 	releaseNow = true
	// }
	//
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

	_, err = articleModel.New(ctx.Request().Context(), models.NewArticlePayload{
		Title:       postPayload.Title,
		HeaderTitle: postPayload.HeaderTitle,
		Filename:    postPayload.Filename,
		Slug:        postPayload.Slug,
		Excerpt:     postPayload.EstimatedReadTime,
		Readtime:    int32(estimatedReadTime),
		TagIDs:      tagIDs,
	})
	if err != nil {
		return err
	}

	ctx.Response().Writer.Header().Add("HX-Redirect", "/dashboard/articles")
	return nil
}

type newTagFormPayload struct {
	Name             string   `form:"tag-name"`
	SelectedKeywords []string `form:"selected-keyword"`
}

func TagStore(ctx echo.Context,
	tagModel models.TagService,
) error {
	var tagPayload newTagFormPayload
	if err := ctx.Bind(&tagPayload); err != nil {
		return err
	}

	if _, err := tagModel.New(ctx.Request().Context(), tagPayload.Name); err != nil {
		return err
	}

	tags, err := tagModel.All(ctx.Request().Context())
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
