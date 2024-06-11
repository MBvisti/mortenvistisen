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
	readtime int32,
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
		ReadTime:    readtime,
		Tags:        tags,
	}

	if err := article.Validate(); err != nil {
		return Article{}, err
	}

	return article, nil
}

var ArticleValidations = map[string][]Rule{
	"ID":          {RequiredRule},
	"Title":       {RequiredRule, MinLenRule(2)},
	"HeaderTitle": {RequiredRule, MinLenRule(2)},
	"Excerpt":     {RequiredRule, MinLenRule(130), MaxLenRule(160)},
	"ReadTime":    {RequiredRule},
	"Filename":    {RequiredRule},
}

func (a Article) Validate() error {
	var errors []ValidationErr
	for field, rules := range ArticleValidations {
		switch field {
		case "ID":
			idValidationErr := ErrValidation{
				FieldName:  "ID",
				FieldValue: a.ID,
			}
			for _, rule := range rules {
				if err := checkRule(a.ID, rule); err != nil {
					idValidationErr.Violations = append(
						idValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, idValidationErr)
		case "Title":
			titleValidationErr := ErrValidation{
				FieldName:  "Title",
				FieldValue: a.Title,
			}
			for _, rule := range rules {
				if err := checkRule(a.Title, rule); err != nil {
					titleValidationErr.Violations = append(
						titleValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, titleValidationErr)
		case "HeaderTitle":
			headerTitleValidationErr := ErrValidation{
				FieldName:  "HeaderTitle",
				FieldValue: a.HeaderTitle,
			}
			for _, rule := range rules {
				if err := checkRule(a.HeaderTitle, rule); err != nil {
					headerTitleValidationErr.Violations = append(
						headerTitleValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, headerTitleValidationErr)
		case "Excerpt":
			excerptValidationErr := ErrValidation{
				FieldName:  "Excerpt",
				FieldValue: a.Excerpt,
			}
			for _, rule := range rules {
				if err := checkRule(a.Excerpt, rule); err != nil {
					excerptValidationErr.Violations = append(
						excerptValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, excerptValidationErr)
		case "ReadTime":
			readTimeValidationErr := ErrValidation{
				FieldName:  "ReadTime",
				FieldValue: a.ReadTime,
			}
			for _, rule := range rules {
				if err := checkRule(a.ReadTime, rule); err != nil {
					readTimeValidationErr.Violations = append(
						readTimeValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, readTimeValidationErr)
		}
	}

	e := constructValidationErrors(
		errors...,
	)
	if len(e) > 0 {
		return e
	}

	return nil
}
