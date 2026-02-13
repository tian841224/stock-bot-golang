package entity

import (
	"regexp"

	domainerror "github.com/tian841224/stock-bot/internal/domain/error"
)

type StockSymbol struct {
	ID     uint
	Symbol string
	Name   string
	Market string
}

// Validate 驗證股票代號的合法性
func (s *StockSymbol) Validate() error {
	if s.Symbol == "" {
		return domainerror.ErrInvalidArgument
	}
	if s.Market == "" {
		return domainerror.ErrInvalidArgument
	}
	if !s.IsValidMarket() {
		return domainerror.NewInvalidMarketError(s.Market)
	}
	if !s.IsValidSymbol() {
		return domainerror.NewInvalidStockSymbolError(s.Symbol)
	}
	return nil
}

// IsValidMarket 檢查市場是否有效
func (s *StockSymbol) IsValidMarket() bool {
	validMarkets := map[string]bool{
		"TWSE": true,
		"TPEX": true,
		"US":   true,
	}
	return validMarkets[s.Market]
}

// IsValidSymbol 檢查股票代號格式是否有效
func (s *StockSymbol) IsValidSymbol() bool {
	if s.Market == "TWSE" || s.Market == "TPEX" {
		// 台股代號應為 4-6 位數字
		matched, _ := regexp.MatchString(`^[0-9]{4,6}$`, s.Symbol)
		return matched
	}
	if s.Market == "US" {
		// 美股代號應為 1-5 位英文字母
		matched, _ := regexp.MatchString(`^[A-Z]{1,5}$`, s.Symbol)
		return matched
	}
	return false
}

// IsTaiwanStock 檢查是否為台股
func (s *StockSymbol) IsTaiwanStock() bool {
	return s.Market == "TWSE" || s.Market == "TPEX"
}
