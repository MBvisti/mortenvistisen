package main

import (
	"errors"
	"log"
	"time"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/google/uuid"
)

func main() {
	// Create individual validation errors
	now := time.Now()
	usr := domain.User{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "",
		Mail:      "",
		Password:  "dsajkldsajdljksaldjksa",
	}
	if err := usr.Validate("dsajkldsajdljksaldjksa"); err != nil {
		var validationErrs domain.ValidationErrs
		if errors.As(err, &validationErrs) {
			for _, valiErr := range validationErrs {
				log.Print(valiErr.Error())
			}
		} else {
			log.Print(usr)
		}
	}
}
