package psql

import (
	"context"
	"database/sql"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p Postgres) InsertArticle(
	ctx context.Context,
	data models.Article,
) error {
	createdAt := pgtype.Timestamp{
		Time:  data.CreatedAt,
		Valid: true,
	}
	updatedAt := pgtype.Timestamp{
		Time:  data.UpdatedAt,
		Valid: true,
	}

	releasedAt := pgtype.Timestamp{
		Time:  data.ReleaseDate,
		Valid: true,
	}

	if _, err := p.Queries.InsertPost(ctx, database.InsertPostParams{
		ID:          data.ID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Title:       data.Title,
		HeaderTitle: sql.NullString{String: data.HeaderTitle, Valid: true},
		Filename:    data.Filename,
		Slug:        data.Slug,
		Excerpt:     data.Excerpt,
		Draft:       data.Draft,
		ReleasedAt:  releasedAt,
		ReadTime:    sql.NullInt32{Int32: data.ReadTime, Valid: true},
	}); err != nil {
		return err
	}

	return nil
}

func (p Postgres) QueryAllArticles(ctx context.Context) ([]models.Article, error) {
	articles, err := p.Queries.QueryAllPosts(ctx, 0)
	if err != nil {
		return nil, err
	}

	var a []models.Article
	for _, article := range articles {
		var convertedTags []models.Tag
		tags, err := p.Queries.QueryTagsByPost(ctx, article.ID)
		if err != nil {
			return nil, err
		}

		for _, t := range tags {
			convertedTags = append(convertedTags, models.Tag{
				ID:   t.ID,
				Name: t.Name,
			})
		}
		a = append(a, models.Article{
			ID:          article.ID,
			CreatedAt:   article.CreatedAt.Time,
			UpdatedAt:   article.UpdatedAt.Time,
			Title:       article.Title,
			HeaderTitle: article.HeaderTitle.String,
			Filename:    article.Filename,
			Slug:        article.Slug,
			Excerpt:     article.Excerpt,
			Draft:       article.Draft,
			ReleaseDate: article.ReleasedAt.Time,
			ReadTime:    article.ReadTime.Int32,
			Tags:        convertedTags,
		})
	}

	return a, nil
}

func (p Postgres) QueryArticleByID(
	ctx context.Context,
	id uuid.UUID,
) (models.Article, error) {
	article, err := p.Queries.QueryPostByID(ctx, id)
	if err != nil {
		return models.Article{}, err
	}

	var convertedTags []models.Tag
	tags, err := p.Queries.QueryTagsByPost(ctx, id)
	if err != nil {
		return models.Article{}, err
	}

	for _, t := range tags {
		convertedTags = append(convertedTags, models.Tag{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	return models.Article{
		ID:          article.ID,
		CreatedAt:   article.CreatedAt.Time,
		UpdatedAt:   article.UpdatedAt.Time,
		Title:       article.Title,
		HeaderTitle: article.HeaderTitle.String,
		Filename:    article.Filename,
		Slug:        article.Slug,
		Excerpt:     article.Excerpt,
		Draft:       article.Draft,
		ReleaseDate: article.ReleasedAt.Time,
		ReadTime:    article.ReadTime.Int32,
		Tags:        convertedTags,
	}, nil
}

func (p Postgres) QueryArticleBySlug(
	ctx context.Context,
	slug string,
) (models.Article, error) {
	article, err := p.Queries.QueryPostBySlug(ctx, slug)
	if err != nil {
		return models.Article{}, err
	}

	var convertedTags []models.Tag
	tags, err := p.Queries.QueryTagsByPost(ctx, article.ID)
	if err != nil {
		return models.Article{}, err
	}

	for _, t := range tags {
		convertedTags = append(convertedTags, models.Tag{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	return models.Article{
		ID:          article.ID,
		CreatedAt:   article.CreatedAt.Time,
		UpdatedAt:   article.UpdatedAt.Time,
		Title:       article.Title,
		HeaderTitle: article.HeaderTitle.String,
		Filename:    article.Filename,
		Slug:        article.Slug,
		Excerpt:     article.Excerpt,
		Draft:       article.Draft,
		ReleaseDate: article.ReleasedAt.Time,
		ReadTime:    article.ReadTime.Int32,
		Tags:        convertedTags,
	}, nil
}

func (p Postgres) UpdateArticle(ctx context.Context, data models.Article) error {
	updatedAt := pgtype.Timestamp{
		Time:  data.UpdatedAt,
		Valid: true,
	}

	releasedAt := pgtype.Timestamp{
		Time:  data.ReleaseDate,
		Valid: true,
	}
	_, err := p.Queries.UpdatePost(ctx, database.UpdatePostParams{
		UpdatedAt:   updatedAt,
		Title:       data.Title,
		HeaderTitle: sql.NullString{String: data.HeaderTitle, Valid: true},
		Slug:        data.Slug,
		Excerpt:     data.Excerpt,
		Draft:       data.Draft,
		ReleasedAt:  releasedAt,
		ReadTime:    sql.NullInt32{Int32: data.ReadTime, Valid: true},
		ID:          data.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) CountArticles(ctx context.Context) (int64, error) {
	return p.Queries.QueryPostsCount(ctx)
}

func (p Postgres) ListArticles(
	ctx context.Context,
	filters models.QueryFilters,
	opts ...models.PaginationOption,
) ([]models.Article, error) {
	options := &models.PaginationOptions{}

	for _, opt := range opts {
		opt(options)
	}

	params := database.QueryPostsParams{
		Offset: sql.NullInt32{Int32: options.Offset, Valid: true},
		Limit:  sql.NullInt32{Int32: options.Limit, Valid: true},
	}

	posts, err := p.Queries.QueryPosts(ctx, params)
	if err != nil {
		return nil, err
	}

	articles := make([]models.Article, len(posts))
	for i, post := range posts {
		tags, err := p.Queries.QueryTagsByPost(ctx, post.ID)
		if err != nil {
			return nil, err
		}

		var convertedTags []models.Tag
		for _, t := range tags {
			convertedTags = append(convertedTags, models.Tag{
				ID:   t.ID,
				Name: t.Name,
			})
		}

		articles[i] = models.Article{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt.Time,
			UpdatedAt:   post.UpdatedAt.Time,
			Title:       post.Title,
			HeaderTitle: post.HeaderTitle.String,
			Filename:    post.Filename,
			Slug:        post.Slug,
			Excerpt:     post.Excerpt,
			Draft:       post.Draft,
			ReleaseDate: post.ReleasedAt.Time,
			ReadTime:    post.ReadTime.Int32,
			Tags:        convertedTags,
		}
	}

	return articles, nil
}
