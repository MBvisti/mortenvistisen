package domain

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewsletter(t *testing.T) {
	tests := map[string]struct {
		title       string
		edition     int32
		releasedAt  time.Time
		released    bool
		paragraphs  []string
		articleSlug string
		expected    error
	}{
		"should create a new newsletter without errors": {
			title:       "Test Newsletter",
			edition:     1,
			releasedAt:  time.Now().Add(5 * time.Minute),
			released:    true,
			paragraphs:  []string{"a paragrah"},
			articleSlug: "/test-newsletter",
			expected:    nil,
		},
		"should return errors ErrIsRequired, ErrTooShort, ErrIsRequired": {
			title:       "",
			edition:     0,
			releasedAt:  time.Now().Add(5 * time.Minute),
			released:    true,
			paragraphs:  []string{"a paragrah"},
			articleSlug: "/test-newsletter",
			expected:    errors.Join(ErrIsRequired, ErrTooShort, ErrIsRequired),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, actualErr := NewNewsletter(
				test.title,
				test.edition,
				test.releasedAt,
				test.released,
				test.paragraphs,
				test.articleSlug,
			)

			if actualErr == nil {
				assert.Equal(
					t,
					test.expected,
					nil,
					fmt.Sprintf(
						"error don't match: expected '%v', got '%v'",
						test.expected,
						actualErr,
					),
				)
			}

			var gotErr error

			var valiErrs ValidationErrs
			if errors.As(actualErr, &valiErrs) {
				var innerErr []error
				for _, valiErr := range valiErrs {
					innerErr = append(innerErr, valiErr.Causes()...)
				}

				gotErr = errors.Join(innerErr...)
			}

			assert.Equal(t, test.expected, gotErr,
				fmt.Sprintf(
					"errors don't match: expected '%v', got '%v'",
					test.expected,
					gotErr,
				),
			)
		})
	}
}
