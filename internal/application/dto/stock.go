package dto

import (
	"time"

	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

// MarketIndexSnapshot 表示單日的大盤快照。
type MarketIndexSnapshot struct {
	Date         time.Time
	Volume       float64
	Amount       float64
	Transactions int64
	WeightedIdx  float64
	Change       float64
}

// VolumeLeaderboardEntry 包含交易量排行資料。
type VolumeLeaderboardEntry struct {
	Symbol           string
	Name             string
	OpenPrice        float64
	ClosePrice       float64
	HighPrice        float64
	LowPrice         float64
	ChangeAmount     float64
	PercentageChange float64
	Volume           int64
	Transactions     int64
	UpDownSign       string
}

// // StockProfile 包含股票的基本資料。
// type StockProfile struct {
// 	Symbol      string
// 	Name        string
// 	Industry    string
// 	Market      string
// 	Description string
// 	Website     string

// 	CurrentPrice float64
// 	Change       float64
// 	ChangeRate   float64
// 	OpenPrice    float64
// 	HighPrice    float64
// 	LowPrice     float64
// 	PrevClose    float64

// 	Volume      int64
// 	Turnover    float64
// 	VolumeRatio float64
// 	Amplitude   float64

// 	PE           float64
// 	PB           float64
// 	MarketCap    float64
// 	BookValue    float64
// 	EPS          float64
// 	QuarterEPS   float64
// 	Dividend     float64
// 	DividendRate float64
// 	GrossMargin  float64
// 	OperMargin   float64
// 	NetMargin    float64

// 	UpperLimit  float64
// 	LowerLimit  float64
// 	High52W     float64
// 	Low52W      float64
// 	High52WDate time.Time
// 	Low52WDate  time.Time

// 	BidPrices []float64
// 	AskPrices []float64
// 	OutVolume int64
// 	InVolume  int64
// 	OutRatio  float64
// 	InRatio   float64
// }

// // RevenueRecord 表示單月營收資料。
// type RevenueRecord struct {
// 	Month          time.Time
// 	MonthlyRevenue float64
// 	MonthlyYoY     float64
// 	AccumRevenue   float64
// 	AccumYoY       float64
// }

// // RevenueSeries 聚合營收序列。
// type RevenueSeries struct {
// 	Symbol  string
// 	Name    string
// 	Records []RevenueRecord
// }

// UserAccount 表示 Bot 端的使用者。
type UserAccount struct {
	ID        uint
	AccountID string
	Type      valueobject.UserType
	Status    bool
}

// SubscriptionPreference 使用者訂閱偏好。
type SubscriptionPreference struct {
	Key         string
	DisplayName string
	Enabled     bool
}

// StockSubscription 股票訂閱狀態。
type StockSubscription struct {
	Symbol    string
	Name      string
	CreatedAt time.Time
}
