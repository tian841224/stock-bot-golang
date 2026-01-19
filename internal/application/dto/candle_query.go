package dto

import "time"

type CandleQuery struct {
	Symbol   string
	From     time.Time
	To       time.Time
	Interval string
	Fields   []string
}
