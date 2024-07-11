package psql

import (
	"context"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/google/uuid"
)

func (p Postgres) InsertTag(ctx context.Context, data models.Tag) error {
	if err := p.Queries.InsertTag(ctx, database.InsertTagParams{
		ID:   data.ID,
		Name: data.Name,
	}); err != nil {
		return err
	}

	return nil
}

func (p Postgres) QueryAllTags(ctx context.Context) ([]models.Tag, error) {
	tags, err := p.Queries.QueryAllTags(ctx)
	if err != nil {
		return nil, err
	}

	var t []models.Tag
	for _, tag := range tags {
		t = append(t, models.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return t, nil
}

func (p Postgres) QueryTagsByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Tag, error) {
	tags, err := p.Queries.QueryTagsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var modelsTags []models.Tag
	for _, t := range tags {
		modelsTags = append(modelsTags, models.Tag{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	return modelsTags, nil
}

func (p Postgres) DeleteTagsFromPost(ctx context.Context, id uuid.UUID) error {
	return p.Queries.DeleteTagsFromPost(ctx, id)
}
