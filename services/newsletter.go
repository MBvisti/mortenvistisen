package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	objectstorage "mortenvistisen/clients/storage"
	"mortenvistisen/config"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/queue/jobs"
)

type Newsletter struct {
	db           storage.Pool
	insertOnly   storage.InsertQueue
	coverStore   objectstorage.CoverStore
	maxCoverSize int64
	now          func() time.Time
}

func NewNewsletter(
	db storage.Pool,
	insertOnly queue.InsertOnly,
	coverStore objectstorage.CoverStore,
	cfg config.Config,
) Newsletter {
	return Newsletter{
		db:           db,
		insertOnly:   &insertOnly,
		coverStore:   coverStore,
		maxCoverSize: cfg.R2.MaxUploadBytes,
		now:          time.Now,
	}
}

func (n Newsletter) Create(
	ctx context.Context,
	data models.CreateNewsletterData,
	cover, metaImage *Cover,
) (models.NewsletterEntity, error) {
	data.ReleasedAt, _ = newsletterReleaseState(
		data.IsPublished.Bool,
		sql.NullTime{},
		n.now(),
	)

	_, uploaded, err := prepareAndUploadImages(
		ctx,
		n.coverStore,
		n.maxCoverSize,
		cover,
		metaImage,
		imageLinks{},
		func(links imageLinks) error {
			data.ImageLink = links.cover
			data.MetaImageLink = links.meta
			return data.Validate()
		},
	)
	if err != nil {
		return models.NewsletterEntity{}, err
	}
	cleanupUploads := true
	defer func() {
		if cleanupUploads {
			cleanupUploadedImages(ctx, n.coverStore, uploaded)
		}
	}()

	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return models.NewsletterEntity{}, fmt.Errorf(
			"begin newsletter creation transaction: %w",
			err,
		)
	}

	newsletter, err := models.Newsletter.Create(ctx, tx, data)
	if err != nil {
		_ = tx.Rollback()
		return models.NewsletterEntity{}, err
	}

	if err := tx.Commit(); err != nil {
		if n.shouldRetainAnyImage(ctx, newsletter.ID, uploaded) {
			cleanupUploads = false
		}
		return models.NewsletterEntity{}, fmt.Errorf(
			"commit newsletter creation transaction: %w",
			err,
		)
	}
	cleanupUploads = false

	return newsletter, nil
}

func (n Newsletter) Update(
	ctx context.Context,
	data models.UpdateNewsletterData,
	cover, metaImage *Cover,
	removeCover, removeMetaImage bool,
) (models.NewsletterEntity, error) {
	if removeCover && cover != nil {
		return models.NewsletterEntity{}, coverValidationError(
			"conflict",
			"Choose a new cover or remove the existing cover, but not both.",
		)
	}
	if removeMetaImage && metaImage != nil {
		return models.NewsletterEntity{}, imageValidationError(
			"metaImage",
			"conflict",
			"Choose a new social image or remove the existing social image, but not both.",
		)
	}
	currentBeforeUpload, err := models.Newsletter.Find(ctx, n.db.Executor(), data.ID)
	if err != nil {
		return models.NewsletterEntity{}, fmt.Errorf("find newsletter before cover upload: %w", err)
	}
	data.ReleasedAt, _ = newsletterReleaseState(
		data.IsPublished.Bool,
		currentBeforeUpload.ReleasedAt,
		n.now(),
	)
	links := imageLinks{
		cover: currentBeforeUpload.ImageLink,
		meta:  currentBeforeUpload.MetaImageLink,
	}
	if removeCover {
		links.cover = sql.NullString{}
	}
	if removeMetaImage {
		links.meta = sql.NullString{}
	}

	_, uploaded, err := prepareAndUploadImages(
		ctx,
		n.coverStore,
		n.maxCoverSize,
		cover,
		metaImage,
		links,
		func(links imageLinks) error {
			data.ImageLink = links.cover
			data.MetaImageLink = links.meta
			return data.Validate()
		},
	)
	if err != nil {
		return models.NewsletterEntity{}, err
	}
	cleanupUploads := true
	defer func() {
		if cleanupUploads {
			cleanupUploadedImages(ctx, n.coverStore, uploaded)
		}
	}()

	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return models.NewsletterEntity{}, fmt.Errorf("begin newsletter update transaction: %w", err)
	}

	current, err := models.Newsletter.FindForUpdate(ctx, tx, data.ID)
	if err != nil {
		_ = tx.Rollback()
		return models.NewsletterEntity{}, fmt.Errorf("find newsletter for update: %w", err)
	}

	data.ReleasedAt, _ = newsletterReleaseState(
		data.IsPublished.Bool,
		current.ReleasedAt,
		n.now(),
	)
	if removeCover {
		data.ImageLink = sql.NullString{}
	} else if cover == nil {
		data.ImageLink = current.ImageLink
	}
	if removeMetaImage {
		data.MetaImageLink = sql.NullString{}
	} else if metaImage == nil {
		data.MetaImageLink = current.MetaImageLink
	}
	if err := data.Validate(); err != nil {
		_ = tx.Rollback()
		return models.NewsletterEntity{}, mapImageValidation(err)
	}

	newsletter, err := models.Newsletter.Update(ctx, tx, data)
	if err != nil {
		_ = tx.Rollback()
		return models.NewsletterEntity{}, err
	}
	if !imageLinksEqual(current.ImageLink, newsletter.ImageLink) {
		if err := n.enqueueImageDeletion(ctx, tx.Tx, current.ID, current.ImageLink); err != nil {
			_ = tx.Rollback()
			return models.NewsletterEntity{}, err
		}
	}
	if !imageLinksEqual(current.MetaImageLink, newsletter.MetaImageLink) {
		if err := n.enqueueImageDeletion(ctx, tx.Tx, current.ID, current.MetaImageLink); err != nil {
			_ = tx.Rollback()
			return models.NewsletterEntity{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		if n.shouldRetainAnyImage(ctx, newsletter.ID, uploaded) {
			cleanupUploads = false
		}
		return models.NewsletterEntity{}, fmt.Errorf(
			"commit newsletter update transaction: %w",
			err,
		)
	}
	cleanupUploads = false

	return newsletter, nil
}

func (n Newsletter) enqueueImageDeletion(
	ctx context.Context,
	sqlTx *sql.Tx,
	newsletterID int32,
	link sql.NullString,
) error {
	if !link.Valid {
		return nil
	}
	key, managed := n.coverStore.CoverKey(link.String)
	if !managed {
		return nil
	}
	if _, err := n.insertOnly.InsertTx(ctx, sqlTx, jobs.DeleteNewsletterCoverArgs{
		NewsletterID: newsletterID,
		Key:          key,
	}, nil); err != nil {
		return fmt.Errorf("enqueue newsletter image deletion: %w", err)
	}
	return nil
}

func (n Newsletter) shouldRetainAnyImage(
	ctx context.Context,
	id int32,
	uploaded []objectstorage.CoverObject,
) bool {
	if len(uploaded) == 0 {
		return false
	}
	lookupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()
	newsletter, err := models.Newsletter.Find(lookupCtx, n.db.Executor(), id)
	if err == nil {
		for _, object := range uploaded {
			if (newsletter.ImageLink.Valid && newsletter.ImageLink.String == object.URL) ||
				(newsletter.MetaImageLink.Valid && newsletter.MetaImageLink.String == object.URL) {
				return true
			}
		}
		return false
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false
	}
	slog.ErrorContext(
		lookupCtx,
		"could not determine whether newsletter images were committed",
		"newsletter_id", id,
		"error", err,
	)
	return true
}

func newsletterReleaseState(
	published bool,
	current sql.NullTime,
	now time.Time,
) (sql.NullTime, bool) {
	if current.Valid || !published {
		return current, false
	}

	return sql.NullTime{Time: now, Valid: true}, true
}
