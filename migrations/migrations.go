package migrations

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gdey/goose/v3"
)

func Up(conn *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := goose.Up(conn, fmt.Sprintf("%v/migrations", cwd)); err != nil {
		return err
	}

	return nil
}

func Status(conn *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := goose.Status(conn, fmt.Sprintf("%v/migrations", cwd)); err != nil {
		return err
	}

	return nil
}
