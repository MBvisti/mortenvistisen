package services

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"strings"
	"testing"

	objectstorage "mortenvistisen/clients/storage"
	"mortenvistisen/internal/validation"
)

type fakeCoverStore struct{}

func (fakeCoverStore) CoverURL(filename string) (string, error) {
	return "https://media.mortenvistisen.com/covers/" + filename, nil
}

func (fakeCoverStore) CoverKey(rawURL string) (string, bool) {
	const prefix = "https://media.mortenvistisen.com/covers/"
	if !strings.HasPrefix(rawURL, prefix) {
		return "", false
	}

	return "covers/" + strings.TrimPrefix(rawURL, prefix), true
}

func (fakeCoverStore) UploadCover(
	context.Context,
	objectstorage.CoverUpload,
) (objectstorage.CoverObject, error) {
	return objectstorage.CoverObject{}, errors.New("not implemented")
}

func (fakeCoverStore) Delete(context.Context, string) error {
	return errors.New("not implemented")
}

func TestPrepareCoverPreservesFilenameAndBody(t *testing.T) {
	png := append([]byte(nil), []byte("\x89PNG\r\n\x1a\n")...)
	png = append(png, []byte("remaining image bytes")...)
	upload, publicURL, err := prepareCover(fakeCoverStore{}, 10<<20, &Cover{
		Filename: "A cover image.png",
		Size:     int64(len(png)),
		Body:     strings.NewReader(string(png)),
	})
	if err != nil {
		t.Fatalf("prepareCover() error = %v", err)
	}
	if upload.Filename != "A cover image.png" || upload.ContentType != "image/png" {
		t.Fatalf("prepareCover() upload = %#v", upload)
	}
	if publicURL != "https://media.mortenvistisen.com/covers/A cover image.png" {
		t.Fatalf("prepareCover() URL = %q", publicURL)
	}
	gotBody, err := io.ReadAll(upload.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotBody) != string(png) {
		t.Fatal("prepareCover() did not preserve the complete image body")
	}
}

func TestPrepareCoverRejectsMismatchedContent(t *testing.T) {
	_, _, err := prepareCover(fakeCoverStore{}, 10<<20, &Cover{
		Filename: "not-really-an-image.png",
		Size:     10,
		Body:     strings.NewReader("plain text"),
	})
	validationErrors, ok := validation.As(err)
	if !ok || len(validationErrors) != 1 || validationErrors[0].Field != "cover" {
		t.Fatalf("prepareCover() error = %v", err)
	}
}

func TestDetectImageContentTypeFindsCompatibleAVIFBrand(t *testing.T) {
	header := make([]byte, 24)
	binary.BigEndian.PutUint32(header[:4], uint32(len(header)))
	copy(header[4:8], "ftyp")
	copy(header[8:12], "mif1")
	copy(header[16:20], "avif")

	if got := detectImageContentType(header); got != "image/avif" {
		t.Fatalf("detectImageContentType() = %q", got)
	}
}
