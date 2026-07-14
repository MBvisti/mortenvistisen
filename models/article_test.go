package models

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"mortenvistisen/internal/validation"
)

func validArticleEntity() ArticleEntity {
	return ArticleEntity{
		Title: strings.Repeat("a", articleTitleMaxLength),
		Slug:  "valid-article",
	}
}

func TestArticleEntityValidate(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*ArticleEntity)
		wantField string
		wantCode  string
	}{
		{
			name: "valid draft",
		},
		{
			name: "draft allows empty optional URL",
			mutate: func(article *ArticleEntity) {
				article.ImageLink = normalizeOptionalString(sql.NullString{Valid: true})
			},
		},
		{
			name: "title is required",
			mutate: func(article *ArticleEntity) {
				article.Title = " "
			},
			wantField: "title",
			wantCode:  "required",
		},
		{
			name: "title is too long",
			mutate: func(article *ArticleEntity) {
				article.Title = strings.Repeat("a", articleTitleMaxLength+1)
			},
			wantField: "title",
			wantCode:  "max",
		},
		{
			name: "slug is required",
			mutate: func(article *ArticleEntity) {
				article.Slug = ""
			},
			wantField: "slug",
			wantCode:  "required",
		},
		{
			name: "meta title is too long",
			mutate: func(article *ArticleEntity) {
				article.MetaTitle = sql.NullString{
					String: strings.Repeat("a", articleMetaTitleMaxLength+1),
					Valid:  true,
				}
			},
			wantField: "metaTitle",
			wantCode:  "max",
		},
		{
			name: "meta description is too long",
			mutate: func(article *ArticleEntity) {
				article.MetaDescription = sql.NullString{
					String: strings.Repeat("a", articleMetaDescriptionMaxLength+1),
					Valid:  true,
				}
			},
			wantField: "metaDescription",
			wantCode:  "max",
		},
		{
			name: "image link is invalid",
			mutate: func(article *ArticleEntity) {
				article.ImageLink = sql.NullString{String: "not-a-url", Valid: true}
			},
			wantField: "imageLink",
			wantCode:  "url",
		},
		{
			name: "social image link is invalid",
			mutate: func(article *ArticleEntity) {
				article.MetaImageLink = sql.NullString{String: "not-a-url", Valid: true}
			},
			wantField: "metaImageLink",
			wantCode:  "url",
		},
		{
			name: "read time is negative",
			mutate: func(article *ArticleEntity) {
				article.ReadTime = sql.NullInt32{Int32: -1, Valid: true}
			},
			wantField: "readTime",
			wantCode:  "min",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := validArticleEntity()
			if tt.mutate != nil {
				tt.mutate(&article)
			}

			err := article.Validate()
			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}

			errors, ok := validation.As(err)
			if !ok {
				t.Fatalf("Validate() error = %v, want validation errors", err)
			}

			for _, validationErr := range errors {
				if validationErr.Field == tt.wantField && validationErr.Code == tt.wantCode {
					return
				}
			}

			t.Fatalf(
				"Validate() errors = %#v, want field %q with code %q",
				errors,
				tt.wantField,
				tt.wantCode,
			)
		})
	}
}

func TestPublishedArticleRequiresCompleteContent(t *testing.T) {
	validPublishedArticle := func() ArticleEntity {
		return ArticleEntity{
			FirstPublishedAt: sql.NullTime{Time: time.Now(), Valid: true},
			Published:        true,
			Title:            strings.Repeat("a", articleTitleMaxLength),
			Excerpt:          sql.NullString{String: "A useful summary", Valid: true},
			MetaTitle:        sql.NullString{String: "Search title", Valid: true},
			MetaDescription:  sql.NullString{String: "Search description", Valid: true},
			Slug:             "complete-article",
			ImageLink:        sql.NullString{String: "https://example.com/image.jpg", Valid: true},
			ReadTime:         sql.NullInt32{Int32: 5, Valid: true},
			Content:          sql.NullString{String: "# Complete article", Valid: true},
		}
	}

	tests := []struct {
		name      string
		mutate    func(*ArticleEntity)
		wantField string
		wantCode  string
	}{
		{name: "complete article"},
		{
			name: "first publication is required",
			mutate: func(article *ArticleEntity) {
				article.FirstPublishedAt = sql.NullTime{}
			},
			wantField: "firstPublishedAt",
			wantCode:  "required",
		},
		{
			name: "excerpt is required",
			mutate: func(article *ArticleEntity) {
				article.Excerpt = sql.NullString{}
			},
			wantField: "excerpt",
			wantCode:  "required",
		},
		{
			name: "meta title is required",
			mutate: func(article *ArticleEntity) {
				article.MetaTitle = sql.NullString{}
			},
			wantField: "metaTitle",
			wantCode:  "required",
		},
		{
			name: "meta description is required",
			mutate: func(article *ArticleEntity) {
				article.MetaDescription = sql.NullString{}
			},
			wantField: "metaDescription",
			wantCode:  "required",
		},
		{
			name: "image is required",
			mutate: func(article *ArticleEntity) {
				article.ImageLink = sql.NullString{}
			},
			wantField: "imageLink",
			wantCode:  "required",
		},
		{
			name: "read time must be positive",
			mutate: func(article *ArticleEntity) {
				article.ReadTime = sql.NullInt32{Int32: 0, Valid: true}
			},
			wantField: "readTime",
			wantCode:  "min",
		},
		{
			name: "content is required",
			mutate: func(article *ArticleEntity) {
				article.Content = sql.NullString{}
			},
			wantField: "content",
			wantCode:  "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := validPublishedArticle()
			if tt.mutate != nil {
				tt.mutate(&article)
			}

			err := article.Validate()
			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}

			errors, ok := validation.As(err)
			if !ok {
				t.Fatalf("Validate() error = %v, want validation errors", err)
			}
			for _, validationErr := range errors {
				if validationErr.Field == tt.wantField && validationErr.Code == tt.wantCode {
					return
				}
			}

			t.Fatalf(
				"Validate() errors = %#v, want field %q with code %q",
				errors,
				tt.wantField,
				tt.wantCode,
			)
		})
	}
}
