package controllers

import (
	"context"
	"strings"
	"testing"

	"mortenvistisen/internal/hypermedia"
	"mortenvistisen/router/routes"
	"mortenvistisen/views"
)

func TestNewsletterIndexSubscriptionForm(t *testing.T) {
	html, err := hypermedia.RenderHTML(context.Background(), views.NewsletterIndex{}.Page())
	if err != nil {
		t.Fatal(err)
	}

	for _, attribute := range []string{
		`action="` + routes.SubscriberCreate.URL() + `"`,
		`method="post"`,
		`type="email"`,
		`name="email"`,
		`required`,
		`maxlength="254"`,
	} {
		if !strings.Contains(html, attribute) {
			t.Fatalf("newsletter subscription form is missing %s", attribute)
		}
	}
}
