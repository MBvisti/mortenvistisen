package services

import (
	"database/sql"
	"testing"
	"time"
)

func TestNewsletterReleaseState(t *testing.T) {
	now := time.Date(2026, time.July, 9, 12, 0, 0, 0, time.UTC)
	original := time.Date(2026, time.July, 1, 9, 30, 0, 0, time.UTC)

	tests := []struct {
		name         string
		published    bool
		current      sql.NullTime
		want         sql.NullTime
		wantReleased bool
	}{
		{
			name:      "draft remains unreleased",
			published: false,
		},
		{
			name:         "first publication records the release time",
			published:    true,
			want:         sql.NullTime{Time: now, Valid: true},
			wantReleased: true,
		},
		{
			name:      "unpublishing preserves the first release",
			published: false,
			current:   sql.NullTime{Time: original, Valid: true},
			want:      sql.NullTime{Time: original, Valid: true},
		},
		{
			name:      "republishing preserves the first release",
			published: true,
			current:   sql.NullTime{Time: original, Valid: true},
			want:      sql.NullTime{Time: original, Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, released := newsletterReleaseState(tt.published, tt.current, now)
			if got.Valid != tt.want.Valid || !got.Time.Equal(tt.want.Time) {
				t.Fatalf("newsletterReleaseState() time = %#v, want %#v", got, tt.want)
			}
			if released != tt.wantReleased {
				t.Fatalf(
					"newsletterReleaseState() released = %t, want %t",
					released,
					tt.wantReleased,
				)
			}
		})
	}
}
