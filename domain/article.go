package domain

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	HeaderTitle string
	Filename    string
	Slug        string
	Excerpt     string
	Draft       bool
	ReleaseDate time.Time
	ReadTime    int32
	Tags        []Tag
}

func NewArticle(
	title string,
	headerTitle string,
	filename string,
	slug string,
	excerpt string,
	estimatedReadTime int32,
	tags []Tag,
) (Article, error) {
	now := time.Now()
	article := Article{
		ID:          uuid.New(),
		CreatedAt:   now,
		UpdatedAt:   now,
		Title:       title,
		HeaderTitle: headerTitle,
		Filename:    filename,
		Slug:        slug,
		Excerpt:     excerpt,
		Draft:       true,
		ReadTime:    estimatedReadTime,
		Tags:        tags,
	}

	if err := article.Validate(); err != nil {
		return Article{}, err
	}

	return article, nil
}

func (a Article) Validate() error {
	titleValidationErr := ErrValidation{
		FieldName:  "Title",
		FieldValue: a.Title,
	}
	if a.Title == "" {
		titleValidationErr.Violations = append(titleValidationErr.Violations, ErrIsRequired)
	}
	if len(a.Title) < 2 {
		titleValidationErr.Violations = append(titleValidationErr.Violations, ErrTooShort)
	}

	headerTitleValidationErr := ErrValidation{
		FieldName:  "HeaderTitle",
		FieldValue: a.HeaderTitle,
	}
	if a.HeaderTitle == "" {
		headerTitleValidationErr.Violations = append(
			headerTitleValidationErr.Violations,
			ErrIsRequired,
		)
	}
	if len(a.HeaderTitle) < 2 {
		headerTitleValidationErr.Violations = append(
			headerTitleValidationErr.Violations,
			ErrTooShort,
		)
	}

	excerptValidationErr := ErrValidation{
		FieldName:  "Excerpt",
		FieldValue: a.Excerpt,
	}
	if a.Excerpt == "" {
		excerptValidationErr.Violations = append(
			excerptValidationErr.Violations,
			ErrIsRequired,
		)
	}
	if len(a.Excerpt) < 130 {
		excerptValidationErr.Violations = append(
			excerptValidationErr.Violations,
			ErrTooShort,
		)
	}
	if len(a.Excerpt) > 160 {
		excerptValidationErr.Violations = append(
			excerptValidationErr.Violations,
			ErrTooLong,
		)
	}

	estimatedReadingTimeValidationErr := ErrValidation{
		FieldName:  "ReadTime",
		FieldValue: a.ReadTime,
	}
	if a.ReadTime == 0 {
		estimatedReadingTimeValidationErr.Violations = append(
			estimatedReadingTimeValidationErr.Violations,
			ErrIsRequired,
		)
	}

	filenameValidationErr := ErrValidation{
		FieldName:  "Filename",
		FieldValue: a.Filename,
	}
	if a.Filename == "" {
		filenameValidationErr.Violations = append(
			filenameValidationErr.Violations,
			ErrIsRequired,
		)
	}

	e := constructValidationErrors(
		titleValidationErr,
		headerTitleValidationErr,
		excerptValidationErr,
		estimatedReadingTimeValidationErr,
		filenameValidationErr,
	)
	if len(e) > 0 {
		return e
	}

	return nil
}

// type UpdateArticle struct {
// 	ID                uuid.UUID `validate:"required"`
// 	Title             string    `validate:"required,gte=3"`
// 	HeaderTitle       string    `validate:"required"`
// 	Excerpt           string    `validate:"required,lte=160,gte=130"`
// 	Slug              string    `validate:"required"`
// 	ReleaedAt         time.Time `validate:"required"`
// 	ReleaseNow        bool
// 	EstimatedReadTime int32 `validate:"required"`
// }
