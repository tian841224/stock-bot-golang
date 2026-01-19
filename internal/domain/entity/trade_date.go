package entity

import "time"

type TradeDate struct {
	ID       uint
	Date     time.Time
	Exchange string
}
