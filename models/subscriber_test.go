package models

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"mortenvistisen/internal/validation"
)

func validSubscriberEntity() SubscriberEntity {
	return SubscriberEntity{
		Email:        sql.NullString{String: "reader@example.com", Valid: true},
		SubscribedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Referer:      sql.NullString{String: "newsletter campaign", Valid: true},
	}
}

func TestSubscriberEntityValidate(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*SubscriberEntity)
		wantField string
		wantCode  string
	}{
		{name: "valid subscriber"},
		{
			name: "email is required",
			mutate: func(subscriber *SubscriberEntity) {
				subscriber.Email = sql.NullString{}
			},
			wantField: "email",
			wantCode:  "required",
		},
		{
			name: "email must be valid",
			mutate: func(subscriber *SubscriberEntity) {
				subscriber.Email = sql.NullString{String: "not-an-email", Valid: true}
			},
			wantField: "email",
			wantCode:  "email",
		},
		{
			name: "email is too long",
			mutate: func(subscriber *SubscriberEntity) {
				subscriber.Email = sql.NullString{
					String: strings.Repeat("a", subscriberEmailMaxLength+1),
					Valid:  true,
				}
			},
			wantField: "email",
			wantCode:  "max",
		},
		{
			name: "subscription date is required",
			mutate: func(subscriber *SubscriberEntity) {
				subscriber.SubscribedAt = sql.NullTime{}
			},
			wantField: "subscribedAt",
			wantCode:  "required",
		},
		{
			name: "referrer is too long",
			mutate: func(subscriber *SubscriberEntity) {
				subscriber.Referer = sql.NullString{
					String: "https://example.com/" + strings.Repeat(
						"a",
						subscriberRefererMaxLength,
					),
					Valid: true,
				}
			},
			wantField: "referer",
			wantCode:  "max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscriber := validSubscriberEntity()
			if tt.mutate != nil {
				tt.mutate(&subscriber)
			}

			err := subscriber.Validate()
			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}

			errors, ok := validation.As(err)
			if !ok {
				t.Fatalf("Validate() error = %v, want validation errors", err)
			}
			for _, validationErr := range errors {
				if validationErr.Field == tt.wantField && validationErr.Code == tt.wantCode {
					return
				}
			}

			t.Fatalf(
				"Validate() errors = %#v, want field %q with code %q",
				errors,
				tt.wantField,
				tt.wantCode,
			)
		})
	}
}

func TestSubscriberValidationRules(t *testing.T) {
	rules := SubscriberValidationRules()

	expected := map[string][]string{
		"email":        {"required", "max", "email"},
		"subscribedAt": {"required"},
		"referer":      {"max"},
	}

	for field, expectedCodes := range expected {
		fieldRules := rules[field]
		for _, expectedCode := range expectedCodes {
			found := false
			for _, rule := range fieldRules {
				if rule.Code == expectedCode {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("SubscriberValidationRules()[%q] missing %q rule", field, expectedCode)
			}
		}
	}

	for _, rule := range rules["referer"] {
		if rule.Code == "url" {
			t.Error("SubscriberValidationRules()[\"referer\"] includes unexpected URL rule")
		}
	}
}
