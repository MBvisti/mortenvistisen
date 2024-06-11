package domain

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	tests := map[string]struct {
		name            string
		mail            string
		mailVerifiedAt  time.Time
		password        string
		confirmPassword string
		expected        error
	}{
		"should create a new subscriber without errors": {
			name:            "Test user",
			mail:            "validemail@gmail.com",
			mailVerifiedAt:  time.Now(),
			password:        "dkslajdklsajkdlsajkldsjkl",
			confirmPassword: "dkslajdklsajkdlsajkldsjkl",
			expected:        nil,
		},
		"should return errors ErrIsRequired ErrInvalidEmail": {
			name:            "Test user",
			mail:            "validemail@gmail.com",
			mailVerifiedAt:  time.Now(),
			password:        "dkslajdklsajkdlsajkldsjkl",
			confirmPassword: "dkslajddsadsadsadklsajkdlsajkldsjkl",
			expected:        errors.Join(ErrPasswordDontMatch),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, actualErr := NewUser(
				test.name,
				test.mail,
				test.password,
				test.confirmPassword,
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
