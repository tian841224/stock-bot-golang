package error

import (
	"errors"
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

// 預定義錯誤 (保持向後相容)
var (
	// ErrNotFound 代表資源不存在
	ErrNotFound = &DomainError{Code: "NOT_FOUND", Message: "資源不存在"}
	// ErrInvalidArgument 代表參數不合法
	ErrInvalidArgument = &DomainError{Code: "INVALID_ARGUMENT", Message: "參數不合法"}
	// ErrAlreadyExists 代表資源已存在
	ErrAlreadyExists = &DomainError{Code: "ALREADY_EXISTS", Message: "資源已存在"}
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

// NewMarketClosedError 建立市場休市錯誤
func NewMarketClosedError(market string) *DomainError {
	return &DomainError{
		Code:    "MARKET_CLOSED",
		Message: fmt.Sprintf("市場 %s 目前休市", market),
	}
}

// 使用者相關錯誤

// NewUserNotFoundError 建立使用者不存在錯誤
func NewUserNotFoundError(userID uint) *DomainError {
	return &DomainError{
		Code:    "USER_NOT_FOUND",
		Message: fmt.Sprintf("找不到使用者 ID: %d", userID),
	}
}

// NewInvalidUserTypeError 建立無效使用者類型錯誤
func NewInvalidUserTypeError(userType string) *DomainError {
	return &DomainError{
		Code:    "INVALID_USER_TYPE",
		Message: fmt.Sprintf("無效的使用者類型: %s", userType),
	}
}

// 訂閱相關錯誤

// NewSubscriptionNotFoundError 建立訂閱不存在錯誤
func NewSubscriptionNotFoundError(subscriptionID uint) *DomainError {
	return &DomainError{
		Code:    "SUBSCRIPTION_NOT_FOUND",
		Message: fmt.Sprintf("找不到訂閱 ID: %d", subscriptionID),
	}
}

// NewDuplicateSubscriptionError 建立重複訂閱錯誤
func NewDuplicateSubscriptionError(userID, featureID uint) *DomainError {
	return &DomainError{
		Code:    "DUPLICATE_SUBSCRIPTION",
		Message: fmt.Sprintf("使用者 %d 已訂閱功能 %d", userID, featureID),
	}
}

// IsNotFound 檢查錯誤是否為資源不存在錯誤
func IsNotFound(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Code == "NOT_FOUND" || domainErr.Code == "USER_NOT_FOUND" || domainErr.Code == "SUBSCRIPTION_NOT_FOUND"
	}
	return false
}

// IsInvalidArgument 檢查錯誤是否為參數不合法錯誤
func IsInvalidArgument(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Code == "INVALID_ARGUMENT" || domainErr.Code == "INVALID_STOCK_SYMBOL" || domainErr.Code == "INVALID_MARKET" || domainErr.Code == "INVALID_USER_TYPE"
	}
	return false
}
