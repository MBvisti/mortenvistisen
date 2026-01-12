package factories

import (
	"context"
	"fmt"
	"time"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

// ArticleFactory wraps models.Article for testing
type ArticleFactory struct {
	models.Article // Embedded
	tagIDs         []uuid.UUID
}

type ArticleOption func(*ArticleFactory)

// BuildArticle creates an in-memory Article with default test values
func BuildArticle(opts ...ArticleOption) models.Article {
	f := &ArticleFactory{
		Article: models.Article{
			ID:               uuid.New(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			FirstPublishedAt: time.Time{}, // Optional timestamp - zero by default
			Title:            faker.Word(),
			Excerpt:          faker.Word(),
			MetaTitle:        faker.Word(),
			MetaDescription:  faker.Word(),
			Slug:             faker.Word(),
			ImageLink:        faker.Word(),
			ReadTime:         randomInt(1, 1000, 100),
			Content:          faker.Word(),
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.Article
}

// CreateArticle creates and persists a Article to the database
func CreateArticle(ctx context.Context, exec storage.Executor, opts ...ArticleOption) (models.Article, error) {
	f := &ArticleFactory{
		Article: BuildArticle(opts...),
	}

	for _, opt := range opts {
		opt(f)
	}

	// Prepare creation data
	data := models.CreateArticleData{
		FirstPublishedAt: f.Article.FirstPublishedAt,
		Title:            f.Article.Title,
		Excerpt:          f.Article.Excerpt,
		MetaTitle:        f.Article.MetaTitle,
		MetaDescription:  f.Article.MetaDescription,
		Slug:             f.Article.Slug,
		ImageLink:        f.Article.ImageLink,
		ReadTime:         f.Article.ReadTime,
		Content:          f.Article.Content,
	}

	// Use model's Create function
	article, err := models.CreateArticle(ctx, exec, data)
	if err != nil {
		return models.Article{}, err
	}

	// Associate tags if provided
	if len(f.tagIDs) > 0 {
		if err := models.AssociateTagsWithArticle(ctx, exec, article.ID, f.tagIDs); err != nil {
			return models.Article{}, err
		}
	}

	return article, nil
}

// CreateArticles creates multiple Article records at once
func CreateArticles(ctx context.Context, exec storage.Executor, count int, opts ...ArticleOption) ([]models.Article, error) {
	articles := make([]models.Article, 0, count)

	for i := 0; i < count; i++ {
		article, err := CreateArticle(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create article %d: %w", i+1, err)
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// Option functions

// WithArticlesFirstPublishedAt sets the FirstPublishedAt field
func WithArticlesFirstPublishedAt(value time.Time) ArticleOption {
	return func(f *ArticleFactory) {
		f.Article.FirstPublishedAt = value
	}
}

// WithArticlesTitle sets the Title field
func WithArticlesTitle(value string) ArticleOption {
	return func(f *ArticleFactory) {
		f.Article.Title = value
	}
}

// WithArticlesExcerpt sets the Excerpt field
func WithArticlesExcerpt(value string) ArticleOption {
	return func(f *ArticleFactory) {
		f.Article.Excerpt = value
	}
}

// WithArticlesMetaTitle sets the MetaTitle field
func WithArticlesMetaTitle(value string) ArticleOption {
	return func(f *ArticleFactory) {
		f.Article.MetaTitle = value
	}
}

// WithArticlesMetaDescription sets the MetaDescription field
func WithArticlesMetaDescription(value string) ArticleOption {
	return func(f *ArticleFactory) {
		f.Article.MetaDescription = value
	}
}

// WithArticlesSlug sets the Slug field
func WithArticlesSlug(value string) ArticleOption {
	return func(f *ArticleFactory) {
		f.Article.Slug = value
	}
}

// WithArticlesImageLink sets the ImageLink field
func WithArticlesImageLink(value string) ArticleOption {
	return func(f *ArticleFactory) {
		f.Article.ImageLink = value
	}
}

// WithArticlesReadTime sets the ReadTime field
func WithArticlesReadTime(value int32) ArticleOption {
	return func(f *ArticleFactory) {
		f.Article.ReadTime = value
	}
}

// WithArticlesContent sets the Content field
func WithArticlesContent(value string) ArticleOption {
	return func(f *ArticleFactory) {
		f.Article.Content = value
	}
}

// WithArticlesTags sets the tags to associate with the article
func WithArticlesTags(tagIDs []uuid.UUID) ArticleOption {
	return func(f *ArticleFactory) {
		f.tagIDs = tagIDs
	}
}
