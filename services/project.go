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

type Project struct {
	db           storage.Pool
	insertOnly   storage.InsertQueue
	coverStore   objectstorage.CoverStore
	maxCoverSize int64
}

func NewProject(
	db storage.Pool,
	insertOnly queue.InsertOnly,
	coverStore objectstorage.CoverStore,
	cfg config.Config,
) Project {
	return Project{
		db:           db,
		insertOnly:   &insertOnly,
		coverStore:   coverStore,
		maxCoverSize: cfg.R2.MaxUploadBytes,
	}
}

func (p Project) Create(
	ctx context.Context,
	data models.CreateProjectData,
	cover, metaImage *Cover,
) (models.ProjectEntity, error) {
	_, uploaded, err := prepareAndUploadImages(
		ctx,
		p.coverStore,
		p.maxCoverSize,
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
		return models.ProjectEntity{}, err
	}
	cleanupUploads := true
	defer func() {
		if cleanupUploads {
			cleanupUploadedImages(ctx, p.coverStore, uploaded)
		}
	}()

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ProjectEntity{}, fmt.Errorf("begin project creation transaction: %w", err)
	}
	project, err := models.Project.Create(ctx, tx, data)
	if err != nil {
		_ = tx.Rollback()
		return models.ProjectEntity{}, err
	}
	if err := tx.Commit(); err != nil {
		if p.shouldRetainAnyImage(ctx, project.ID, uploaded) {
			cleanupUploads = false
		}
		return models.ProjectEntity{}, fmt.Errorf("commit project creation transaction: %w", err)
	}
	cleanupUploads = false
	return project, nil
}

func (p Project) Update(
	ctx context.Context,
	data models.UpdateProjectData,
	cover, metaImage *Cover,
	removeCover, removeMetaImage bool,
) (models.ProjectEntity, error) {
	if removeCover && cover != nil {
		return models.ProjectEntity{}, coverValidationError(
			"conflict",
			"Choose a new cover or remove the existing cover, but not both.",
		)
	}
	if removeMetaImage && metaImage != nil {
		return models.ProjectEntity{}, imageValidationError(
			"metaImage",
			"conflict",
			"Choose a new social image or remove the existing social image, but not both.",
		)
	}

	currentBeforeUpload, err := models.Project.Find(ctx, p.db.Executor(), data.ID)
	if err != nil {
		return models.ProjectEntity{}, fmt.Errorf("find project before image upload: %w", err)
	}
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
		p.coverStore,
		p.maxCoverSize,
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
		return models.ProjectEntity{}, err
	}
	cleanupUploads := true
	defer func() {
		if cleanupUploads {
			cleanupUploadedImages(ctx, p.coverStore, uploaded)
		}
	}()

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ProjectEntity{}, fmt.Errorf("begin project update transaction: %w", err)
	}
	current, err := models.Project.FindForUpdate(ctx, tx, data.ID)
	if err != nil {
		_ = tx.Rollback()
		return models.ProjectEntity{}, fmt.Errorf("find project for update: %w", err)
	}
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
		return models.ProjectEntity{}, mapImageValidation(err)
	}

	project, err := models.Project.Update(ctx, tx, data)
	if err != nil {
		_ = tx.Rollback()
		return models.ProjectEntity{}, err
	}
	if !imageLinksEqual(current.ImageLink, project.ImageLink) {
		if err := p.enqueueImageDeletion(ctx, tx.Tx, current.ID, current.ImageLink); err != nil {
			_ = tx.Rollback()
			return models.ProjectEntity{}, err
		}
	}
	if !imageLinksEqual(current.MetaImageLink, project.MetaImageLink) {
		if err := p.enqueueImageDeletion(ctx, tx.Tx, current.ID, current.MetaImageLink); err != nil {
			_ = tx.Rollback()
			return models.ProjectEntity{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		if p.shouldRetainAnyImage(ctx, project.ID, uploaded) {
			cleanupUploads = false
		}
		return models.ProjectEntity{}, fmt.Errorf("commit project update transaction: %w", err)
	}
	cleanupUploads = false
	return project, nil
}

func (p Project) enqueueImageDeletion(
	ctx context.Context,
	sqlTx *sql.Tx,
	projectID int32,
	link sql.NullString,
) error {
	if !link.Valid {
		return nil
	}
	key, managed := p.coverStore.CoverKey(link.String)
	if !managed {
		return nil
	}
	if _, err := p.insertOnly.InsertTx(ctx, sqlTx, jobs.DeleteProjectCoverArgs{
		ProjectID: projectID,
		Key:       key,
	}, nil); err != nil {
		return fmt.Errorf("enqueue project image deletion: %w", err)
	}
	return nil
}

func (p Project) shouldRetainAnyImage(
	ctx context.Context,
	id int32,
	uploaded []objectstorage.CoverObject,
) bool {
	if len(uploaded) == 0 {
		return false
	}
	lookupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()
	project, err := models.Project.Find(lookupCtx, p.db.Executor(), id)
	if err == nil {
		for _, object := range uploaded {
			if (project.ImageLink.Valid && project.ImageLink.String == object.URL) ||
				(project.MetaImageLink.Valid && project.MetaImageLink.String == object.URL) {
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
		"could not determine whether project images were committed",
		"project_id", id,
		"error", err,
	)
	return true
}
