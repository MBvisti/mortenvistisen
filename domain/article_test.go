package domain

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	tests := map[string]struct {
		title       string
		headerTitle string
		filename    string
		slug        string
		excerpt     string
		draft       bool
		releaseDate time.Time
		readTime    int32
		tags        []Tag
		expected    error
	}{
		"should create a new article without errors": {
			title:       "Test Article",
			headerTitle: "Test Article",
			filename:    "article.md",
			slug:        "http://localhost:8080/post/test-article",
			excerpt:     "this is a test article with an excerpt that is just looong enough to pass these validations that i have built myself by stealing like an artist",
			draft:       false,
			releaseDate: time.Now(),
			readTime:    4,
			tags:        []Tag{{uuid.New(), "testing"}},
			expected:    nil,
		},
		"should return errors ErrIsRequired ErrTooShort": {
			title:       "",
			headerTitle: "Test Article",
			filename:    "article.md",
			slug:        "http://localhost:8080/post/test-article",
			excerpt:     "this is a test article with an excerpt that is just looong enough to pass these validations that i have built myself by stealing like an artist",
			draft:       false,
			releaseDate: time.Now(),
			readTime:    4,
			tags:        []Tag{{uuid.New(), "testing"}},
			expected:    errors.Join(ErrIsRequired, ErrTooShort),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, actualErr := NewArticle(
				test.title,
				test.headerTitle,
				test.filename,
				test.slug,
				test.excerpt,
				test.readTime,
				test.tags,
			)

			if actualErr == nil {
				assert.Equal(
					t,
					test.expected,
					nil,
					fmt.Sprintf(
						"errors don't match: expected %v, got %v",
						test.expected,
						actualErr,
					),
				)
			}

			var gotErr error

			var valiErrs ValidationErrs
			if errors.As(actualErr, &valiErrs) {
				for _, valiErr := range valiErrs {
					gotErr = errors.Join(valiErr.Causes()...)
				}
			}

			assert.Equal(t, test.expected, gotErr,
				fmt.Sprintf(
					"errors don't match: expected %v, got %v",
					test.expected,
					gotErr,
				),
			)
		})
	}
}
