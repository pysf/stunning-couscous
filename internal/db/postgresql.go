package db

import (
	"database/sql"
	"fmt"
	"os"
)

func NewPostgreConnection() (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRESQL_HOST"),
		os.Getenv("POSTGRESQL_PORT"),
		os.Getenv("POSTGRESQL_USERNAME"),
		os.Getenv("POSTGRESQL_PASSWORD"),
		os.Getenv("POSTGRESQL_DATABASE"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("NewPostgreConnection: connect to postgre err= %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("NewPostgreConnection: ping postgre err= %w", err)
	}

	return db, nil
}
