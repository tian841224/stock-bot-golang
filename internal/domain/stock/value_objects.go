package stock

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// StockID 股票代號值物件
type StockID struct {
	value string
}

// NewStockID 建立股票代號
func NewStockID(symbol string) (*StockID, error) {
	if symbol == "" {
		return nil, fmt.Errorf("股票代號不能為空")
	}

	// 清理股票代號（移除前綴等）
	cleanSymbol := strings.TrimSpace(symbol)
	cleanSymbol = strings.TrimPrefix(cleanSymbol, "TWS:")
	cleanSymbol = strings.TrimSuffix(cleanSymbol, ":STOCK")

	// 驗證股票代號格式
	if !isValidStockSymbol(cleanSymbol) {
		return nil, fmt.Errorf("無效的股票代號格式: %s", symbol)
	}

	return &StockID{value: cleanSymbol}, nil
}

// String 取得股票代號字串
func (s StockID) String() string {
	return s.value
}

// Equals 比較股票代號是否相等
func (s StockID) Equals(other StockID) bool {
	return s.value == other.value
}

// isValidStockSymbol 驗證股票代號格式
func isValidStockSymbol(symbol string) bool {
	// 台股代號通常是 4 位數字
	if len(symbol) != 4 {
		return false
	}

	// 檢查是否為數字
	_, err := strconv.Atoi(symbol)
	return err == nil
}

// Price 價格值物件
type Price struct {
	value float64
}

// NewPrice 建立價格
func NewPrice(value float64) (*Price, error) {
	if value < 0 {
		return nil, fmt.Errorf("價格不能為負數")
	}

	return &Price{value: value}, nil
}

// Value 取得價格數值
func (p Price) Value() float64 {
	return p.value
}

// Add 價格相加
func (p Price) Add(other Price) Price {
	return Price{value: p.value + other.value}
}

// Subtract 價格相減
func (p Price) Subtract(other Price) Price {
	return Price{value: p.value - other.value}
}

// Multiply 價格相乘
func (p Price) Multiply(factor float64) Price {
	return Price{value: p.value * factor}
}

// IsZero 檢查是否為零
func (p Price) IsZero() bool {
	return p.value == 0
}

// Percentage 百分比值物件
type Percentage struct {
	value float64
}

// NewPercentage 建立百分比
func NewPercentage(value float64) (*Percentage, error) {
	return &Percentage{value: value}, nil
}

// Value 取得百分比數值
func (p Percentage) Value() float64 {
	return p.value
}

// AsDecimal 取得小數形式
func (p Percentage) AsDecimal() float64 {
	return p.value / 100
}

// IsPositive 是否為正數
func (p Percentage) IsPositive() bool {
	return p.value > 0
}

// IsNegative 是否為負數
func (p Percentage) IsNegative() bool {
	return p.value < 0
}

// IsZero 是否為零
func (p Percentage) IsZero() bool {
	return p.value == 0
}

// String 格式化為字串
func (p Percentage) String() string {
	return fmt.Sprintf("%.2f%%", p.value)
}

// Volume 成交量值物件
type Volume struct {
	value int64
}

// NewVolume 建立成交量
func NewVolume(value int64) (*Volume, error) {
	if value < 0 {
		return nil, fmt.Errorf("成交量不能為負數")
	}

	return &Volume{value: value}, nil
}

// Value 取得成交量數值
func (v Volume) Value() int64 {
	return v.value
}

// Add 成交量相加
func (v Volume) Add(other Volume) Volume {
	return Volume{value: v.value + other.value}
}

// IsZero 檢查是否為零
func (v Volume) IsZero() bool {
	return v.value == 0
}

// String 格式化為字串（帶單位）
func (v Volume) String() string {
	if v.value >= 100000000 {
		return fmt.Sprintf("%.2f億", float64(v.value)/100000000)
	} else if v.value >= 10000 {
		return fmt.Sprintf("%.2f萬", float64(v.value)/10000)
	}
	return fmt.Sprintf("%d", v.value)
}

// TradingDate 交易日值物件
type TradingDate struct {
	date time.Time
}

// NewTradingDate 建立交易日
func NewTradingDate(date time.Time) *TradingDate {
	return &TradingDate{date: date}
}

// Date 取得日期
func (t TradingDate) Date() time.Time {
	return t.date
}

// IsToday 是否為今天
func (t TradingDate) IsToday() bool {
	now := time.Now()
	return t.date.Year() == now.Year() &&
		t.date.Month() == now.Month() &&
		t.date.Day() == now.Day()
}

// IsWeekend 是否為週末
func (t TradingDate) IsWeekend() bool {
	weekday := t.date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// Format 格式化日期
func (t TradingDate) Format(layout string) string {
	return t.date.Format(layout)
}

// String 預設格式化
func (t TradingDate) String() string {
	return t.date.Format("2006-01-02")
}
