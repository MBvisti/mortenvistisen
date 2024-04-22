package entity

import (
	"time"

	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Tag struct {
	ID   uuid.UUID
	Name string
}

type Post struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Filename    string
	Slug        string
	Excerpt     string
	Draft       bool
	ReleaseDate time.Time
	ReadTime    int32
	Tags        []Tag
}

type NewPost struct {
	Title             string `validate:"required,gte=3"`
	HeaderTitle       string `validate:"required"`
	Excerpt           string `validate:"required,lte=160,gte=130"`
	ReleaseNow        bool
	EstimatedReadTime int64  `validate:"required"`
	Filename          string `validate:"required"`
}

func FilenameValidation(sl validator.StructLevel) {
	data := sl.Current().Interface().(NewPost)

	if err := posts.CheckIfFileExist(data.Filename); err != nil {
		sl.ReportError(
			data.Filename,
			"",
			"Filename",
			"",
			"filename not in assets",
		)
	}
}
