package models

import (
	"context"
	"time"

	"github.com/MBvisti/mortenvistisen/models/internal/database"
	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Filename    string
	Slug        string
	Excerpt     string
	Draft       bool
	ReleasedAt  time.Time
	ReadTime    int32
	HeaderTitle string
}

type ArticleModel struct {
	db *database.Queries
}

func NewArticle(db *database.Queries) ArticleModel {
	return ArticleModel{
		db: db,
	}
}

func (a *ArticleModel) List(
	ctx context.Context,
	opts ...listOpt,
) ([]Article, error) {
	options := &listOptions{}

	for _, opt := range opts {
		opt(options)
	}

	articleModels, err := a.db.QueryPosts(ctx, database.QueryPostsParams{
		Offset: options.offset,
		Limit:  options.limit,
	})
	if err != nil {
		return nil, err
	}

	articles := make([]Article, len(articleModels))
	for i, article := range articleModels {
		articles[i] = Article{
			ID:          article.ID,
			UpdatedAt:   article.CreatedAt.Time,
			CreatedAt:   article.CreatedAt.Time,
			Title:       article.Title,
			Filename:    article.Filename,
			Slug:        article.Slug,
			Excerpt:     article.Excerpt,
			Draft:       article.Draft,
			ReleasedAt:  article.ReleasedAt.Time,
			ReadTime:    article.ReadTime.Int32,
			HeaderTitle: article.HeaderTitle.String,
		}
	}

	return articles, nil
}
