package main

import (
	"errors"
	"log"
	"time"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/google/uuid"
)

func main() {
	usr := domain.User{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Name:           "m",
		Mail:           "sdjkasldjklads@gmail.com",
		MailVerifiedAt: time.Now(),
		Password:       "dksjlkdjaljd",
	}

	vali := domain.BuildUserValidations("dksjlkdjaljd")

	var e domain.ValidationErrs
	if errors.As(usr.Validate(vali), &e) {
		for _, err := range e {
			log.Print(err.ErrorForHumans())
			log.Print(err.Error())
		}
	}
}
