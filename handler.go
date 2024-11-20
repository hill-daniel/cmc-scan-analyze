package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	password        = os.Getenv("DB_PASSWORD")
	user            = os.Getenv("DB_USER")
	host            = os.Getenv("DB_HOST")
	port            = os.Getenv("DB_PORT")
	dbName          = os.Getenv("DB_NAME")
	emailRecipients = os.Getenv("EMAIL_RECIPIENTS")
	emailSender     = os.Getenv("EMAIL_SENDER")
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

	ranker := Ranker{db: db}
	timeRange := TimeRange{
		From: time.Now().Add(-10 * 24 * time.Hour),
		To:   time.Now(),
	}
	changes, err := ranker.calcRankChanges(timeRange, 10)
	if err != nil {
		return fmt.Errorf("failed to calculate rank changes: %v", err)
	}

	var builder strings.Builder
	writeChanges(&builder, changes)
	content := builder.String()
	htmlContent, err := createHtml(changes)
	if err != nil {
		return fmt.Errorf("failed to calculate rank changes: %v", err)
	}

	recipients := strings.Split(emailRecipients, ",")

	err = sendEmail(EmailInput{
		Sender:     emailSender,
		Recipients: recipients,
		Subject:    "New coinmarketcap rank changes",
		Body:       content,
		Html:       htmlContent,
	})
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
