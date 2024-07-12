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

func TestNewSubscriber(t *testing.T) {
	tests := map[string]struct {
		subscriber models.Subscriber
		expected   []error
	}{
		"should create a new subscriber without errors": {
			subscriber: models.Subscriber{
				ID:           uuid.New(),
				Email:        "validemail@gmail.com",
				SubscribedAt: time.Now().Add(5 * time.Minute),
				Referer:      "awesome-post",
				IsVerified:   true,
			},
			expected: nil,
		},
		"should return errors ErrIsRequired ErrInvalidEmail": {
			subscriber: models.Subscriber{
				Email:        "",
				SubscribedAt: time.Now().Add(5 * time.Minute),
				Referer:      "awesome-post",
				IsVerified:   true,
			},
			expected: []error{
				validation.ErrIsRequired,
				validation.ErrIsRequired,
				validation.ErrInvalidEmail,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualErr := validation.Validate(test.subscriber, models.CreateSubscriberValidations())

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
