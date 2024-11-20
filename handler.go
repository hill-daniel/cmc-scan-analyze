package main

import (
	"context"
	"database/sql"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"os"
	"text/tabwriter"
	"time"
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

	ranker := Ranker{db: db}
	timeRange := TimeRange{
		From: time.Now().Add(-10 * 24 * time.Hour),
		To:   time.Now(),
	}
	changes, err := ranker.calcRankChanges(timeRange, 10)
	if err != nil {
		return fmt.Errorf("failed to calculate rank changes: %v", err)
	}
	writeToConsole(changes)
	return nil
}

func writeToConsole(rankChanges []RankChange) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	_, err := fmt.Fprintln(w, "#\t CMC-Rank\t Name\t Symbol\t Change 24h\t Change 7d\t Change 30d\t Market Cap\t Price")
	if err != nil {
		log.Fatalf("failed to format header: %v", err)
	}
	p := message.NewPrinter(language.English)

	for i, change := range rankChanges {
		_, err = fmt.Fprintf(w, "#%d\t %d\t %s\t %s\t %.2f%%\t %.2f%%\t  %.2f%%\t %s\t %s \n", i+1, change.Change, change.RecentQuote.Name, change.RecentQuote.Symbol, change.RecentQuote.PercentChange24H,
			change.RecentQuote.PercentChange7D, change.RecentQuote.PercentChange30D, p.Sprintf("%.2f", change.RecentQuote.MarketCap), p.Sprintf("%.8f", change.RecentQuote.Price))
		if err != nil {
			log.Fatalf("failed to format rankChange: %v", err)
		}
	}
	err = w.Flush()
	if err != nil {
		log.Fatalf("failed to flush writer: %v", err)
	}
}
