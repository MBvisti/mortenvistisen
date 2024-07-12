package models

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
)

type articleStorage interface {
	InsertArticle(
		ctx context.Context,
		data Article,
	) error
	QueryArticleByID(
		ctx context.Context,
		id uuid.UUID,
	) (Article, error)
	QueryArticleBySlug(
		ctx context.Context,
		slug string,
	) (Article, error)
	UpdateArticle(ctx context.Context, data Article) error
	QueryTagsByIDs(ctx context.Context, ids []uuid.UUID) ([]Tag, error)
	AssociateTagsWithPost(
		ctx context.Context,
		postID uuid.UUID,
		tagIDs []uuid.UUID,
	) error
	ListArticles(
		ctx context.Context,
		filters QueryFilters,
		opts ...PaginationOption,
	) ([]Article, error)
	UpdateTagsPostAssociations(
		ctx context.Context,
		postID uuid.UUID,
		tagIDs []uuid.UUID,
	) error
	CountArticles(ctx context.Context) (int64, error)
	QueryAllArticles(ctx context.Context) ([]Article, error)
}

type ArticleService struct {
	articleStorage articleStorage
}

func NewArticleSvc(articleStorage articleStorage) ArticleService {
	return ArticleService{articleStorage}
}

type NewArticlePayload struct {
	ReleaseNow  bool
	Title       string
	HeaderTitle string
	Filename    string
	Excerpt     string
	Readtime    int32
	TagIDs      []uuid.UUID
}

func (a ArticleService) New(
	ctx context.Context,
	payload NewArticlePayload,
) (Article, error) {
	tags, err := a.articleStorage.QueryTagsByIDs(ctx, payload.TagIDs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Article{}, errors.Join(ErrNoRowWithIdentifier, err)
		}

		return Article{}, err
	}

	t := time.Now()
	article := Article{
		ID:          uuid.New(),
		CreatedAt:   t,
		UpdatedAt:   t,
		Title:       payload.Title,
		HeaderTitle: payload.HeaderTitle,
		Filename:    payload.Filename,
		Slug:        slug.MakeLang(payload.Title, "en"),
		Excerpt:     payload.Excerpt,
		Draft:       true,
		ReadTime:    payload.Readtime,
		Tags:        tags,
	}
	if err := validation.Validate(article, CreateArticleValidations()); err != nil {
		return Article{}, errors.Join(ErrFailValidation, err)
	}

	if payload.ReleaseNow {
		article.Draft = false
		article.ReleaseDate = time.Now()
	}

	if err := a.articleStorage.InsertArticle(ctx, article); err != nil {
		return Article{}, errors.Join(ErrUnrecoverableEvent, err)
	}

	if err := a.articleStorage.AssociateTagsWithPost(ctx, article.ID, payload.TagIDs); err != nil {
		return Article{}, errors.Join(ErrUnrecoverableEvent, err)
	}

	return article, nil
}

type UpdateArticlePayload struct {
	ID          uuid.UUID
	Title       string
	HeaderTitle string
	Filename    string
	Excerpt     string
	Readtime    int32
	TagIDs      []uuid.UUID
}

func (a ArticleService) Update(
	ctx context.Context,
	payload UpdateArticlePayload,
) (Article, error) {
	tags, err := a.articleStorage.QueryTagsByIDs(ctx, payload.TagIDs)
	if err != nil {
		return Article{}, err
	}

	article, err := a.articleStorage.QueryArticleByID(ctx, payload.ID)
	if err != nil {
		return Article{}, err
	}

	article.Title = payload.Title
	article.HeaderTitle = payload.HeaderTitle
	article.Filename = payload.Filename
	article.Slug = slug.MakeLang(payload.Title, "en")
	article.Excerpt = payload.Excerpt
	article.ReadTime = payload.Readtime
	article.Tags = tags

	if err := a.articleStorage.UpdateArticle(ctx, article); err != nil {
		return Article{}, err
	}

	if err := a.articleStorage.UpdateTagsPostAssociations(ctx, article.ID, payload.TagIDs); err != nil {
		return Article{}, err
	}

	return article, nil
}

func (a ArticleService) List(
	ctx context.Context,
	offset int32,
	limit int32,
) ([]Article, error) {
	articles, err := a.articleStorage.ListArticles(ctx, nil, WithPagination(limit, offset))
	if err != nil {
		slog.ErrorContext(ctx, "could not get list of articles", "error", err)
		return nil, err
	}

	return articles, nil
}

func (a ArticleService) BySlug(ctx context.Context, slug string) (Article, error) {
	return a.articleStorage.QueryArticleBySlug(ctx, slug)
}

func (a ArticleService) Count(ctx context.Context) (int64, error) {
	count, err := a.articleStorage.CountArticles(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a ArticleService) All(ctx context.Context) ([]Article, error) {
	return a.articleStorage.QueryAllArticles(ctx)
}
