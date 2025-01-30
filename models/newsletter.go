package models

import (
	"time"

	"github.com/google/uuid"
)

type Newsletter struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string
	Content    string
	ReleasedAt time.Time
	Released   bool
}
