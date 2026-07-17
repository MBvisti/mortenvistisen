package models

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"mortenvistisen/internal/validation"
)

func validNewsletterEntity() NewsletterEntity {
	return NewsletterEntity{
		Title: strings.Repeat("a", newsletterTitleMaxLength),
		Slug:  sql.NullString{String: "valid-newsletter", Valid: true},
	}
}

func TestNewsletterEntityValidate(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*NewsletterEntity)
		wantField string
		wantCode  string
	}{
		{name: "valid draft"},
		{
			name: "title is required",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.Title = " "
			},
			wantField: "title",
			wantCode:  "required",
		},
		{
			name: "title is too long",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.Title = strings.Repeat("a", newsletterTitleMaxLength+1)
			},
			wantField: "title",
			wantCode:  "max",
		},
		{
			name: "slug is required",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.Slug = sql.NullString{}
			},
			wantField: "slug",
			wantCode:  "required",
		},
		{
			name: "meta title is too long",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.MetaTitle = strings.Repeat("a", newsletterMetaTitleMaxLength+1)
			},
			wantField: "metaTitle",
			wantCode:  "max",
		},
		{
			name: "social image link is invalid",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.MetaImageLink = sql.NullString{String: "not-a-url", Valid: true}
			},
			wantField: "metaImageLink",
			wantCode:  "url",
		},
		{
			name: "meta description is too long",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.MetaDescription = strings.Repeat(
					"a",
					newsletterMetaDescriptionMaxLength+1,
				)
			},
			wantField: "metaDescription",
			wantCode:  "max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newsletter := validNewsletterEntity()
			if tt.mutate != nil {
				tt.mutate(&newsletter)
			}

			err := newsletter.Validate()
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

func TestPublishedNewsletterRequiresCompleteContent(t *testing.T) {
	validPublishedNewsletter := func() NewsletterEntity {
		return NewsletterEntity{
			Title:           strings.Repeat("a", newsletterTitleMaxLength),
			Slug:            sql.NullString{String: "complete-newsletter", Valid: true},
			MetaTitle:       "Search title",
			MetaDescription: "Search description",
			IsPublished:     sql.NullBool{Bool: true, Valid: true},
			ReleasedAt:      sql.NullTime{Time: time.Now(), Valid: true},
			ImageLink: sql.NullString{
				String: "https://media.example.com/cover.png",
				Valid:  true,
			},
			Content: sql.NullString{String: "# Complete newsletter", Valid: true},
		}
	}

	tests := []struct {
		name      string
		mutate    func(*NewsletterEntity)
		wantField string
	}{
		{name: "complete newsletter"},
		{
			name: "release time is required",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.ReleasedAt = sql.NullTime{}
			},
			wantField: "releasedAt",
		},
		{
			name: "meta title is required",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.MetaTitle = ""
			},
			wantField: "metaTitle",
		},
		{
			name: "meta description is required",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.MetaDescription = ""
			},
			wantField: "metaDescription",
		},
		{
			name: "image is required",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.ImageLink = sql.NullString{}
			},
			wantField: "imageLink",
		},
		{
			name: "content is required",
			mutate: func(newsletter *NewsletterEntity) {
				newsletter.Content = sql.NullString{}
			},
			wantField: "content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newsletter := validPublishedNewsletter()
			if tt.mutate != nil {
				tt.mutate(&newsletter)
			}

			err := newsletter.Validate()
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
				if validationErr.Field == tt.wantField && validationErr.Code == "required" {
					return
				}
			}

			t.Fatalf("Validate() errors = %#v, want required error for %q", errors, tt.wantField)
		})
	}
}
