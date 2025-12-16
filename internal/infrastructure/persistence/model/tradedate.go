package models

import "time"

type TradeDate struct {
	Model
	// 交易所代號
	Exchange string `gorm:"column:exchange;type:varchar(255);not null" json:"exchange"`
	// 日期
	Date time.Time `gorm:"column:date;type:date;not null" json:"date"`
}

func (TradeDate) TableName() string {
	return "trade_dates"
}

func init() {
	RegisterModel(&TradeDate{})
}
