package models

import (
	"context"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/google/uuid"
)

type tagStorage interface {
	InsertTag(ctx context.Context, data domain.Tag) error
	QueryAllTags(ctx context.Context) ([]domain.Tag, error)
}

type TagService struct {
	storage tagStorage
}

func NewTagSvc(storage tagStorage) TagService {
	return TagService{
		storage,
	}
}

func (t TagService) New(ctx context.Context, name string) (domain.Tag, error) {
	tag := domain.Tag{
		ID:   uuid.New(),
		Name: name,
	}
	if err := t.storage.InsertTag(ctx, tag); err != nil {
		return domain.Tag{}, err
	}

	return tag, nil
}

func (t TagService) All(ctx context.Context) ([]domain.Tag, error) {
	return t.storage.QueryAllTags(ctx)
}
