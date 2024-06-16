package domain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewsletter(t *testing.T) {
	tests := map[string]struct {
		title       string
		edition     int32
		paragraphs  []string
		articleSlug string
		expected    error
	}{
		"should create a new newsletter without errors": {
			title:       "Test Newsletter",
			edition:     1,
			paragraphs:  []string{"a paragrah"},
			articleSlug: "/test-newsletter",
			expected:    nil,
		},
		"should return errors ErrIsRequired, ErrTooShort, ErrIsRequired": {
			title:       "",
			edition:     0,
			paragraphs:  []string{"a paragrah"},
			articleSlug: "/test-newsletter",
			expected:    errors.Join(ErrIsRequired, ErrTooShort, ErrIsRequired),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, actualErr := CreateNewsletter(
				test.title,
				test.edition,
				test.paragraphs,
				test.articleSlug,
			).Release()

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
