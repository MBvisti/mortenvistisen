package main

import (
	"log"

	"github.com/google/uuid"
)

type MinLength int

// Validate implements validators.
func (m MinLength) Violated(val any) bool {
	id, ok := val.(uuid.UUID)
	log.Print(id, ok)
	return val != int(m)
}

func main() {
	// id := uuid.New()
	id := uuid.UUID{}
	MinLength(2).Violated(id)
}
