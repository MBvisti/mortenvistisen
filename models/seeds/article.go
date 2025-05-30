package seeds

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/mbvisti/mortenvistisen/models"
)

type articleSeedData struct {
	ID               uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	FirstPublishedAt time.Time
	IsPublised       bool
	Title            string
	Excerpt          string
	MetaTitle        string
	MetaDescription  string
	Slug             string
	ImageLink        string
	Content          string
	ReadTime         int32
	Tags             []models.ArticleTag
}

type articleSeedOption func(*articleSeedData)

func WithArticleID(id uuid.UUID) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.ID = id
	}
}

func WithArticleCreatedAt(createdAt time.Time) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.CreatedAt = createdAt
	}
}

func WithArticleUpdatedAt(updatedAt time.Time) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.UpdatedAt = updatedAt
	}
}

func WithArticlePublishedAt(publishedAt time.Time) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.FirstPublishedAt = publishedAt
	}
}

func WithArticleTitle(title string) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.Title = title
		// Auto-generate slug from title if not set
		if asd.Slug == "" {
			asd.Slug = slug.Make(title)
		}
		// Auto-generate meta title if not set
		if asd.MetaTitle == "" {
			asd.MetaTitle = title
		}
	}
}

func WithArticleExcerpt(excerpt string) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.Excerpt = excerpt
	}
}

func WithArticleMetaTitle(metaTitle string) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.MetaTitle = metaTitle
	}
}

func WithArticleMetaDescription(metaDescription string) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.MetaDescription = metaDescription
	}
}

func WithArticleSlug(slug string) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.Slug = slug
	}
}

func WithArticleImageLink(imageLink string) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.ImageLink = imageLink
	}
}

func WithArticleContent(content string) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.Content = content
	}
}

func WithPublishedArticle() articleSeedOption {
	return func(asd *articleSeedData) {
		publishedAt := asd.CreatedAt.Add(
			time.Duration(
				rand.IntN(24),
			) * time.Hour, //nolint:gosec // G404: Weak random for test data is acceptable
		)
		asd.FirstPublishedAt = publishedAt
		asd.IsPublised = true
	}
}

func WithDraftArticle() articleSeedOption {
	return func(asd *articleSeedData) {
		asd.FirstPublishedAt = time.Time{}
	}
}

func WithArticleTags(tags []models.ArticleTag) articleSeedOption {
	return func(asd *articleSeedData) {
		asd.Tags = tags
	}
}

func (s Seeder) PlantArticle(
	ctx context.Context,
	opts ...articleSeedOption,
) (models.Article, error) {
	title := generateTechTitle()
	data := &articleSeedData{
		ID: uuid.New(),
		CreatedAt: time.Now().
			Add(-time.Duration(rand.IntN(365)) * 24 * time.Hour),
		//nolint:gosec // G404: Weak random for test data is acceptable
		// Random date in past year
		UpdatedAt:       time.Now(),
		Title:           title,
		Excerpt:         generateExcerpt(),
		MetaTitle:       title,
		MetaDescription: generateMetaDescription(),
		Slug:            slug.Make(title),
		Content:         generateContent(),
		ReadTime:        1,
	}

	// 70% chance of being published
	if rand.Float32() < 0.7 { //nolint:gosec // G404: Weak random for test data is acceptable
		publishedAt := data.CreatedAt.Add(
			time.Duration(
				rand.IntN(24),
			) * time.Hour, //nolint:gosec // G404: Weak random for test data is acceptable
		)
		data.FirstPublishedAt = publishedAt
		data.IsPublised = true
	}

	for _, opt := range opts {
		opt(data)
	}

	article, err := models.NewArticle(ctx, s.dbtx, models.NewArticlePayload{
		Title:           data.Title,
		Excerpt:         data.Excerpt,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            data.Slug,
		ImageLink:       data.ImageLink,
		Content:         data.Content,
		ReadTime:        1,
	})
	if err != nil {
		return models.Article{}, err
	}

	if !data.FirstPublishedAt.IsZero() && data.IsPublised {
		article, err = models.PublishArticle(
			ctx,
			s.dbtx,
			models.PublishArticlePayload{
				ID:  article.ID,
				Now: data.FirstPublishedAt,
			},
		)
		if err != nil {
			return models.Article{}, err
		}
	}

	// Assign tags if provided
	if len(data.Tags) > 0 {
		tagIDs := make([]uuid.UUID, len(data.Tags))
		for i, tag := range data.Tags {
			tagIDs[i] = tag.ID
		}
		_, err = s.PlantArticleTagConnections(ctx, article.ID, tagIDs)
		if err != nil {
			return models.Article{}, err
		}
	}

	return article, nil
}

func (s Seeder) PlantArticles(
	ctx context.Context,
	amount int,
) ([]models.Article, error) {
	articles := make([]models.Article, amount)

	for i := range amount {
		article, err := s.PlantArticle(ctx)
		if err != nil {
			return nil, err
		}

		articles[i] = article
	}

	return articles, nil
}

func (s Seeder) PlantArticlesWithRandomTags(
	ctx context.Context,
	amount int,
	availableTags []models.ArticleTag,
	minTags, maxTags int,
) ([]models.Article, error) {
	articles := make([]models.Article, amount)

	for i := range amount {
		article, err := s.PlantArticle(ctx)
		if err != nil {
			return nil, err
		}

		// Assign random tags
		if len(availableTags) > 0 {
			_, err = s.PlantRandomArticleTagConnections(
				ctx,
				article.ID,
				availableTags,
				minTags,
				maxTags,
			)
			if err != nil {
				return nil, err
			}
		}

		articles[i] = article
	}

	return articles, nil
}

// Helper functions for generating realistic content

func generateTechTitle() string {
	prefixes := []string{
		"How to Build",
		"Getting Started with",
		"A Complete Guide to",
		"Understanding",
		"Mastering",
		"Introduction to",
		"Advanced",
		"Building",
		"Creating",
		"Implementing",
		"Deploying",
		"Optimizing",
		"Debugging",
		"Testing",
		"Scaling",
	}

	topics := []string{
		"Go Applications",
		"REST APIs",
		"Microservices",
		"Docker Containers",
		"Kubernetes Clusters",
		"PostgreSQL Databases",
		"React Components",
		"TypeScript Projects",
		"CI/CD Pipelines",
		"Authentication Systems",
		"Web Servers",
		"GraphQL APIs",
		"Serverless Functions",
		"Database Migrations",
		"Unit Tests",
		"Integration Tests",
		"Performance Monitoring",
		"Error Handling",
		"Logging Systems",
		"Caching Strategies",
	}

	suffixes := []string{
		"",
		"in 2024",
		"with Best Practices",
		"from Scratch",
		"Step by Step",
		"for Beginners",
		"for Production",
		"at Scale",
		"with Examples",
		"the Right Way",
	}

	prefix := prefixes[rand.IntN(len(prefixes))] //nolint:gosec // G404: Weak random for test data is acceptable
	topic := topics[rand.IntN(len(topics))]      //nolint:gosec // G404: Weak random for test data is acceptable
	suffix := suffixes[rand.IntN(len(suffixes))] //nolint:gosec // G404: Weak random for test data is acceptable

	if suffix != "" {
		return fmt.Sprintf("%s %s %s", prefix, topic, suffix)
	}
	return fmt.Sprintf("%s %s", prefix, topic)
}

func generateExcerpt() string {
	excerpts := []string{
		"Learn how to build scalable applications with modern development practices and industry-standard.",
		"A comprehensive guide covering everything you need to know to get started with professional.",
		"Discover best practices and common pitfalls to avoid when building production-ready systems.",
		"Step-by-step tutorial with practical examples and real-world use cases.",
		"Master the fundamentals and advanced concepts with hands-on examples and detailed explanations.",
		"Explore modern development techniques and learn how to implement them in your projects.",
		"Complete walkthrough from basic concepts to advanced implementation strategies.",
		"Learn industry best practices and how to apply them to your development workflow.",
		"Practical guide with code examples and detailed explanations of key concepts.",
		"Everything you need to know to build robust, maintainable, and scalable solutions.",
	}
	return excerpts[rand.IntN(len(excerpts))] //nolint:gosec // G404: Weak random for test data is acceptable
}

func generateMetaDescription() string {
	descriptions := []string{
		"Learn essential development skills with practical examples and best practices.",
		"Comprehensive tutorial covering modern development techniques and industry standards.",
		"Step-by-step guide with real-world examples. Master the tools and techniques used by professional.",
		"Practical development guide with code examples and detailed explanations.",
		"Complete tutorial covering everything from basics to advanced concepts.",
		"Learn modern development practices with hands-on examples.",
		"Detailed guide covering best practices and common patterns.",
		"Professional development tutorial with practical examples.",
		"Comprehensive guide to building robust applications.",
		"Master essential development concepts with practical examples and real-world use cases.",
	}
	return descriptions[rand.IntN(len(descriptions))] //nolint:gosec // G404: Weak random for test data is acceptable
}

func generateContent() string {
	return fmt.Sprintf(
		`# Introduction

%s

## Getting Started

%s

### Prerequisites

Before we begin, make sure you have the following installed:

- Go 1.21 or later
- PostgreSQL 14+
- Docker (optional)
- Git

### Installation

First, let's set up our development environment:

%s

## Implementation

%s

### Step 1: Project Setup

%s

### Step 2: Configuration

%s

### Step 3: Implementation

%s

## Best Practices

%s

## Conclusion

%s

Happy coding!`,
		faker.Paragraph(),
		faker.Paragraph(),
		generateCodeBlock(
			"bash",
			"go mod init example.com/myproject\ngo get github.com/labstack/echo/v4",
		),
		faker.Paragraph(),
		faker.Paragraph(),
		faker.Paragraph(),
		faker.Paragraph(),
		faker.Paragraph(),
		faker.Paragraph(),
	)
}

func generateCodeBlock(language, code string) string {
	return fmt.Sprintf("```%s\n%s\n```", language, code)
}
