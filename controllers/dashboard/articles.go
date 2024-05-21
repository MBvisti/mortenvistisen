package dashboard

import (
	"slices"
	"strconv"

	"github.com/MBvisti/mortenvistisen/controllers/misc"
	"github.com/MBvisti/mortenvistisen/entity"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/MBvisti/mortenvistisen/views/validation"
	"github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func ArticlesIndex(ctx echo.Context, db database.Queries) error {
	page := ctx.QueryParam("page")

	var currentPage int
	if page == "" {
		currentPage = 1
	}
	if page != "" {
		cp, err := strconv.Atoi(page)
		if err != nil {
			return err
		}

		currentPage = cp
	}

	offset := 0
	if currentPage == 2 {
		offset = 7
	}

	if currentPage > 2 {
		offset = 7 * (currentPage - 1)
	}

	articles, err := db.QueryAllPosts(ctx.Request().Context(), int32(offset))
	if err != nil {
		return err
	}

	var totalPostCount int

	viewData := make([]dashboard.ArticleViewData, 0, len(articles))
	for i, article := range articles {
		if i == 0 {
			totalPostCount = int(article.TotalPostsCount)
		}
		viewData = append(viewData, dashboard.ArticleViewData{
			ID:         article.ID.String(),
			Title:      article.Title,
			Draft:      article.Draft,
			ReleasedAt: article.ReleasedAt.Time.String(),
			Slug:       article.Slug,
		})
	}

	pagination := components.PaginationPayload{
		CurrentPage:     currentPage,
		NextPage:        currentPage + 1,
		PrevPage:        currentPage - 1,
		HasNextNextPage: totalPostCount-7 >= 7,
	}

	if len(articles) < 7 {
		pagination.NoNextPage = true
	}

	return dashboard.Articles(viewData, pagination, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func ArticleCreate(ctx echo.Context, db database.Queries) error {
	tags, err := db.QueryAllTags(ctx.Request().Context())
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

	usedFilenames, err := db.QueryAllFilenames(ctx.Request().Context())
	if err != nil {
		return err
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
	db database.Queries,
	postManager posts.PostManager,
) error {
	slug := ctx.Param("slug")

	post, err := db.GetPostBySlug(ctx.Request().Context(), slug)
	if err != nil {
		return err
	}

	postContent, err := postManager.Parse(post.Filename)
	if err != nil {
		return err
	}

	tags, err := db.GetTagsForPost(ctx.Request().Context(), post.ID)
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

type articleUpdatePayload struct {
	Title       string `form:"title"`
	HeaderTitle string `form:"header-title"`
	Slug        string `form:"slug"`
	Excerpt     string `form:"excerpt"`
	EstReadTime int32  `form:"est-read-time"`
	UpdatedAt   string `form:"updated-at"`
	ReleasedAt  string `form:"released-at"`
	IsLive      string `form:"is-live"`
}

func ArticleUpdate(ctx echo.Context, db database.Queries, v *validator.Validate) error {
	id := ctx.Param("id")

	var updateArticlePayload articleUpdatePayload
	if err := ctx.Bind(&updateArticlePayload); err != nil {
		return err
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	parsedReleasedAt := carbon.Parse(updateArticlePayload.ReleasedAt).ToStdTime()

	var release bool
	if updateArticlePayload.IsLive == "on" {
		release = true
	}

	post, err := services.UpdatePost(
		ctx.Request().Context(),
		&db,
		v,
		entity.UpdatePost{
			ID:                parsedID,
			Title:             updateArticlePayload.Title,
			HeaderTitle:       updateArticlePayload.HeaderTitle,
			Excerpt:           updateArticlePayload.Excerpt,
			Slug:              updateArticlePayload.Slug,
			ReleaedAt:         parsedReleasedAt,
			ReleaseNow:        release,
			EstimatedReadTime: updateArticlePayload.EstReadTime,
		},
	)
	if err != nil {
		return err
	}

	return dashboard.ArticleEditForm(dashboard.ArticleEditViewData{
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
}

func ArticleStore(ctx echo.Context, db database.Queries, v *validator.Validate) error {
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

	if err := services.NewPost(ctx.Request().Context(), &db, v, entity.NewPost{
		Title:             postPayload.Title,
		HeaderTitle:       postPayload.HeaderTitle,
		Excerpt:           postPayload.Excerpt,
		ReleaseNow:        releaseNow,
		EstimatedReadTime: int64(estimatedReadTime),
		Filename:          postPayload.Filename,
	}, postPayload.SelectedKeywords); err != nil {
		e, ok := err.(validator.ValidationErrors)
		if !ok {
			return misc.InternalError(ctx)
		}

		tags, err := db.QueryAllTags(ctx.Request().Context())
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

		usedFilenames, err := db.QueryAllFilenames(ctx.Request().Context())
		if err != nil {
			return err
		}

		var unusedFileNames []string
		for _, filename := range filenames {
			if !slices.Contains(usedFilenames, filename) {
				unusedFileNames = append(unusedFileNames, filename)
			}
		}

		props := dashboard.CreateArticleFormProps{
			Title: validation.InputField{
				OldValue: postPayload.Title,
			},
			HeaderTitle: validation.InputField{
				OldValue: postPayload.HeaderTitle,
			},
			Excerpt: validation.InputField{
				OldValue: postPayload.Excerpt,
			},
			Filename: validation.InputField{
				OldValue: postPayload.Filename,
			},
			EstimatedReadTime: validation.InputField{
				OldValue: postPayload.EstimatedReadTime,
			},
			ReleaseNow: releaseNow,
		}

		for _, validationError := range e {
			switch validationError.StructField() {
			case "Title":
				props.Title.Invalid = true
				props.Title.InvalidMsg = "Title is not long enough"
			case "HeaderTitle":
				props.HeaderTitle.Invalid = true
				props.HeaderTitle.InvalidMsg = "Header title is not long enough"
			case "Excerpt":
				props.Excerpt.Invalid = true
				props.Excerpt.InvalidMsg = "Excerpt has to be between 130 and 160 chars long"
			case "EstimatedReadTime":
				props.EstimatedReadTime.Invalid = true
				props.EstimatedReadTime.InvalidMsg = "Est. time cannot be 0"
			case "Filename":
				props.Filename.Invalid = true
				props.Filename.InvalidMsg = "All filenames must end with .md"
			}
		}

		return dashboard.CreateArticleFormContent(
			props,
			keywords,
			unusedFileNames,
		).Render(views.ExtractRenderDeps(ctx))
	}

	ctx.Response().Writer.Header().Add("HX-Redirect", "/dashboard/articles")
	return nil
}

type newTagFormPayload struct {
	Name             string   `form:"tag-name"`
	SelectedKeywords []string `form:"selected-keyword"`
}

func TagStore(ctx echo.Context, db database.Queries) error {
	var tagPayload newTagFormPayload
	if err := ctx.Bind(&tagPayload); err != nil {
		return err
	}

	if err := db.InsertTag(ctx.Request().Context(), database.InsertTagParams{
		ID:   uuid.New(),
		Name: tagPayload.Name,
	}); err != nil {
		return err
	}

	tags, err := db.QueryAllTags(ctx.Request().Context())
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
