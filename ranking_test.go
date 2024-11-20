package main

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_console_printout(t *testing.T) {
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
	if len(changes) == 0 {
		t.Errorf("expected changes, got nothing")
	}
	writeChanges(os.Stdout, changes)
}

func Test_send_email(t *testing.T) {
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
	if len(changes) == 0 {
		t.Errorf("expected changes, got nothing")
	}
	var builder strings.Builder
	writeChanges(&builder, changes)
	content := builder.String()
	htmlContent, err := createHtml(changes)
	if err != nil {
		t.Fatal(err)
	}

	err = sendEmail(EmailInput{
		Sender:     "kiwilisk@gmail.com",
		Recipients: []string{"daniel@uphill.dev"},
		Subject:    "New coinmarketcap rank changes",
		Body:       content,
		Html:       htmlContent,
	})
	if err != nil {
		t.Errorf("failed to send email: %v", err)
	}
}
