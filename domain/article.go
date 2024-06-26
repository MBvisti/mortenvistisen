package domain

import (
	"reflect"
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

	if err := article.Validate(ArticleValidations()); err != nil {
		return Article{}, err
	}

	return article, nil
}

var ArticleValidations = func() map[string][]Rule {
	return map[string][]Rule{
		"ID":          {RequiredRule},
		"Title":       {RequiredRule, MinLenRule(2)},
		"HeaderTitle": {RequiredRule, MinLenRule(2)},
		"Excerpt":     {RequiredRule, MinLenRule(130), MaxLenRule(160)},
		"ReadTime":    {RequiredRule},
		"Filename":    {RequiredRule},
	}
}

func (a Article) Validate(validations map[string][]Rule) error {
	val := reflect.ValueOf(a)
	typ := reflect.TypeOf(a)
	var errors []ValidationErr
	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i)
		name := typ.Field(i).Name

		errVal := ErrValidation{
			FieldValue: value,
			FieldName:  name,
		}

		for _, rule := range validations[name] {
			if rule.IsViolated(GetFieldValue(value)) {
				errVal.Violations = append(
					errVal.Violations,
					rule.Violation(),
				)
				errVal.ViolationsForHuman = append(
					errVal.ViolationsForHuman,
					rule.ViolationForHumans(name),
				)
			}
		}

		if len(errVal.Violations) > 0 {
			errors = append(errors, errVal)
		}
	}

	e := constructValidationErrors(errors...)
	if len(e) > 0 {
		return e
	}

	return nil
}

func (a *Article) Update(
	title string,
	headerTitle string,
	filename string,
	slug string,
	excerpt string,
	readtime int32,
	tags []Tag,
) error {
	a.Title = title
	a.HeaderTitle = headerTitle
	a.Filename = filename
	a.Slug = slug
	a.Excerpt = excerpt
	a.ReadTime = readtime
	a.Tags = tags

	if err := a.Validate(ArticleValidations()); err != nil {
		return err
	}

	return nil
}
