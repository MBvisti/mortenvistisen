package models

import "github.com/MBvisti/mortenvistisen/models/internal/database"

type User struct{}

type UserModel struct {
	db *database.Queries
}

func NewUser(db *database.Queries) UserModel {
	return UserModel{
		db: db,
	}
}
