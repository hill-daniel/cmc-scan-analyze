package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	connStr := "postgres://myuser:mypassword@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ranker := Ranker{db: db}

	changes, err := ranker.calcRankChanges(TimeRange{
		From: time.Now().Add(-10 * 24 * time.Hour),
		To:   time.Now(),
	}, 10)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(changes)
	if len(changes) == 0 {
		t.Errorf("expected changes, got nothing")
	}
}
