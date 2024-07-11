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

func TestNewUser(t *testing.T) {
	tests := map[string]struct {
		user            models.User
		confirmPassword string
		expected        error
	}{
		"should create a new user without errors": {
			user: models.User{
				ID:             uuid.New(),
				Name:           "Test user",
				CreatedAt:      time.Now(),
				Mail:           "validemail@gmail.com",
				MailVerifiedAt: time.Now(),
				Password:       "dkslajdklsajkdlsajkldsjkl",
			},
			confirmPassword: "dkslajdklsajkdlsajkldsjkl",
			expected:        nil,
		},
		"should return errors ErrIsRequired": {
			user: models.User{
				ID:             uuid.New(),
				CreatedAt:      time.Now(),
				Name:           "Test user",
				Mail:           "validemail@gmail.com",
				MailVerifiedAt: time.Now(),
				Password:       "dkslajdklsajkdlsajkldsjkl",
			},
			confirmPassword: "dkslajddsadsadsadklsajkdlsajkldsjkl",
			expected:        errors.Join(validation.ErrPasswordDontMatch),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualErr := validation.Validate(
				test.user,
				models.CreateUserValidations(test.confirmPassword),
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

			var valiErrs validation.ValidationErrs
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
