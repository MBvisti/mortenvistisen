package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type postDatabase interface {
	InsertPost(ctx context.Context, arg database.InsertPostParams) (uuid.UUID, error)
	UpdatePost(ctx context.Context, arg database.UpdatePostParams) (database.Post, error)
	AssociateTagWithPost(ctx context.Context, arg database.AssociateTagWithPostParams) error
}

func NewPost(
	ctx context.Context,
	db postDatabase,
	v *validator.Validate,
	newPost domain.NewPost,
	associatedTags []string,
) error {
	if err := v.Struct(newPost); err != nil {
		telemetry.Logger.Error("provided post data did not pass the validation", "error", err)
		return err
	}

	now := time.Now()

	args := database.InsertPostParams{
		ID:          uuid.New(),
		CreatedAt:   database.ConvertToPGTimestamp(now),
		UpdatedAt:   database.ConvertToPGTimestamp(now),
		Title:       newPost.Title,
		HeaderTitle: sql.NullString{Valid: true, String: newPost.HeaderTitle},
		Filename:    newPost.Filename,
		Slug:        slug.MakeLang(newPost.Title, "en"),
		Excerpt:     newPost.Excerpt,
		Draft:       true,
	}

	if newPost.ReleaseNow {
		args.ReleasedAt = database.ConvertToPGTimestamp(now)
		args.Draft = false
	}

	id, err := db.InsertPost(ctx, args)
	if err != nil {
		return err
	}

	// TODO: run in transaction
	for _, associatedTag := range associatedTags {
		tagID, err := uuid.Parse(associatedTag)
		if err != nil {
			return err
		}

		if err := db.AssociateTagWithPost(
			ctx,
			database.AssociateTagWithPostParams{
				ID:     uuid.New(),
				PostID: id,
				TagID:  tagID,
			}); err != nil {
			return err
		}
	}

	return nil
}

func UpdatePost(
	ctx context.Context,
	db postDatabase,
	v *validator.Validate,
	updatePost domain.UpdatePost,
) (domain.Post, error) {
	if err := v.Struct(updatePost); err != nil {
		telemetry.Logger.Error("provided post data did not pass the validation", "error", err)
		return domain.Post{}, err
	}

	now := time.Now()
	args := database.UpdatePostParams{
		ID:          updatePost.ID,
		UpdatedAt:   database.ConvertToPGTimestamp(now),
		Title:       updatePost.Title,
		HeaderTitle: sql.NullString{Valid: true, String: updatePost.HeaderTitle},
		Slug:        slug.MakeLang(updatePost.Title, "en"),
		Excerpt:     updatePost.Excerpt,
		Draft:       updatePost.ReleaseNow,
		ReleasedAt:  database.ConvertToPGTimestamp(now),
		ReadTime:    sql.NullInt32{Int32: updatePost.EstimatedReadTime, Valid: true},
	}

	updatedPost, err := db.UpdatePost(ctx, args)
	if err != nil {
		return domain.Post{}, err
	}

	return domain.Post{
		ID:          updatedPost.ID,
		CreatedAt:   database.ConvertFromPGTimestampToTime(updatedPost.CreatedAt),
		UpdatedAt:   database.ConvertFromPGTimestampToTime(updatedPost.UpdatedAt),
		Title:       updatedPost.Title,
		Filename:    updatedPost.Filename,
		Slug:        updatedPost.Slug,
		Excerpt:     updatedPost.Excerpt,
		Draft:       updatedPost.Draft,
		ReleaseDate: database.ConvertFromPGTimestampToTime(updatedPost.ReleasedAt),
		ReadTime:    updatedPost.ReadTime.Int32,
	}, nil
}
