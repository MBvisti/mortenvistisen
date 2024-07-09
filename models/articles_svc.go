package models

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
)

type articleStorage interface {
	InsertArticle(
		ctx context.Context,
		data domain.Article,
	) error
	QueryArticleByID(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Article, error)
	QueryArticleBySlug(
		ctx context.Context,
		slug string,
	) (domain.Article, error)
	UpdateArticle(ctx context.Context, data domain.Article) error
	QueryTagsByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Tag, error)
	AssociateTagsWithPost(
		ctx context.Context,
		postID uuid.UUID,
		tagIDs []uuid.UUID,
	) error
	ListArticles(
		ctx context.Context,
		filters QueryFilters,
		opts ...PaginationOption,
	) ([]domain.Article, error)
	UpdateTagsPostAssociations(
		ctx context.Context,
		postID uuid.UUID,
		tagIDs []uuid.UUID,
	) error
	CountArticles(ctx context.Context) (int64, error)
	QueryAllArticles(ctx context.Context) ([]domain.Article, error)
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
	Slug        string
	Excerpt     string
	Readtime    int32
	TagIDs      []uuid.UUID
}

func (a ArticleService) New(
	ctx context.Context,
	payload NewArticlePayload,
) (domain.Article, error) {
	tags, err := a.articleStorage.QueryTagsByIDs(ctx, payload.TagIDs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Article{}, errors.Join(ErrNoRowWithIdentifier, err)
		}

		return domain.Article{}, err
	}

	article, err := domain.NewArticle(
		payload.Title,
		payload.HeaderTitle,
		payload.Filename,
		slug.MakeLang(payload.Title, "en"),
		payload.Excerpt,
		payload.Readtime,
		tags,
	)
	if err != nil {
		return domain.Article{}, errors.Join(ErrFailValidation, err)
	}

	if payload.ReleaseNow {
		article.Draft = false
		article.ReleaseDate = time.Now()
	}

	if err := a.articleStorage.InsertArticle(ctx, article); err != nil {
		return domain.Article{}, errors.Join(ErrUnrecoverableEvent, err)
	}

	if err := a.articleStorage.AssociateTagsWithPost(ctx, article.ID, payload.TagIDs); err != nil {
		return domain.Article{}, errors.Join(ErrUnrecoverableEvent, err)
	}

	return article, nil
}

type UpdateArticlePayload struct {
	ID          uuid.UUID
	Title       string
	HeaderTitle string
	Filename    string
	Slug        string
	Excerpt     string
	Readtime    int32
	TagIDs      []uuid.UUID
}

func (a ArticleService) Update(
	ctx context.Context,
	payload UpdateArticlePayload,
) (domain.Article, error) {
	tags, err := a.articleStorage.QueryTagsByIDs(ctx, payload.TagIDs)
	if err != nil {
		return domain.Article{}, err
	}

	article, err := a.articleStorage.QueryArticleByID(ctx, payload.ID)
	if err != nil {
		return domain.Article{}, err
	}

	if err := article.Update(
		payload.Title,
		payload.HeaderTitle,
		payload.Filename,
		payload.Slug,
		payload.Excerpt,
		payload.Readtime,
		tags,
	); err != nil {
		return domain.Article{}, err
	}

	if err := a.articleStorage.UpdateArticle(ctx, article); err != nil {
		return domain.Article{}, err
	}

	if err := a.articleStorage.UpdateTagsPostAssociations(ctx, article.ID, payload.TagIDs); err != nil {
		return domain.Article{}, err
	}

	return article, nil
}

func (a ArticleService) List(
	ctx context.Context,
	offset int32,
	limit int32,
) ([]domain.Article, error) {
	articles, err := a.articleStorage.ListArticles(ctx, nil, WithPagination(limit, offset))
	if err != nil {
		slog.ErrorContext(ctx, "could not get list of articles", "error", err)
		return nil, err
	}

	return articles, nil
}

func (a ArticleService) BySlug(ctx context.Context, slug string) (domain.Article, error) {
	return a.articleStorage.QueryArticleBySlug(ctx, slug)
}

func (a ArticleService) Count(ctx context.Context) (int64, error) {
	count, err := a.articleStorage.CountArticles(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a ArticleService) All(ctx context.Context) ([]domain.Article, error) {
	return a.articleStorage.QueryAllArticles(ctx)
}
