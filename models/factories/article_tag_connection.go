package factories

import (
	"context"
	"fmt"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
)

// ArticleTagConnectionFactory wraps models.ArticleTagConnectionEntity for testing
type ArticleTagConnectionFactory struct {
	models.ArticleTagConnectionEntity
}

type ArticleTagConnectionOption func(*ArticleTagConnectionFactory)

// BuildArticleTagConnection creates an in-memory ArticleTagConnection with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateArticleTagConnection.
func BuildArticleTagConnection(
	articleID int32,
	tagID int32,
	opts ...ArticleTagConnectionOption,
) models.ArticleTagConnectionEntity {
	f := &ArticleTagConnectionFactory{
		ArticleTagConnectionEntity: models.ArticleTagConnectionEntity{
			ArticleID: articleID,
			TagID:     tagID,
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.ArticleTagConnectionEntity
}

// CreateArticleTagConnection creates and persists a ArticleTagConnection to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateArticleTagConnection(
	ctx context.Context,
	exec storage.Executor,
	articleID int32,
	tagID int32,
	opts ...ArticleTagConnectionOption,
) (models.ArticleTagConnectionEntity, error) {
	built := BuildArticleTagConnection(articleID, tagID, opts...)

	entity := models.ArticleTagConnectionEntity{
		ArticleID: built.ArticleID,
		TagID:     built.TagID,
	}

	if err := exec.NewInsert().Model(&entity).Returning("*").Scan(ctx); err != nil {
		return models.ArticleTagConnectionEntity{}, err
	}

	return entity, nil
}

// CreateArticleTagConnections creates multiple ArticleTagConnection records at once
func CreateArticleTagConnections(
	ctx context.Context,
	exec storage.Executor,
	articleID int32,
	tagID int32,
	count int,
	opts ...ArticleTagConnectionOption,
) ([]models.ArticleTagConnectionEntity, error) {
	articletagconnections := make([]models.ArticleTagConnectionEntity, 0, count)

	for i := range count {
		entity, err := CreateArticleTagConnection(ctx, exec, articleID, tagID, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create articletagconnection %d: %w", i+1, err)
		}
		articletagconnections = append(articletagconnections, entity)
	}

	return articletagconnections, nil
}

// Option functions

// WithArticleTagConnectionsArticleID sets the ArticleID field
func WithArticleTagConnectionsArticleID(value int32) ArticleTagConnectionOption {
	return func(f *ArticleTagConnectionFactory) {
		f.ArticleTagConnectionEntity.ArticleID = value
	}
}

// WithArticleTagConnectionsTagID sets the TagID field
func WithArticleTagConnectionsTagID(value int32) ArticleTagConnectionOption {
	return func(f *ArticleTagConnectionFactory) {
		f.ArticleTagConnectionEntity.TagID = value
	}
}
