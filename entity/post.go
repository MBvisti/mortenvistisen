package entity

import (
	"time"

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
