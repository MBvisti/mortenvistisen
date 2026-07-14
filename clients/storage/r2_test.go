package objectstorage

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type fakeR2API struct {
	putInput    *s3.PutObjectInput
	putError    error
	deleteInput *s3.DeleteObjectInput
}

func (f *fakeR2API) PutObject(
	_ context.Context,
	input *s3.PutObjectInput,
	_ ...func(*s3.Options),
) (*s3.PutObjectOutput, error) {
	f.putInput = input
	return &s3.PutObjectOutput{}, f.putError
}

func (f *fakeR2API) DeleteObject(
	_ context.Context,
	input *s3.DeleteObjectInput,
	_ ...func(*s3.Options),
) (*s3.DeleteObjectOutput, error) {
	f.deleteInput = input
	return &s3.DeleteObjectOutput{}, nil
}

func testR2(t *testing.T, api r2API) *R2 {
	t.Helper()
	baseURL, err := url.Parse("https://media.mortenvistisen.com")
	if err != nil {
		t.Fatal(err)
	}

	return &R2{
		client:        api,
		bucket:        "mbv-blog",
		coverPrefix:   "covers",
		publicBaseURL: baseURL,
	}
}

func TestR2UploadCoverUsesConditionalCreate(t *testing.T) {
	api := &fakeR2API{}
	client := testR2(t, api)
	body := "image bytes"

	object, err := client.UploadCover(context.Background(), CoverUpload{
		Filename:    "One Month With Claude.png",
		ContentType: "image/png",
		Size:        int64(len(body)),
		Body:        strings.NewReader(body),
	})
	if err != nil {
		t.Fatalf("UploadCover() error = %v", err)
	}
	if object.Key != "covers/One Month With Claude.png" {
		t.Fatalf("UploadCover() key = %q", object.Key)
	}
	if object.URL != "https://media.mortenvistisen.com/covers/One%20Month%20With%20Claude.png" {
		t.Fatalf("UploadCover() URL = %q", object.URL)
	}
	if aws.ToString(api.putInput.Bucket) != "mbv-blog" ||
		aws.ToString(api.putInput.Key) != object.Key {
		t.Fatalf("PutObject input = %#v", api.putInput)
	}
	if aws.ToString(api.putInput.IfNoneMatch) != "*" {
		t.Fatalf("PutObject IfNoneMatch = %q", aws.ToString(api.putInput.IfNoneMatch))
	}
}

func TestR2UploadCoverMapsCollision(t *testing.T) {
	api := &fakeR2API{putError: &smithy.GenericAPIError{
		Code:    "PreconditionFailed",
		Message: "already exists",
	}}
	client := testR2(t, api)

	_, err := client.UploadCover(context.Background(), CoverUpload{
		Filename:    "existing.png",
		ContentType: "image/png",
		Size:        1,
		Body:        strings.NewReader("x"),
	})
	if !errors.Is(err, ErrObjectExists) {
		t.Fatalf("UploadCover() error = %v, want ErrObjectExists", err)
	}
}

func TestR2CoverKeyOnlyAcceptsManagedCoverURLs(t *testing.T) {
	client := testR2(t, &fakeR2API{})

	key, ok := client.CoverKey(
		"https://media.mortenvistisen.com/covers/One%20Month%20With%20Claude.png",
	)
	if !ok || key != "covers/One Month With Claude.png" {
		t.Fatalf("CoverKey() = %q, %t", key, ok)
	}

	rejected := []string{
		"https://example.com/covers/image.png",
		"https://media.mortenvistisen.com/other/image.png",
		"https://media.mortenvistisen.com/covers/nested/image.png",
		"https://media.mortenvistisen.com/covers/image.png?unsafe=true",
	}
	for _, rawURL := range rejected {
		if key, ok := client.CoverKey(rawURL); ok {
			t.Fatalf("CoverKey(%q) = %q, true", rawURL, key)
		}
	}
}
