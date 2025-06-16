package models

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvisti/mortenvistisen/models/internal/db"
)

type Article struct {
	ID               uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	FirstPublishedAt time.Time
	IsPublished      bool
	Title            string
	Excerpt          string
	MetaTitle        string
	MetaDescription  string
	Slug             string
	ImageLink        string
	Content          string
	ReadTime         int32
	Tags             []ArticleTag
}

func populateArticleTags(
	ctx context.Context,
	dbtx db.DBTX,
	articles []Article,
) error {
	for i := range articles {
		tags, err := GetArticleTagsByArticleID(ctx, dbtx, articles[i].ID)
		if err != nil {
			return err
		}
		articles[i].Tags = tags
	}
	return nil
}

func GetArticleByID(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (Article, error) {
	row, err := db.Stmts.QueryArticleByID(ctx, dbtx, id)
	if err != nil {
		return Article{}, err
	}

	tags, err := GetArticleTagsByArticleID(ctx, dbtx, id)
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:               row.ID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		FirstPublishedAt: row.FirstPublishedAt.Time,
		IsPublished:      row.IsPublished.Bool,
		Title:            row.Title,
		Excerpt:          row.Excerpt,
		MetaTitle:        row.MetaTitle,
		MetaDescription:  row.MetaDescription,
		Slug:             row.Slug,
		ImageLink:        row.ImageLink.String,
		Content:          row.Content.String,
		ReadTime:         row.ReadTime.Int32,
		Tags:             tags,
	}, nil
}

func GetArticleByTitle(
	ctx context.Context,
	dbtx db.DBTX,
	title string,
) (Article, error) {
	row, err := db.Stmts.QueryArticleByTitle(ctx, dbtx, title)
	if err != nil {
		return Article{}, err
	}

	tags, err := GetArticleTagsByArticleID(ctx, dbtx, row.ID)
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:               row.ID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		FirstPublishedAt: row.FirstPublishedAt.Time,
		IsPublished:      row.IsPublished.Bool,
		Title:            row.Title,
		Excerpt:          row.Excerpt,
		MetaTitle:        row.MetaTitle,
		MetaDescription:  row.MetaDescription,
		Slug:             row.Slug,
		ImageLink:        row.ImageLink.String,
		Content:          row.Content.String,
		ReadTime:         row.ReadTime.Int32,
		Tags:             tags,
	}, nil
}

func GetArticleBySlug(
	ctx context.Context,
	dbtx db.DBTX,
	slug string,
) (Article, error) {
	row, err := db.Stmts.QueryArticleBySlug(ctx, dbtx, slug)
	if err != nil {
		return Article{}, err
	}

	tags, err := GetArticleTagsByArticleID(ctx, dbtx, row.ID)
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:               row.ID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		FirstPublishedAt: row.FirstPublishedAt.Time,
		IsPublished:      row.IsPublished.Bool,
		Title:            row.Title,
		Excerpt:          row.Excerpt,
		MetaTitle:        row.MetaTitle,
		MetaDescription:  row.MetaDescription,
		Slug:             row.Slug,
		ImageLink:        row.ImageLink.String,
		Content:          row.Content.String,
		ReadTime:         row.ReadTime.Int32,
		Tags:             tags,
	}, nil
}

func GetArticles(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Article, error) {
	rows, err := db.Stmts.QueryArticles(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		articles[i] = Article{
			ID:               row.ID,
			CreatedAt:        row.CreatedAt.Time,
			UpdatedAt:        row.UpdatedAt.Time,
			FirstPublishedAt: row.FirstPublishedAt.Time,
			IsPublished:      row.IsPublished.Bool,
			Title:            row.Title,
			Excerpt:          row.Excerpt,
			MetaTitle:        row.MetaTitle,
			MetaDescription:  row.MetaDescription,
			Slug:             row.Slug,
			ImageLink:        row.ImageLink.String,
			Content:          row.Content.String,
			ReadTime:         row.ReadTime.Int32,
		}
	}

	if err := populateArticleTags(ctx, dbtx, articles); err != nil {
		return nil, err
	}

	return articles, nil
}

func GetPublishedArticles(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Article, error) {
	rows, err := db.Stmts.QueryPublishedArticles(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		articles[i] = Article{
			ID:               row.ID,
			CreatedAt:        row.CreatedAt.Time,
			UpdatedAt:        row.UpdatedAt.Time,
			FirstPublishedAt: row.FirstPublishedAt.Time,
			IsPublished:      row.IsPublished.Bool,
			Title:            row.Title,
			Excerpt:          row.Excerpt,
			MetaTitle:        row.MetaTitle,
			MetaDescription:  row.MetaDescription,
			Slug:             row.Slug,
			ImageLink:        row.ImageLink.String,
			Content:          row.Content.String,
			ReadTime:         row.ReadTime.Int32,
		}
	}

	if err := populateArticleTags(ctx, dbtx, articles); err != nil {
		return nil, err
	}

	return articles, nil
}

func GetDraftArticles(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Article, error) {
	rows, err := db.Stmts.QueryDraftArticles(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		articles[i] = Article{
			ID:               row.ID,
			CreatedAt:        row.CreatedAt.Time,
			UpdatedAt:        row.UpdatedAt.Time,
			FirstPublishedAt: row.FirstPublishedAt.Time,
			IsPublished:      row.IsPublished.Bool,
			Title:            row.Title,
			Excerpt:          row.Excerpt,
			MetaTitle:        row.MetaTitle,
			MetaDescription:  row.MetaDescription,
			Slug:             row.Slug,
			ImageLink:        row.ImageLink.String,
			Content:          row.Content.String,
			ReadTime:         row.ReadTime.Int32,
		}
	}

	if err := populateArticleTags(ctx, dbtx, articles); err != nil {
		return nil, err
	}

	return articles, nil
}

type PaginationResult struct {
	Articles    []Article
	TotalCount  int64
	Page        int
	PageSize    int
	TotalPages  int
	HasNext     bool
	HasPrevious bool
}

type SortConfig struct {
	Field string
	Order string
}

var allowedArticleSortFields = map[string]string{
	"title":      "title",
	"created_at": "created_at",
	"updated_at": "updated_at",
	"status":     "is_published",
	"published":  "first_published_at",
	"read_time":  "read_time",
}

func GetArticlesPaginated(
	ctx context.Context,
	dbtx db.DBTX,
	page int,
	pageSize int,
) (PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Limit max page size
	}

	offset := (page - 1) * pageSize

	// Get total count
	totalCount, err := db.Stmts.CountArticles(ctx, dbtx)
	if err != nil {
		return PaginationResult{}, err
	}

	// Get paginated articles
	rows, err := db.Stmts.QueryArticlesPaginated(
		ctx,
		dbtx,
		db.QueryArticlesPaginatedParams{
			//nolint:gosec // pageSize is bounded above
			Limit: int32(pageSize),
			//nolint:gosec // offset is calculated from bounded values
			Offset: int32(
				offset,
			),
		},
	)
	if err != nil {
		return PaginationResult{}, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		articles[i] = Article{
			ID:               row.ID,
			CreatedAt:        row.CreatedAt.Time,
			UpdatedAt:        row.UpdatedAt.Time,
			FirstPublishedAt: row.FirstPublishedAt.Time,
			IsPublished:      row.IsPublished.Bool,
			Title:            row.Title,
			Excerpt:          row.Excerpt,
			MetaTitle:        row.MetaTitle,
			MetaDescription:  row.MetaDescription,
			Slug:             row.Slug,
			ImageLink:        row.ImageLink.String,
			Content:          row.Content.String,
			ReadTime:         row.ReadTime.Int32,
		}
	}

	if err := populateArticleTags(ctx, dbtx, articles); err != nil {
		return PaginationResult{}, err
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return PaginationResult{
		Articles:    articles,
		TotalCount:  totalCount,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}, nil
}

func GetArticlesSorted(
	ctx context.Context,
	dbtx db.DBTX,
	page int,
	pageSize int,
	sort SortConfig,
) (PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// Get total count first
	totalCount, err := db.Stmts.CountArticles(ctx, dbtx)
	if err != nil {
		return PaginationResult{}, err
	}

	// Build the sortable query using Squirrel
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query := psql.Select(
		"id", "created_at", "updated_at", "first_published_at",
		"title", "excerpt", "meta_title", "meta_description",
		"slug", "image_link", "content", "read_time", "is_published",
	).From("articles")

	// Add sorting if valid field provided
	if sort.Field != "" && sort.Order != "" {
		if field, ok := allowedArticleSortFields[sort.Field]; ok {
			orderClause := field
			if sort.Order == "desc" {
				orderClause += " DESC"
			} else {
				orderClause += " ASC"
			}
			query = query.OrderBy(orderClause)
		} else {
			// Default sorting if invalid field
			query = query.OrderBy("created_at DESC")
		}
	} else {
		// Default sorting
		query = query.OrderBy("created_at DESC")
	}

	// Add pagination
	if pageSize >= 0 && offset >= 0 {
		//nolint:gosec // not needed
		query = query.Limit(uint64(pageSize)).
			Offset(uint64(offset))
	}

	// Build SQL
	sql, args, err := query.ToSql()
	if err != nil {
		return PaginationResult{}, err
	}

	// Execute query
	rows, err := dbtx.Query(ctx, sql, args...)
	if err != nil {
		return PaginationResult{}, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		var createdAt, updatedAt, firstPublishedAt pgtype.Timestamptz
		var imageLink, content pgtype.Text
		var readTime pgtype.Int4
		var isPublished pgtype.Bool

		err := rows.Scan(
			&a.ID, &createdAt, &updatedAt, &firstPublishedAt,
			&a.Title, &a.Excerpt, &a.MetaTitle, &a.MetaDescription,
			&a.Slug, &imageLink, &content, &readTime, &isPublished,
		)
		if err != nil {
			return PaginationResult{}, err
		}

		// Convert pgtype values
		a.CreatedAt = createdAt.Time
		a.UpdatedAt = updatedAt.Time
		a.FirstPublishedAt = firstPublishedAt.Time
		a.ImageLink = imageLink.String
		a.Content = content.String
		a.ReadTime = readTime.Int32
		a.IsPublished = isPublished.Bool

		articles = append(articles, a)
	}

	if err = rows.Err(); err != nil {
		return PaginationResult{}, err
	}

	// Populate tags for all articles
	if err := populateArticleTags(ctx, dbtx, articles); err != nil {
		return PaginationResult{}, err
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return PaginationResult{
		Articles:    articles,
		TotalCount:  totalCount,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}, nil
}

type NewArticlePayload struct {
	Title           string `validate:"required,max=100"`
	Excerpt         string `validate:"required,max=255"`
	MetaTitle       string `validate:"required,max=100"`
	MetaDescription string `validate:"required,max=100"`
	Slug            string `validate:"required,max=255"`
	ImageLink       string `validate:"omitempty,max=255"`
	Content         string
	ReadTime        int32
	TagIDs          []string
}

func NewArticle(
	ctx context.Context,
	dbtx db.DBTX,
	data NewArticlePayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		slog.ErrorContext(
			ctx,
			"could not validate new article payload",
			"error",
			err,
			"data",
			data,
		)
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	article := Article{
		ID:              uuid.New(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Title:           data.Title,
		Excerpt:         data.Excerpt,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            data.Slug,
		ImageLink:       data.ImageLink,
		Content:         data.Content,
		ReadTime:        data.ReadTime,
	}

	_, err := db.Stmts.InsertArticle(ctx, dbtx, db.InsertArticleParams{
		ID: article.ID,
		CreatedAt: pgtype.Timestamptz{
			Time:  article.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  article.UpdatedAt,
			Valid: true,
		},
		Title:           article.Title,
		Excerpt:         article.Excerpt,
		MetaTitle:       article.MetaTitle,
		MetaDescription: article.MetaDescription,
		Slug:            article.Slug,
		ImageLink: sql.NullString{
			String: article.ImageLink,
			Valid:  article.ImageLink != "",
		},
		Content: sql.NullString{
			String: article.Content,
			Valid:  article.Content != "",
		},
		ReadTime: sql.NullInt32{
			Int32: article.ReadTime,
			Valid: article.ReadTime > 0,
		},
	})
	if err != nil {
		return Article{}, err
	}

	// Create tag connections
	for _, tagIDStr := range data.TagIDs {
		tagID, err := uuid.Parse(tagIDStr)
		if err != nil {
			continue // Skip invalid UUIDs
		}
		_, err = NewArticleTagConnection(ctx, dbtx, article.ID, tagID)
		if err != nil {
			return Article{}, err
		}
	}

	// Fetch the created article with tags
	return GetArticleByID(ctx, dbtx, article.ID)
}

type UpdateArticlePayload struct {
	ID              uuid.UUID `validate:"required,uuid"`
	UpdatedAt       time.Time `validate:"required"`
	IsPublished     bool
	Title           string `validate:"required,max=100"`
	Excerpt         string `validate:"required,max=255"`
	MetaTitle       string `validate:"required,max=100"`
	MetaDescription string `validate:"required,max=100"`
	Slug            string `validate:"required,max=255"`
	ImageLink       string `validate:"omitempty,max=255"`
	Content         string
	ReadTime        int32 `validate:"min=1,max=999"`
	TagIDs          []string
}

func UpdateArticle(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateArticlePayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	_, err := db.Stmts.UpdateArticle(ctx, dbtx, db.UpdateArticleParams{
		ID:              data.ID,
		UpdatedAt:       pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
		Title:           data.Title,
		Excerpt:         data.Excerpt,
		IsPublished:     sql.NullBool{Bool: data.IsPublished, Valid: true},
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            data.Slug,
		ImageLink: sql.NullString{
			String: data.ImageLink,
			Valid:  data.ImageLink != "",
		},
		Content: sql.NullString{
			String: data.Content,
			Valid:  data.Content != "",
		},
		ReadTime: sql.NullInt32{
			Int32: data.ReadTime,
			Valid: data.ReadTime > 0,
		},
	})
	if err != nil {
		return Article{}, err
	}

	err = DeleteArticleTagConnectionsByArticleID(ctx, dbtx, data.ID)
	if err != nil {
		return Article{}, err
	}

	for _, tagIDStr := range data.TagIDs {
		tagID, err := uuid.Parse(tagIDStr)
		if err != nil {
			continue // Skip invalid UUIDs
		}
		_, err = NewArticleTagConnection(ctx, dbtx, data.ID, tagID)
		if err != nil {
			return Article{}, err
		}
	}

	return GetArticleByID(ctx, dbtx, data.ID)
}

type UpdateArticleContentPayload struct {
	ID        uuid.UUID `validate:"required,uuid"`
	UpdatedAt time.Time `validate:"required"`
	Content   string
}

func UpdateArticleContent(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateArticleContentPayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.UpdateArticleContent(
		ctx,
		dbtx,
		db.UpdateArticleContentParams{
			ID:        data.ID,
			UpdatedAt: pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
			Content: sql.NullString{
				String: data.Content,
				Valid:  data.Content != "",
			},
		},
	)
	if err != nil {
		return Article{}, err
	}

	tags, err := GetArticleTagsByArticleID(ctx, dbtx, row.ID)
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:               row.ID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		FirstPublishedAt: row.FirstPublishedAt.Time,
		IsPublished:      row.IsPublished.Bool,
		Title:            row.Title,
		Excerpt:          row.Excerpt,
		MetaTitle:        row.MetaTitle,
		MetaDescription:  row.MetaDescription,
		Slug:             row.Slug,
		ImageLink:        row.ImageLink.String,
		Content:          row.Content.String,
		ReadTime:         row.ReadTime.Int32,
		Tags:             tags,
	}, nil
}

type UpdateArticleMetadataPayload struct {
	ID              uuid.UUID `validate:"required,uuid"`
	UpdatedAt       time.Time `validate:"required"`
	Title           string    `validate:"required,max=100"`
	Excerpt         string    `validate:"required,max=255"`
	MetaTitle       string    `validate:"required,max=100"`
	MetaDescription string    `validate:"required,max=100"`
	Slug            string    `validate:"required,max=255"`
	ImageLink       string    `validate:"omitempty,max=255"`
}

func UpdateArticleMetadata(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateArticleMetadataPayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.UpdateArticleMetadata(
		ctx,
		dbtx,
		db.UpdateArticleMetadataParams{
			ID: data.ID,
			UpdatedAt: pgtype.Timestamptz{
				Time:  data.UpdatedAt,
				Valid: true,
			},
			Title:           data.Title,
			Excerpt:         data.Excerpt,
			MetaTitle:       data.MetaTitle,
			MetaDescription: data.MetaDescription,
			Slug:            data.Slug,
			ImageLink: sql.NullString{
				String: data.ImageLink,
				Valid:  data.ImageLink != "",
			},
		},
	)
	if err != nil {
		return Article{}, err
	}

	tags, err := GetArticleTagsByArticleID(ctx, dbtx, row.ID)
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:               row.ID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		FirstPublishedAt: row.FirstPublishedAt.Time,
		IsPublished:      row.IsPublished.Bool,
		Title:            row.Title,
		Excerpt:          row.Excerpt,
		MetaTitle:        row.MetaTitle,
		MetaDescription:  row.MetaDescription,
		Slug:             row.Slug,
		ImageLink:        row.ImageLink.String,
		Content:          row.Content.String,
		ReadTime:         row.ReadTime.Int32,
		Tags:             tags,
	}, nil
}

type PublishArticlePayload struct {
	ID  uuid.UUID `validate:"required,uuid"`
	Now time.Time // TODO: validate
}

func PublishArticle(
	ctx context.Context,
	dbtx db.DBTX,
	data PublishArticlePayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.PublishArticle(ctx, dbtx, db.PublishArticleParams{
		ID:        data.ID,
		UpdatedAt: pgtype.Timestamptz{Time: data.Now, Valid: true},
		FirstPublishedAt: pgtype.Timestamptz{
			Time:  data.Now,
			Valid: true,
		},
		IsPublished: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return Article{}, err
	}

	tags, err := GetArticleTagsByArticleID(ctx, dbtx, row.ID)
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:               row.ID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		FirstPublishedAt: row.FirstPublishedAt.Time,
		IsPublished:      row.IsPublished.Bool,
		Title:            row.Title,
		Excerpt:          row.Excerpt,
		MetaTitle:        row.MetaTitle,
		MetaDescription:  row.MetaDescription,
		Slug:             row.Slug,
		ImageLink:        row.ImageLink.String,
		Content:          row.Content.String,
		ReadTime:         row.ReadTime.Int32,
		Tags:             tags,
	}, nil
}

// type UnpublishArticlePayload struct {
// 	ID        uuid.UUID `validate:"required,uuid"`
// 	UpdatedAt time.Time `validate:"required"`
// }
//
// func UnpublishArticle(
// 	ctx context.Context,
// 	dbtx db.DBTX,
// 	data UnpublishArticlePayload,
// ) (Article, error) {
// 	if err := validate.Struct(data); err != nil {
// 		return Article{}, errors.Join(ErrDomainValidation, err)
// 	}
//
// 	row, err := db.Stmts.UnpublishArticle(ctx, dbtx, db.UnpublishArticleParams{
// 		ID:        data.ID,
// 		UpdatedAt: pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
// 	})
// 	if err != nil {
// 		return Article{}, err
// 	}
//
// 	tags, err := GetArticleTagsByArticleID(ctx, dbtx, row.ID)
// 	if err != nil {
// 		return Article{}, err
// 	}
//
// 	return Article{
// 		ID:               row.ID,
// 		CreatedAt:        row.CreatedAt.Time,
// 		UpdatedAt:        row.UpdatedAt.Time,
// 		FirstPublishedAt: row.PublishedAt.Time,
// 		Title:            row.Title,
// 		Excerpt:          row.Excerpt,
// 		MetaTitle:        row.MetaTitle,
// 		MetaDescription:  row.MetaDescription,
// 		Slug:             row.Slug,
// 		ImageLink:        row.ImageLink.String,
// 		Content:          row.Content.String,
// 		ReadTime:         row.ReadTime.Int32,
// 		Tags:             tags,
// 	}, nil
// }

func DeleteArticle(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.Stmts.DeleteArticle(ctx, dbtx, id)
}

// CountPublishedArticles returns the total count of published articles
func CountPublishedArticles(
	ctx context.Context,
	dbtx db.DBTX,
) (int64, error) {
	return db.Stmts.CountPublishedArticles(ctx, dbtx)
}

// CountDraftArticles returns the total count of draft articles
func CountDraftArticles(
	ctx context.Context,
	dbtx db.DBTX,
) (int64, error) {
	return db.Stmts.CountDraftArticles(ctx, dbtx)
}
