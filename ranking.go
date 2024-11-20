package main

import (
	"database/sql"
	"fmt"
	"time"
)

const rangeQuoteQuery = "SELECT * FROM quotes WHERE timestamp between $1 AND $2 ORDER BY asset_id, timestamp"

type RankChange struct {
	AssetID     int
	Change      int
	RecentQuote Quote
}

type Ranker struct {
	db *sql.DB
}

// TimeRange holds data for querying data within a time interval.
type TimeRange struct {
	From time.Time
	To   time.Time
}

func (r *Ranker) calcRankChanges(timeRange TimeRange, limit int) ([]RankChange, error) {
	db := r.db
	rows, err := db.Query(rangeQuoteQuery, timeRange.From, timeRange.To)

	if err != nil {
		return nil, fmt.Errorf("failed to query ranks: %w", err)
	}

	var quotes []Quote
	pq := newPriorityQueue[RankChange](limit, func(c1, c2 RankChange) int {
		return c1.Change - c2.Change
	})
	var curId int

	for rows.Next() {
		quote := Quote{}
		err := rows.Scan(&quote.ID, &quote.AssetID, &quote.Rank, &quote.Price, &quote.MarketCap, &quote.PercentChange1H, &quote.PercentChange24H, &quote.PercentChange7D, &quote.PercentChange30D, &quote.PercentChange60D, &quote.PercentChange90D, &quote.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ranks: %w", err)
		}
		if curId == 0 {
			curId = quote.AssetID
		}
		if curId != quote.AssetID {
			change := RankChange{AssetID: curId, Change: calcRankChange(quotes), RecentQuote: quotes[len(quotes)-1]}
			pq.Add(change)
			quotes = nil
			curId = quote.AssetID
		}
		quotes = append(quotes, quote)
	}
	return pq.GetAll(), nil
}

func calcRankChange(quotes []Quote) int {
	if len(quotes) > 1 {
		return quotes[0].Rank - quotes[len(quotes)-1].Rank
	}
	return 0
}
