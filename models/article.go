package models

import (
	"context"
	"time"

	"github.com/MBvisti/mortenvistisen/models/internal/db"
	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	HeaderTitle string
	Filename    string
	Slug        string
	Excerpt     string
	Draft       bool
	ReleaseDate time.Time
	ReadTime    int32
	Tags        []Tag
}

type PaginatedArticles struct {
	Articles []Article
	Total    int64
	Page     int32
	PageSize int32
}

func GetArticlesPage(
	ctx context.Context,
	page int32,
	pageSize int32,
	dbtx db.DBTX,
) (PaginatedArticles, error) {
	offset := (page - 1) * pageSize

	total, err := db.Stmts.QueryArticlesCount(ctx, dbtx)
	if err != nil {
		return PaginatedArticles{}, err
	}

	articles, err := db.Stmts.QueryArticlesPage(
		ctx,
		dbtx,
		db.QueryArticlesPageParams{
			Limit:  pageSize,
			Offset: offset,
		},
	)
	if err != nil {
		return PaginatedArticles{}, err
	}

	result := make([]Article, len(articles))
	for i, article := range articles {
		result[i] = Article{
			ID:          article.ID,
			CreatedAt:   article.CreatedAt.Time,
			UpdatedAt:   article.UpdatedAt.Time,
			Title:       article.Title,
			Filename:    article.Filename,
			Slug:        article.Slug,
			Excerpt:     article.Excerpt,
			Draft:       article.Draft,
			ReleaseDate: article.ReleaseDate.Time,
			ReadTime:    article.ReadTime.Int32,
		}

		// Get tags for each article
		tags, err := db.Stmts.QueryArticleTags(ctx, dbtx, article.ID)
		if err != nil {
			return PaginatedArticles{}, err
		}

		result[i].Tags = make([]Tag, len(tags))
		for j, tag := range tags {
			result[i].Tags[j] = Tag{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}
	}

	return PaginatedArticles{
		Articles: result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func GetAllArticles(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Article, error) {
	articles, err := db.Stmts.QueryArticles(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	result := make([]Article, len(articles))
	for i, article := range articles {
		result[i] = Article{
			ID:          article.ID,
			CreatedAt:   article.CreatedAt.Time,
			UpdatedAt:   article.UpdatedAt.Time,
			Title:       article.Title,
			Filename:    article.Filename,
			Slug:        article.Slug,
			Excerpt:     article.Excerpt,
			Draft:       article.Draft,
			ReleaseDate: article.ReleaseDate.Time,
			ReadTime:    article.ReadTime.Int32,
		}

		// Get tags for each article
		tags, err := db.Stmts.QueryArticleTags(ctx, dbtx, article.ID)
		if err != nil {
			return nil, err
		}

		result[i].Tags = make([]Tag, len(tags))
		for j, tag := range tags {
			result[i].Tags[j] = Tag{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}
	}

	return result, nil
}

func GetArticleBySlug(
	ctx context.Context,
	slug string,
	dbtx db.DBTX,
) (*Article, error) {
	article, err := db.Stmts.QueryArticleBySlug(ctx, dbtx, slug)
	if err != nil {
		return nil, err
	}

	result := &Article{
		ID:          article.ID,
		CreatedAt:   article.CreatedAt.Time,
		UpdatedAt:   article.UpdatedAt.Time,
		Title:       article.Title,
		Filename:    article.Filename,
		Slug:        article.Slug,
		Excerpt:     article.Excerpt,
		Draft:       article.Draft,
		ReleaseDate: article.ReleaseDate.Time,
		ReadTime:    article.ReadTime.Int32,
	}

	// Get tags for the article
	tags, err := db.Stmts.QueryArticleTags(ctx, dbtx, article.ID)
	if err != nil {
		return nil, err
	}

	result.Tags = make([]Tag, len(tags))
	for i, tag := range tags {
		result.Tags[i] = Tag{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return result, nil
}
