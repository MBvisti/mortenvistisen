package models

import "errors"

var ErrNoRowWithIdentifier = errors.New("could not find requested row in database")

type ValidationErrsMap map[string]string
