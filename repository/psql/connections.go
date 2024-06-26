package psql

import (
	"context"
	"errors"

	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/google/uuid"
)

func (p Postgres) AssociateTagsWithPost(
	ctx context.Context,
	postID uuid.UUID,
	tagIDs []uuid.UUID,
) error {
	tx, err := p.tx.Begin(ctx)
	if err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		if err := p.Queries.WithTx(tx).AssociateTagWithPost(ctx, database.AssociateTagWithPostParams{
			ID:     uuid.New(),
			PostID: postID,
			TagID:  tagID,
		}); err != nil {
			return errors.Join(err, tx.Rollback(ctx))
		}
	}

	return tx.Commit(ctx)
}

func (p Postgres) UpdateTagsPostAssociations(
	ctx context.Context,
	postID uuid.UUID,
	tagIDs []uuid.UUID,
) error {
	tx, err := p.tx.Begin(ctx)
	if err != nil {
		return err
	}

	if err := p.Queries.WithTx(tx).DeleteTagsFromPost(ctx, postID); err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		if err := p.Queries.WithTx(tx).AssociateTagWithPost(ctx, database.AssociateTagWithPostParams{
			ID:     uuid.New(),
			PostID: postID,
			TagID:  tagID,
		}); err != nil {
			return errors.Join(err, tx.Rollback(ctx))
		}
	}

	return tx.Commit(ctx)
}
