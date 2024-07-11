package models

import (
	"context"

	"github.com/google/uuid"
)

type tagStorage interface {
	InsertTag(ctx context.Context, data Tag) error
	QueryAllTags(ctx context.Context) ([]Tag, error)
}

type TagService struct {
	storage tagStorage
}

func NewTagSvc(storage tagStorage) TagService {
	return TagService{
		storage,
	}
}

func (t TagService) New(ctx context.Context, name string) (Tag, error) {
	tag := Tag{
		ID:   uuid.New(),
		Name: name,
	}
	if err := t.storage.InsertTag(ctx, tag); err != nil {
		return Tag{}, err
	}

	return tag, nil
}

func (t TagService) All(ctx context.Context) ([]Tag, error) {
	return t.storage.QueryAllTags(ctx)
}
