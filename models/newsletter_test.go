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

func TestCreateNewsletter(t *testing.T) {
	tests := map[string]struct {
		newsletter models.Newsletter
		expected   []error
	}{
		"should create a new newsletter without errors": {
			newsletter: models.Newsletter{
				Title:       "Test Newsletter",
				Edition:     1,
				Paragraphs:  []string{"a paragrah"},
				ArticleSlug: "/test-newsletter",
			},
			expected: nil,
		},
		"should return errors ErrIsRequired, ErrTooShort, ErrIsRequired": {
			newsletter: models.Newsletter{
				Title:       "empty",
				Edition:     0,
				Paragraphs:  []string{"a paragrah"},
				ArticleSlug: "",
			},
			expected: []error{validation.ErrIsRequired},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualErr := validation.Validate(test.newsletter, models.CreateNewsletterValidations())

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

func TestReleaseNewsletter(t *testing.T) {
	tests := map[string]struct {
		newsleter models.Newsletter
		expected  []error
	}{
		"should release without errors": {
			newsleter: models.Newsletter{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       "All about unit testing",
				Edition:     32,
				ReleasedAt:  time.Now(),
				Released:    true,
				Paragraphs:  []string{"Welcome to this edition about unit testing"},
				ArticleSlug: "/posts/unit-testing-in-golang",
			},
			expected: nil,
		},
		"should return errors ErrTooShort for title and paragrapsh ": {
			newsleter: models.Newsletter{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       "Al",
				Edition:     32,
				ReleasedAt:  time.Now(),
				Paragraphs:  []string{""},
				ArticleSlug: "/posts/unit-testing-in-golang",
			},
			expected: []error{
				validation.ErrTooShort,
				validation.ErrIsRequired,
				validation.ErrMustBeTrue,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualErr := test.newsleter.CanBeReleased()

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
