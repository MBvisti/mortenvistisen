package services

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"mortenvistisen/config"
	"mortenvistisen/internal/validation"
	"mortenvistisen/models"
)

func TestSubscriberConfirmationJob(t *testing.T) {
	subscriber := models.SubscriberEntity{
		ID:    42,
		Email: sql.NullString{String: "reader@example.com", Valid: true},
	}

	args, err := subscriberConfirmationJob(subscriber, "123456")
	if err != nil {
		t.Fatalf("subscriberConfirmationJob() error = %v", err)
	}

	if args.Data.To != subscriber.Email.String {
		t.Fatalf(
			"subscriberConfirmationJob() recipient = %q, want %q",
			args.Data.To,
			subscriber.Email.String,
		)
	}
	if args.Data.Subject != "Confirm your email address" {
		t.Fatalf("subscriberConfirmationJob() subject = %q", args.Data.Subject)
	}
	if args.Data.Metadata["subscriber_id"] != "42" {
		t.Fatalf("subscriberConfirmationJob() metadata = %#v", args.Data.Metadata)
	}

	wantURL := config.BaseURL + "/subscriptions/confirmation"
	for name, body := range map[string]string{
		"HTML": args.Data.HTMLBody,
		"text": args.Data.TextBody,
	} {
		if !strings.Contains(body, "123456") {
			t.Errorf("%s body does not contain the confirmation code", name)
		}
		if !strings.Contains(body, wantURL) {
			t.Errorf("%s body does not contain confirmation URL %q", name, wantURL)
		}
	}
}

func TestConfirmSubscriberEmailValidatesCodeBeforeDatabaseAccess(t *testing.T) {
	service := SubscriberConfirmation{}

	_, err := service.Confirm(t.Context(), ConfirmSubscriberEmailData{Code: "short"})
	validationErrors, ok := validation.As(err)
	if !ok {
		t.Fatalf("Confirm() error = %v, want validation errors", err)
	}

	for _, validationErr := range validationErrors {
		if validationErr.Field == "code" {
			return
		}
	}
	t.Fatalf("Confirm() validation errors = %#v, want code error", validationErrors)
}

func TestSubscriberCodeLifetime(t *testing.T) {
	if subscriberCodeLifetime != 24*time.Hour {
		t.Fatalf("subscriberCodeLifetime = %s, want 24h", subscriberCodeLifetime)
	}
}
