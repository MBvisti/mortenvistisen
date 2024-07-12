package models_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	tests := map[string]struct {
		article  models.Article
		expected []error
	}{
		"should create a new article without errors": {
			article: models.Article{
				ID:          uuid.New(),
				Title:       "Test Article",
				HeaderTitle: "Test Article",
				Filename:    "article.md",
				Slug:        "http://localhost:8080/post/test-article",
				Excerpt:     "this is a test article with an excerpt that is just looong enough to pass these validations that i have built myself by stealing like an artist",
				Draft:       false,
				ReleaseDate: time.Now(),
				ReadTime:    4,
				Tags:        []models.Tag{{uuid.New(), "testing"}},
			},
			expected: nil,
		},
		"should return errors ErrIsRequired ErrTooShort": {
			article: models.Article{
				ID:          uuid.New(),
				Title:       "",
				HeaderTitle: "Test Article",
				Filename:    "article.md",
				Slug:        "http://localhost:8080/post/test-article",
				Excerpt:     "this is a test article with an excerpt that is just looong enough to pass these validations that i have built myself by stealing like an artist",
				Draft:       false,
				ReleaseDate: time.Now(),
				ReadTime:    4,
				Tags:        []models.Tag{{uuid.New(), "testing"}},
			},
			expected: []error{validation.ErrIsRequired, validation.ErrTooShort},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualErr := validation.Validate(test.article, models.CreateArticleValidations())

			if actualErr == nil {
				assert.Equal(
					t,
					len(test.expected),
					0,
					fmt.Sprintf(
						"errors don't match: expected %v, got %v",
						len(test.expected),
						0,
					),
				)
			}

			if actualErr != nil {
				var valiErrs validation.ValidationErrs
				if ok := errors.As(actualErr, &valiErrs); !ok {
					t.Fail()
				}

				assert.Equal(t, test.expected, valiErrs.UnwrapViolations(),
					fmt.Sprintf(
						"errors don't match: expected %v, got %v",
						test.expected,
						valiErrs.UnwrapViolations(),
					),
				)
			}
		})
	}
}
