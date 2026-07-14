package services

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"mortenvistisen/models"
	"mortenvistisen/queue/jobs"
)

func TestFirstPublicationState(t *testing.T) {
	now := time.Date(2026, time.July, 9, 12, 0, 0, 0, time.UTC)
	original := time.Date(2026, time.July, 1, 9, 30, 0, 0, time.UTC)

	tests := []struct {
		name          string
		published     bool
		current       sql.NullTime
		want          sql.NullTime
		wantFirstTime bool
	}{
		{
			name:      "draft remains unpublished",
			published: false,
		},
		{
			name:          "first publication records the time",
			published:     true,
			want:          sql.NullTime{Time: now, Valid: true},
			wantFirstTime: true,
		},
		{
			name:      "unpublishing preserves the first publication",
			published: false,
			current:   sql.NullTime{Time: original, Valid: true},
			want:      sql.NullTime{Time: original, Valid: true},
		},
		{
			name:      "republishing preserves the first publication",
			published: true,
			current:   sql.NullTime{Time: original, Valid: true},
			want:      sql.NullTime{Time: original, Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, firstTime := firstPublicationState(tt.published, tt.current, now)
			if got.Valid != tt.want.Valid || !got.Time.Equal(tt.want.Time) {
				t.Fatalf("firstPublicationState() time = %#v, want %#v", got, tt.want)
			}
			if firstTime != tt.wantFirstTime {
				t.Fatalf(
					"firstPublicationState() firstTime = %t, want %t",
					firstTime,
					tt.wantFirstTime,
				)
			}
		})
	}
}

func TestUniqueArticleTagIDs(t *testing.T) {
	got := uniqueArticleTagIDs([]int32{4, 1, 4, 0, -2, 9, 1})
	want := []int32{4, 1, 9}
	if len(got) != len(want) {
		t.Fatalf("uniqueArticleTagIDs() = %#v, want %#v", got, want)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("uniqueArticleTagIDs() = %#v, want %#v", got, want)
		}
	}
}

func TestFirstPublicationJobs(t *testing.T) {
	article := models.ArticleEntity{
		ID:        42,
		Title:     "Shipping reliable Go services",
		Excerpt:   sql.NullString{String: "A practical guide to reliable delivery.", Valid: true},
		ImageLink: sql.NullString{String: "https://example.com/article.jpg", Valid: true},
	}
	subscribers := []models.SubscriberEntity{
		{ID: 7, Email: sql.NullString{String: "first@example.com", Valid: true}},
		{ID: 8, Email: sql.NullString{String: "second@example.com", Valid: true}},
	}

	params, err := firstPublicationJobs(article, subscribers)
	if err != nil {
		t.Fatalf("firstPublicationJobs() error = %v", err)
	}
	if len(params) != len(subscribers) {
		t.Fatalf("firstPublicationJobs() count = %d, want %d", len(params), len(subscribers))
	}

	for index, param := range params {
		args, ok := param.Args.(jobs.SendMarketingEmailArgs)
		if !ok {
			t.Fatalf("job %d args type = %T, want SendMarketingEmailArgs", index, param.Args)
		}
		if len(args.Data.To) != 1 || args.Data.To[0] != subscribers[index].Email.String {
			t.Fatalf("job %d recipients = %#v", index, args.Data.To)
		}
		if !strings.Contains(args.Data.UnsubscribeURL, articleUnsubscribeTODOPath) {
			t.Fatalf("job %d unsubscribe URL = %q", index, args.Data.UnsubscribeURL)
		}
		if args.Data.Metadata["article_id"] != "42" {
			t.Fatalf("job %d article metadata = %#v", index, args.Data.Metadata)
		}
		if !strings.Contains(args.Data.HTMLBody, article.Title) ||
			!strings.Contains(args.Data.HTMLBody, article.Excerpt.String) {
			t.Fatalf("job %d HTML body does not contain article details", index)
		}
	}
}

func TestApplyArticlePatchPreservesOmittedFields(t *testing.T) {
	current := models.ArticleEntity{
		ID:              42,
		Title:           "Original title",
		Excerpt:         sql.NullString{String: "Original excerpt", Valid: true},
		MetaTitle:       sql.NullString{String: "Original meta title", Valid: true},
		MetaDescription: sql.NullString{String: "Original meta description", Valid: true},
		Slug:            "original-slug",
		ImageLink:       sql.NullString{String: "https://example.com/original.png", Valid: true},
	}
	title := "Updated title"
	metaDescription := "Updated meta description"
	imageLink := "https://example.com/updated.png"

	got := applyArticlePatch(current, ArticlePatch{
		Title:           &title,
		MetaDescription: &metaDescription,
		ImageLink:       &imageLink,
	})

	if got.Title != title || got.MetaDescription.String != metaDescription ||
		got.ImageLink.String != imageLink {
		t.Fatalf("applyArticlePatch() did not apply supplied fields: %#v", got)
	}
	if got.Excerpt != current.Excerpt || got.MetaTitle != current.MetaTitle || got.Slug != current.Slug {
		t.Fatalf("applyArticlePatch() changed omitted fields: %#v", got)
	}
	if !(ArticlePatch{}).Empty() || (ArticlePatch{ImageLink: &imageLink}).Empty() {
		t.Fatal("ArticlePatch.Empty() returned the wrong result")
	}
}
