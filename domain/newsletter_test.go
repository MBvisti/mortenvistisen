package domain

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
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
			title:       "empty",
			edition:     0,
			paragraphs:  []string{"a paragrah"},
			articleSlug: "",
			expected:    errors.Join(ErrIsRequired),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, actualErr := CreateNewsletter(
				test.title,
				test.edition,
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

func TestReleaseNewsletter(t *testing.T) {
	tests := map[string]struct {
		newsleter Newsletter
		expected  error
	}{
		"should release without errors": {
			newsleter: Newsletter{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       "All about unit testing",
				Edition:     32,
				ReleasedAt:  time.Now(),
				Paragraphs:  []string{"Welcome to this edition about unit testing"},
				ArticleSlug: "/posts/unit-testing-in-golang",
			},
			expected: nil,
		},
		"should return errors ErrTooShort for title and paragrapsh ": {
			newsleter: Newsletter{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       "Al",
				Edition:     32,
				ReleasedAt:  time.Now(),
				Paragraphs:  []string{""},
				ArticleSlug: "/posts/unit-testing-in-golang",
			},
			expected: errors.Join(ErrTooShort),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, actualErr := test.newsleter.Release()

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
