package error

import (
	"fmt"
)

// DomainError 領域錯誤基礎類型
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// 預定義錯誤
var (
	// ErrInvalidArgument 代表參數不合法
	ErrInvalidArgument = &DomainError{Code: "INVALID_ARGUMENT", Message: "參數不合法"}
)

// 股票相關錯誤

// NewInvalidStockSymbolError 建立無效股票代號錯誤
func NewInvalidStockSymbolError(symbol string) *DomainError {
	return &DomainError{
		Code:    "INVALID_STOCK_SYMBOL",
		Message: fmt.Sprintf("無效的股票代號: %s", symbol),
	}
}

// NewInvalidMarketError 建立無效市場錯誤
func NewInvalidMarketError(market string) *DomainError {
	return &DomainError{
		Code:    "INVALID_MARKET",
		Message: fmt.Sprintf("不支援的市場: %s", market),
	}
}

// NewInvalidUserTypeError 建立無效使用者類型錯誤
func NewInvalidUserTypeError(userType string) *DomainError {
	return &DomainError{
		Code:    "INVALID_USER_TYPE",
		Message: fmt.Sprintf("無效的使用者類型: %s", userType),
	}
}
