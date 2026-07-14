package queue

import (
	"context"
	"errors"
	"testing"

	"github.com/riverqueue/river"

	objectstorage "mortenvistisen/clients/storage"
	"mortenvistisen/queue/jobs"
)

type recordingCoverStore struct {
	deletedKey string
	deleteErr  error
}

func (*recordingCoverStore) CoverURL(string) (string, error) {
	return "", errors.New("not implemented")
}

func (*recordingCoverStore) CoverKey(string) (string, bool) {
	return "", false
}

func (*recordingCoverStore) UploadCover(
	context.Context,
	objectstorage.CoverUpload,
) (objectstorage.CoverObject, error) {
	return objectstorage.CoverObject{}, errors.New("not implemented")
}

func (s *recordingCoverStore) Delete(_ context.Context, key string) error {
	s.deletedKey = key
	return s.deleteErr
}

func TestDeleteArticleCoverWorkerDeletesRequestedKey(t *testing.T) {
	store := &recordingCoverStore{}
	worker := NewDeleteArticleCoverWorker(store)

	err := worker.Work(context.Background(), &river.Job[jobs.DeleteArticleCoverArgs]{
		Args: jobs.DeleteArticleCoverArgs{
			ArticleID: 42,
			Key:       "covers/article.png",
		},
	})
	if err != nil {
		t.Fatalf("Work() error = %v", err)
	}
	if store.deletedKey != "covers/article.png" {
		t.Fatalf("Work() deleted key = %q", store.deletedKey)
	}
}

func TestDeleteArticleCoverWorkerReturnsStorageErrorForRetry(t *testing.T) {
	store := &recordingCoverStore{deleteErr: errors.New("R2 unavailable")}
	worker := NewDeleteArticleCoverWorker(store)

	err := worker.Work(context.Background(), &river.Job[jobs.DeleteArticleCoverArgs]{
		Args: jobs.DeleteArticleCoverArgs{ArticleID: 42, Key: "covers/article.png"},
	})
	if err == nil {
		t.Fatal("Work() error = nil, want retryable storage error")
	}
}
