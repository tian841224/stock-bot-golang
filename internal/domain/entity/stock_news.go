package entity

import "time"

// StockNews 股票新聞資料
type StockNews struct {
	StockID     string
	Title       string
	Summary     string
	Link        string
	Source      string
	PublishedAt time.Time
}
