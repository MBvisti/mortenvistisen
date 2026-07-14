package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	objectstorage "mortenvistisen/clients/storage"
	"mortenvistisen/config"
	"mortenvistisen/email"
	"mortenvistisen/internal/storage"
	"mortenvistisen/internal/validation"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/queue/jobs"

	"github.com/riverqueue/river"
)

type Article struct {
	db           storage.Pool
	insertOnly   storage.InsertQueue
	coverStore   objectstorage.CoverStore
	maxCoverSize int64
	now          func() time.Time
}

type ArticleTagSelection struct {
	TagIDs       []int32
	NewTagTitles []string
}

type Cover struct {
	Filename string
	Size     int64
	Body     io.Reader
}

type imageLinks struct {
	cover sql.NullString
	meta  sql.NullString
}

func NewArticle(
	db storage.Pool,
	insertOnly queue.InsertOnly,
	coverStore objectstorage.CoverStore,
	cfg config.Config,
) Article {
	return Article{
		db:           db,
		insertOnly:   &insertOnly,
		coverStore:   coverStore,
		maxCoverSize: cfg.R2.MaxUploadBytes,
		now:          time.Now,
	}
}

func (a Article) Create(
	ctx context.Context,
	data models.CreateArticleData,
	tagSelection ArticleTagSelection,
	cover, metaImage *Cover,
) (models.ArticleEntity, error) {
	var firstPublication bool
	data.FirstPublishedAt, firstPublication = firstPublicationState(
		data.Published,
		sql.NullTime{},
		a.now(),
	)

	_, uploaded, err := prepareAndUploadImages(
		ctx,
		a.coverStore,
		a.maxCoverSize,
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
		return models.ArticleEntity{}, err
	}
	cleanupUploads := true
	defer func() {
		if cleanupUploads {
			cleanupUploadedImages(ctx, a.coverStore, uploaded)
		}
	}()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ArticleEntity{}, fmt.Errorf("begin article creation transaction: %w", err)
	}

	article, err := models.Article.Create(ctx, tx, data)
	if err != nil {
		_ = tx.Rollback()
		return models.ArticleEntity{}, err
	}
	tagIDs, err := resolveArticleTagIDs(ctx, tx, tagSelection)
	if err != nil {
		_ = tx.Rollback()
		return models.ArticleEntity{}, err
	}
	if err := models.ArticleTagConnection.ReplaceForArticle(
		ctx,
		tx,
		article.ID,
		tagIDs,
	); err != nil {
		_ = tx.Rollback()
		return models.ArticleEntity{}, fmt.Errorf("attach tags to article: %w", err)
	}
	if firstPublication {
		if err := a.enqueueFirstPublication(ctx, tx.Tx, tx, article); err != nil {
			_ = tx.Rollback()
			return models.ArticleEntity{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		if a.shouldRetainAnyImage(ctx, article.ID, uploaded) {
			cleanupUploads = false
		}
		return models.ArticleEntity{}, fmt.Errorf("commit article creation transaction: %w", err)
	}
	cleanupUploads = false

	return article, nil
}

func (a Article) Update(
	ctx context.Context,
	data models.UpdateArticleData,
	tagSelection ArticleTagSelection,
	cover, metaImage *Cover,
	removeCover, removeMetaImage bool,
) (models.ArticleEntity, error) {
	if removeCover && cover != nil {
		return models.ArticleEntity{}, coverValidationError(
			"conflict",
			"Choose a new cover or remove the existing cover, but not both.",
		)
	}
	if removeMetaImage && metaImage != nil {
		return models.ArticleEntity{}, imageValidationError(
			"metaImage",
			"conflict",
			"Choose a new social image or remove the existing social image, but not both.",
		)
	}

	currentBeforeUpload, err := models.Article.Find(ctx, a.db.Executor(), data.ID)
	if err != nil {
		return models.ArticleEntity{}, fmt.Errorf("find article before cover upload: %w", err)
	}
	data.FirstPublishedAt = currentBeforeUpload.FirstPublishedAt
	data.FirstPublishedAt, _ = firstPublicationState(
		data.Published,
		currentBeforeUpload.FirstPublishedAt,
		a.now(),
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
		a.coverStore,
		a.maxCoverSize,
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
		return models.ArticleEntity{}, err
	}
	cleanupUploads := true
	defer func() {
		if cleanupUploads {
			cleanupUploadedImages(ctx, a.coverStore, uploaded)
		}
	}()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ArticleEntity{}, fmt.Errorf("begin article update transaction: %w", err)
	}

	current, err := models.Article.FindForUpdate(ctx, tx, data.ID)
	if err != nil {
		_ = tx.Rollback()
		return models.ArticleEntity{}, fmt.Errorf("find article for update: %w", err)
	}

	data.FirstPublishedAt = current.FirstPublishedAt
	var firstPublication bool
	data.FirstPublishedAt, firstPublication = firstPublicationState(
		data.Published,
		current.FirstPublishedAt,
		a.now(),
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
		return models.ArticleEntity{}, mapImageValidation(err)
	}

	article, err := models.Article.Update(ctx, tx, data)
	if err != nil {
		_ = tx.Rollback()
		return models.ArticleEntity{}, err
	}
	if !imageLinksEqual(current.ImageLink, article.ImageLink) {
		if err := a.enqueueImageDeletion(ctx, tx.Tx, current.ID, current.ImageLink); err != nil {
			_ = tx.Rollback()
			return models.ArticleEntity{}, err
		}
	}
	if !imageLinksEqual(current.MetaImageLink, article.MetaImageLink) {
		if err := a.enqueueImageDeletion(ctx, tx.Tx, current.ID, current.MetaImageLink); err != nil {
			_ = tx.Rollback()
			return models.ArticleEntity{}, err
		}
	}
	tagIDs, err := resolveArticleTagIDs(ctx, tx, tagSelection)
	if err != nil {
		_ = tx.Rollback()
		return models.ArticleEntity{}, err
	}
	if err := models.ArticleTagConnection.ReplaceForArticle(
		ctx,
		tx,
		article.ID,
		tagIDs,
	); err != nil {
		_ = tx.Rollback()
		return models.ArticleEntity{}, fmt.Errorf("attach tags to article: %w", err)
	}
	if firstPublication {
		if err := a.enqueueFirstPublication(ctx, tx.Tx, tx, article); err != nil {
			_ = tx.Rollback()
			return models.ArticleEntity{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		if a.shouldRetainAnyImage(ctx, article.ID, uploaded) {
			cleanupUploads = false
		}
		return models.ArticleEntity{}, fmt.Errorf("commit article update transaction: %w", err)
	}
	cleanupUploads = false

	return article, nil
}

type ArticlePatch struct {
	Title           *string
	Excerpt         *string
	MetaTitle       *string
	MetaDescription *string
}

func (p ArticlePatch) Empty() bool {
	return p.Title == nil && p.Excerpt == nil && p.MetaTitle == nil && p.MetaDescription == nil
}

func (a Article) Patch(
	ctx context.Context,
	slug string,
	patch ArticlePatch,
) (models.ArticleEntity, error) {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ArticleEntity{}, fmt.Errorf("begin article patch transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	current, err := models.Article.FindForUpdateBySlug(ctx, tx, slug)
	if err != nil {
		return models.ArticleEntity{}, err
	}

	article, err := models.Article.Update(ctx, tx, applyArticlePatch(current, patch))
	if err != nil {
		return models.ArticleEntity{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.ArticleEntity{}, fmt.Errorf("commit article patch transaction: %w", err)
	}

	return article, nil
}

func applyArticlePatch(current models.ArticleEntity, patch ArticlePatch) models.UpdateArticleData {
	data := models.UpdateArticleData{
		ID:               current.ID,
		FirstPublishedAt: current.FirstPublishedAt,
		Published:        current.Published,
		Title:            current.Title,
		Excerpt:          current.Excerpt,
		MetaTitle:        current.MetaTitle,
		MetaDescription:  current.MetaDescription,
		Slug:             current.Slug,
		ImageLink:        current.ImageLink,
		MetaImageLink:    current.MetaImageLink,
		ReadTime:         current.ReadTime,
		Content:          current.Content,
	}
	if patch.Title != nil {
		data.Title = *patch.Title
	}
	if patch.Excerpt != nil {
		data.Excerpt = sql.NullString{String: *patch.Excerpt, Valid: true}
	}
	if patch.MetaTitle != nil {
		data.MetaTitle = sql.NullString{String: *patch.MetaTitle, Valid: true}
	}
	if patch.MetaDescription != nil {
		data.MetaDescription = sql.NullString{String: *patch.MetaDescription, Valid: true}
	}
	return data
}

func (a Article) enqueueImageDeletion(
	ctx context.Context,
	sqlTx *sql.Tx,
	articleID int32,
	link sql.NullString,
) error {
	if !link.Valid {
		return nil
	}
	key, managed := a.coverStore.CoverKey(link.String)
	if !managed {
		return nil
	}

	if _, err := a.insertOnly.InsertTx(
		ctx,
		sqlTx,
		jobs.DeleteArticleCoverArgs{ArticleID: articleID, Key: key},
		nil,
	); err != nil {
		return fmt.Errorf("enqueue article image deletion: %w", err)
	}

	return nil
}

func prepareCover(
	coverStore objectstorage.CoverStore,
	maxCoverSize int64,
	cover *Cover,
) (*objectstorage.CoverUpload, string, error) {
	if cover == nil {
		return nil, "", nil
	}

	filename := cover.Filename
	if strings.TrimSpace(filename) == "" || strings.TrimSpace(filename) != filename ||
		filename == "." || filename == ".." ||
		strings.ContainsAny(filename, `/\\`) ||
		strings.IndexFunc(filename, unicode.IsControl) >= 0 || len(filename) > 255 {
		return nil, "", coverValidationError("filename", "Choose an image with a valid filename.")
	}
	if cover.Size <= 0 {
		return nil, "", coverValidationError("empty", "The selected image is empty.")
	}
	if cover.Size > maxCoverSize {
		return nil, "", coverValidationError("max_size", "The selected image must be 10 MB or smaller.")
	}

	extension := strings.ToLower(filepath.Ext(filename))
	expectedContentType, allowed := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
		".avif": "image/avif",
		".gif":  "image/gif",
	}[extension]
	if !allowed {
		return nil, "", coverValidationError(
			"type",
			"Use a JPEG, PNG, WebP, AVIF, or GIF image.",
		)
	}

	header := make([]byte, 512)
	read, err := io.ReadFull(cover.Body, header)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, "", fmt.Errorf("read cover header: %w", err)
	}
	header = header[:read]
	contentType := detectImageContentType(header)
	if contentType != expectedContentType {
		return nil, "", coverValidationError(
			"type",
			"The file contents do not match its image extension.",
		)
	}

	publicURL, err := coverStore.CoverURL(filename)
	if err != nil {
		return nil, "", fmt.Errorf("build cover URL: %w", err)
	}

	return &objectstorage.CoverUpload{
		Filename:    filename,
		ContentType: contentType,
		Size:        cover.Size,
		Body:        io.MultiReader(bytes.NewReader(header), cover.Body),
	}, publicURL, nil
}

func prepareAndUploadImages(
	ctx context.Context,
	store objectstorage.CoverStore,
	maxSize int64,
	cover, metaImage *Cover,
	links imageLinks,
	validate func(imageLinks) error,
) (imageLinks, []objectstorage.CoverObject, error) {
	coverUpload, coverURL, err := prepareCover(store, maxSize, cover)
	if err != nil {
		return imageLinks{}, nil, err
	}
	metaUpload, metaURL, err := prepareCover(store, maxSize, metaImage)
	if err != nil {
		return imageLinks{}, nil, remapValidationField(err, "cover", "metaImage")
	}
	if coverUpload != nil {
		links.cover = sql.NullString{String: coverURL, Valid: true}
	}
	if metaUpload != nil {
		links.meta = sql.NullString{String: metaURL, Valid: true}
	}
	if err := validate(links); err != nil {
		return imageLinks{}, nil, mapImageValidation(err)
	}

	uploads := []struct {
		prepared *objectstorage.CoverUpload
		field    string
	}{
		{prepared: coverUpload, field: "cover"},
		{prepared: metaUpload, field: "metaImage"},
	}
	uploaded := make([]objectstorage.CoverObject, 0, 2)
	for _, upload := range uploads {
		if upload.prepared == nil {
			continue
		}
		object, uploadErr := store.UploadCover(ctx, *upload.prepared)
		if uploadErr != nil {
			cleanupUploadedImages(ctx, store, uploaded)
			if errors.Is(uploadErr, objectstorage.ErrObjectExists) {
				return imageLinks{}, nil, imageValidationError(
					upload.field,
					"exists",
					"An image with this filename already exists.",
				)
			}
			return imageLinks{}, nil, uploadErr
		}
		uploaded = append(uploaded, object)
	}

	return links, uploaded, nil
}

func cleanupUploadedImages(
	ctx context.Context,
	store objectstorage.CoverStore,
	uploaded []objectstorage.CoverObject,
) {
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()
	for _, object := range uploaded {
		if err := store.Delete(cleanupCtx, object.Key); err != nil {
			slog.ErrorContext(cleanupCtx, "could not clean up uploaded image", "key", object.Key, "error", err)
		}
	}
}

func remapValidationField(err error, from, to string) error {
	validationErrors, ok := validation.As(err)
	if !ok {
		return err
	}
	mapped := make(validation.ValidationErrors, len(validationErrors))
	copy(mapped, validationErrors)
	for index := range mapped {
		if mapped[index].Field == from {
			mapped[index].Field = to
		}
	}
	return errors.Join(models.ErrDomainValidation, mapped)
}

func imageLinksEqual(left, right sql.NullString) bool {
	return left.Valid == right.Valid && (!left.Valid || left.String == right.String)
}

func detectImageContentType(header []byte) string {
	if len(header) >= 12 && string(header[4:8]) == "ftyp" {
		boxSize := int(binary.BigEndian.Uint32(header[:4]))
		if boxSize <= 0 || boxSize > len(header) {
			boxSize = len(header)
		}
		for offset := 8; offset+4 <= boxSize; offset += 4 {
			brand := string(header[offset : offset+4])
			if brand == "avif" || brand == "avis" {
				return "image/avif"
			}
		}
	}

	return http.DetectContentType(header)
}

func (a Article) shouldRetainAnyImage(
	ctx context.Context,
	articleID int32,
	uploaded []objectstorage.CoverObject,
) bool {
	if len(uploaded) == 0 {
		return false
	}
	lookupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()

	article, err := models.Article.Find(lookupCtx, a.db.Executor(), articleID)
	if err == nil {
		for _, object := range uploaded {
			if (article.ImageLink.Valid && article.ImageLink.String == object.URL) ||
				(article.MetaImageLink.Valid && article.MetaImageLink.String == object.URL) {
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
		"could not determine whether article images were committed",
		"article_id", articleID,
		"error", err,
	)
	return true
}

func imageValidationError(field, code, message string) error {
	return errors.Join(models.ErrDomainValidation, validation.ValidationErrors{{
		Field:   field,
		Code:    code,
		Message: message,
	}})
}

func coverValidationError(code, message string) error {
	return imageValidationError("cover", code, message)
}

func mapImageValidation(err error) error {
	validationErrors, ok := validation.As(err)
	if !ok {
		return err
	}

	mapped := make(validation.ValidationErrors, len(validationErrors))
	copy(mapped, validationErrors)
	for index := range mapped {
		switch mapped[index].Field {
		case "imageLink":
			mapped[index].Field = "cover"
		case "metaImageLink":
			mapped[index].Field = "metaImage"
		}
	}

	return errors.Join(models.ErrDomainValidation, mapped)
}

func resolveArticleTagIDs(
	ctx context.Context,
	db storage.Executor,
	selection ArticleTagSelection,
) ([]int32, error) {
	requestedIDs := uniqueArticleTagIDs(selection.TagIDs)
	if len(requestedIDs) > 0 {
		existing, err := models.Tag.FindByIDs(ctx, db, requestedIDs)
		if err != nil {
			return nil, fmt.Errorf("find selected tags: %w", err)
		}
		if len(existing) != len(requestedIDs) {
			return nil, errors.Join(models.ErrDomainValidation, validation.ValidationErrors{{
				Field:   "tagIds",
				Code:    "invalid",
				Message: "One or more selected tags no longer exist.",
			}})
		}
	}

	resolvedIDs := append([]int32(nil), requestedIDs...)
	seenIDs := make(map[int32]struct{}, len(resolvedIDs))
	for _, id := range resolvedIDs {
		seenIDs[id] = struct{}{}
	}
	seenTitles := make(map[string]struct{}, len(selection.NewTagTitles))
	for _, requestedTitle := range selection.NewTagTitles {
		title := strings.TrimSpace(requestedTitle)
		normalizedTitle := strings.ToLower(title)
		if _, exists := seenTitles[normalizedTitle]; exists {
			continue
		}
		seenTitles[normalizedTitle] = struct{}{}

		tag, err := models.Tag.FindByTitle(ctx, db, title)
		if errors.Is(err, sql.ErrNoRows) {
			tag, err = models.Tag.Create(ctx, db, models.CreateTagData{Title: title})
		}
		if err != nil {
			if tagValidationErrors, ok := validation.As(err); ok {
				mappedErrors := make(validation.ValidationErrors, 0, len(tagValidationErrors))
				for _, tagValidationErr := range tagValidationErrors {
					tagValidationErr.Field = "newTags"
					mappedErrors = append(mappedErrors, tagValidationErr)
				}
				return nil, errors.Join(models.ErrDomainValidation, mappedErrors)
			}
			return nil, fmt.Errorf("resolve tag %q: %w", title, err)
		}
		if _, exists := seenIDs[tag.ID]; exists {
			continue
		}
		seenIDs[tag.ID] = struct{}{}
		resolvedIDs = append(resolvedIDs, tag.ID)
	}

	return resolvedIDs, nil
}

func uniqueArticleTagIDs(ids []int32) []int32 {
	unique := make([]int32, 0, len(ids))
	seen := make(map[int32]struct{}, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}

	return unique
}

func firstPublicationState(
	published bool,
	current sql.NullTime,
	now time.Time,
) (sql.NullTime, bool) {
	if current.Valid || !published {
		return current, false
	}

	return sql.NullTime{Time: now, Valid: true}, true
}

const articleUnsubscribeTODOPath = "/TODO-unsubscribe"

func (a Article) enqueueFirstPublication(
	ctx context.Context,
	sqlTx *sql.Tx,
	db storage.Executor,
	article models.ArticleEntity,
) error {
	subscribers, err := models.Subscriber.Verified(ctx, db)
	if err != nil {
		return fmt.Errorf("find article notification subscribers: %w", err)
	}
	if len(subscribers) == 0 {
		return nil
	}

	params, err := firstPublicationJobs(article, subscribers)
	if err != nil {
		return err
	}

	if _, err := a.insertOnly.InsertManyTx(ctx, sqlTx, params); err != nil {
		return fmt.Errorf("queue article publication emails: %w", err)
	}

	return nil
}

func firstPublicationJobs(
	article models.ArticleEntity,
	subscribers []models.SubscriberEntity,
) ([]river.InsertManyParams, error) {
	// TODO: Replace this with a subscriber-specific unsubscribe URL.
	unsubscribeURL := config.BaseURL + articleUnsubscribeTODOPath
	message := email.ArticlePublished{
		Title:          article.Title,
		Excerpt:        article.Excerpt.String,
		ImageURL:       article.ImageLink.String,
		UnsubscribeURL: unsubscribeURL,
	}
	html, err := message.ToHTML()
	if err != nil {
		return nil, fmt.Errorf("render article publication email HTML: %w", err)
	}
	text, err := message.ToText()
	if err != nil {
		return nil, fmt.Errorf("render article publication email text: %w", err)
	}

	params := make([]river.InsertManyParams, 0, len(subscribers))
	for _, subscriber := range subscribers {
		params = append(params, river.InsertManyParams{
			Args: jobs.SendMarketingEmailArgs{
				Data: email.MarketingData{
					To:             []string{subscriber.Email.String},
					From:           config.DefaultSenderSignature,
					Subject:        "New article: " + article.Title,
					HTMLBody:       html,
					TextBody:       text,
					UnsubscribeURL: unsubscribeURL,
					Tags:           []string{"article-published"},
					Metadata: map[string]string{
						"article_id":    strconv.FormatInt(int64(article.ID), 10),
						"subscriber_id": strconv.FormatInt(int64(subscriber.ID), 10),
					},
				},
			},
		})
	}

	return params, nil
}
