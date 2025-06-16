package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/mbvisti/mortenvistisen/models/internal/db"
)

type ArticleTagConnection struct {
	ID        uuid.UUID
	ArticleID uuid.UUID
	TagID     uuid.UUID
}

func GetArticleTagConnectionByID(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (ArticleTagConnection, error) {
	row, err := db.Stmts.QueryArticleTagConnectionByID(ctx, dbtx, id)
	if err != nil {
		return ArticleTagConnection{}, err
	}

	return ArticleTagConnection{
		ID:        row.ID,
		ArticleID: row.ArticleID,
		TagID:     row.TagID,
	}, nil
}

func GetArticleTagConnectionsByArticleID(
	ctx context.Context,
	dbtx db.DBTX,
	articleID uuid.UUID,
) ([]ArticleTagConnection, error) {
	rows, err := db.Stmts.QueryArticleTagConnectionsByArticleID(ctx, dbtx, articleID)
	if err != nil {
		return nil, err
	}

	connections := make([]ArticleTagConnection, len(rows))
	for i, row := range rows {
		connections[i] = ArticleTagConnection{
			ID:        row.ID,
			ArticleID: row.ArticleID,
			TagID:     row.TagID,
		}
	}

	return connections, nil
}

func GetArticleTagConnectionsByTagID(
	ctx context.Context,
	dbtx db.DBTX,
	tagID uuid.UUID,
) ([]ArticleTagConnection, error) {
	rows, err := db.Stmts.QueryArticleTagConnectionsByTagID(ctx, dbtx, tagID)
	if err != nil {
		return nil, err
	}

	connections := make([]ArticleTagConnection, len(rows))
	for i, row := range rows {
		connections[i] = ArticleTagConnection{
			ID:        row.ID,
			ArticleID: row.ArticleID,
			TagID:     row.TagID,
		}
	}

	return connections, nil
}

func GetArticleTagConnection(
	ctx context.Context,
	dbtx db.DBTX,
	articleID, tagID uuid.UUID,
) (ArticleTagConnection, error) {
	row, err := db.Stmts.QueryArticleTagConnection(ctx, dbtx, db.QueryArticleTagConnectionParams{
		ArticleID: articleID,
		TagID:     tagID,
	})
	if err != nil {
		return ArticleTagConnection{}, err
	}

	return ArticleTagConnection{
		ID:        row.ID,
		ArticleID: row.ArticleID,
		TagID:     row.TagID,
	}, nil
}

func NewArticleTagConnection(
	ctx context.Context,
	dbtx db.DBTX,
	articleID, tagID uuid.UUID,
) (ArticleTagConnection, error) {
	connection := ArticleTagConnection{
		ID:        uuid.New(),
		ArticleID: articleID,
		TagID:     tagID,
	}

	_, err := db.Stmts.InsertArticleTagConnection(ctx, dbtx, db.InsertArticleTagConnectionParams{
		ID:        connection.ID,
		ArticleID: connection.ArticleID,
		TagID:     connection.TagID,
	})
	if err != nil {
		return ArticleTagConnection{}, err
	}

	return connection, nil
}

func DeleteArticleTagConnection(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.Stmts.DeleteArticleTagConnection(ctx, dbtx, id)
}

func DeleteArticleTagConnectionByArticleAndTag(
	ctx context.Context,
	dbtx db.DBTX,
	articleID, tagID uuid.UUID,
) error {
	return db.Stmts.DeleteArticleTagConnectionByArticleAndTag(ctx, dbtx, db.DeleteArticleTagConnectionByArticleAndTagParams{
		ArticleID: articleID,
		TagID:     tagID,
	})
}

func DeleteArticleTagConnectionsByArticleID(
	ctx context.Context,
	dbtx db.DBTX,
	articleID uuid.UUID,
) error {
	return db.Stmts.DeleteArticleTagConnectionsByArticleID(ctx, dbtx, articleID)
}

func DeleteArticleTagConnectionsByTagID(
	ctx context.Context,
	dbtx db.DBTX,
	tagID uuid.UUID,
) error {
	return db.Stmts.DeleteArticleTagConnectionsByTagID(ctx, dbtx, tagID)
}
