package domain

import "github.com/google/uuid"

type Tag struct {
	ID   uuid.UUID
	Name string
}

func NewTag(name string) (Tag, error) {
	tag := Tag{ID: uuid.New(), Name: name}

	if err := tag.Validate(); err != nil {
		return Tag{}, err
	}

	return tag, nil
}

func (t Tag) Validate() error {
	nameValidationErr := ErrValidation{
		FieldName:  "Name",
		FieldValue: t.Name,
	}
	if t.Name == "" {
		nameValidationErr.Violations = append(nameValidationErr.Violations, ErrIsRequired)
	}
	if len(t.Name) < 2 {
		nameValidationErr.Violations = append(nameValidationErr.Violations, ErrTooShort)
	}

	e := constructValidationErrors(
		nameValidationErr,
	)
	if len(e) > 0 {
		return e
	}

	return nil
}
