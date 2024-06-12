package main

import (
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
		Name:           "mm",
		Mail:           "sdjkasldjklads@gmail.com",
		MailVerifiedAt: time.Now(),
		Password:       "dksjlkdjaljd",
	}

	vali := domain.BuildUserValidations("dksjalkdjaljd")

	log.Print(usr.Validate(vali))
}
