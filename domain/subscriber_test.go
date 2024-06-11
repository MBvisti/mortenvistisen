package domain

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSubscriber(t *testing.T) {
	tests := map[string]struct {
		email        string
		subscribedAt time.Time
		referer      string
		isVerified   bool
		expected     error
	}{
		"should create a new subscriber without errors": {
			email:        "validemail@gmail.com",
			subscribedAt: time.Now().Add(5 * time.Minute),
			referer:      "awesome-post",
			isVerified:   true,
			expected:     nil,
		},
		"should return errors ErrIsRequired ErrInvalidEmail": {
			email:        "",
			subscribedAt: time.Now().Add(5 * time.Minute),
			referer:      "awesome-post",
			isVerified:   true,
			expected:     errors.Join(ErrIsRequired, ErrInvalidEmail),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, actualErr := NewSubscriber(
				test.email,
				test.referer,
				test.subscribedAt,
				test.isVerified,
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
