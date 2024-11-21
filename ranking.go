package main

import (
	"database/sql"
	"fmt"
	"time"
)

const rangeQuery = `select * from assets a inner join quotes q on a.id = q.asset_id WHERE q.timestamp between $1 AND $2 ORDER BY q.asset_id, q.timestamp`

type Ranker struct {
	db *sql.DB
}

// TimeRange holds data for querying data within a time interval.
type TimeRange struct {
	From time.Time
	To   time.Time
}

func (r *Ranker) calcRankChanges(timeRange TimeRange, limit int) ([]AssetQuote, error) {
	db := r.db
	rows, err := db.Query(rangeQuery, timeRange.From, timeRange.To)

	if err != nil {
		return nil, fmt.Errorf("failed to query ranks: %w", err)
	}

	var quotes []AssetQuote
	pq := newPriorityQueue[AssetQuote](limit, func(c1, c2 AssetQuote) int {
		return c1.RankChange - c2.RankChange
	})
	var curId int

	for rows.Next() {
		aq := AssetQuote{}
		err := rows.Scan(&aq.Asset.ID, &aq.TokenID, &aq.Name, &aq.Symbol, &aq.Asset.Timestamp, &aq.Quote.ID, &aq.AssetID, &aq.Rank, &aq.Price, &aq.MarketCap, &aq.PercentChange1H, &aq.PercentChange24H, &aq.PercentChange7D, &aq.PercentChange30D, &aq.PercentChange60D, &aq.PercentChange90D, &aq.Quote.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ranks: %w", err)
		}
		if curId == 0 {
			curId = aq.AssetID
		}
		if curId != aq.AssetID {
			recentQuote := quotes[len(quotes)-1]
			recentQuote.RankChange = calcRankChange(quotes)
			pq.Add(recentQuote)
			quotes = nil
			curId = aq.AssetID
		}
		quotes = append(quotes, aq)
	}
	return pq.GetAll(), nil
}

func calcRankChange(quotes []AssetQuote) int {
	if len(quotes) > 1 {
		return quotes[0].Rank - quotes[len(quotes)-1].Rank
	}
	return 0
}
