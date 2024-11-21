package main

import (
	"time"
)

type AssetQuote struct {
	RankChange int
	Asset
	Quote
}

// Asset represents the assets table.
type Asset struct {
	ID        int       `json:"id" db:"id"`
	TokenID   int64     `json:"token_id" db:"token_id"`
	Name      string    `json:"name" db:"name"`
	Symbol    string    `json:"symbol,omitempty" db:"symbol"` // Nullable
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

// Quote represents the quotes table.
type Quote struct {
	ID               int       `json:"id" db:"id"`
	AssetID          int       `json:"asset_id" db:"asset_id"`
	Rank             int       `json:"rank" db:"rank"`
	Price            float64   `json:"price,omitempty" db:"price"`                           // Nullable
	MarketCap        float64   `json:"market_cap,omitempty" db:"market_cap"`                 // Nullable
	PercentChange1H  float64   `json:"percent_change_1h,omitempty" db:"percent_change_1h"`   // Nullable
	PercentChange24H float64   `json:"percent_change_24h,omitempty" db:"percent_change_24h"` // Nullable
	PercentChange7D  float64   `json:"percent_change_7d,omitempty" db:"percent_change_7d"`   // Nullable
	PercentChange30D float64   `json:"percent_change_30d,omitempty" db:"percent_change_30d"` // Nullable
	PercentChange60D float64   `json:"percent_change_60d,omitempty" db:"percent_change_60d"` // Nullable
	PercentChange90D float64   `json:"percent_change_90d,omitempty" db:"percent_change_90d"` // Nullable
	Timestamp        time.Time `json:"timestamp" db:"timestamp"`
}
