package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
)

var (
	password = os.Getenv("DB_PASSWORD")
	user     = os.Getenv("DB_USER")
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	dbName   = os.Getenv("DB_NAME")
)

func handler(ctx context.Context) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}(db)

	return nil
}