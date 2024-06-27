package psql

import (
	"context"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/google/uuid"
)

func (p Postgres) InsertTag(ctx context.Context, data domain.Tag) error {
	if err := p.Queries.InsertTag(ctx, database.InsertTagParams{
		ID:   data.ID,
		Name: data.Name,
	}); err != nil {
		return err
	}

	return nil
}

func (p Postgres) QueryAllTags(ctx context.Context) ([]domain.Tag, error) {
	tags, err := p.Queries.QueryAllTags(ctx)
	if err != nil {
		return nil, err
	}

	var t []domain.Tag
	for _, tag := range tags {
		t = append(t, domain.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return t, nil
}

func (p Postgres) QueryTagsByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Tag, error) {
	tags, err := p.Queries.QueryTagsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var domainTags []domain.Tag
	for _, t := range tags {
		domainTags = append(domainTags, domain.Tag{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	return domainTags, nil
}

func (p Postgres) DeleteTagsFromPost(ctx context.Context, id uuid.UUID) error {
	return p.Queries.DeleteTagsFromPost(ctx, id)
}
