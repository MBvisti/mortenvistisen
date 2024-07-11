package models_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTagArticle(t *testing.T) {
	tests := map[string]struct {
		tag      models.Tag
		expected []error
	}{
		"should create a new tag without errors": {
			tag: models.Tag{
				ID:   uuid.New(),
				Name: "golang",
			},
			expected: nil,
		},
		"should return errors ErrIsRequired ErrTooShort": {
			tag: models.Tag{
				Name: "G",
			},
			expected: []error{validation.ErrIsRequired, validation.ErrTooShort},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualErr := validation.Validate(test.tag, models.CreateTagValidations())

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
		})
	}
}
