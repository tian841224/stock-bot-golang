package stock

import (
	"context"
	"fmt"
	"testing"

	"github.com/tian841224/stock-bot/internal/domain/entity"
)

func TestStockValidationUsecase_ValidateSymbol(t *testing.T) {
	tests := []struct {
		name        string
		symbol      string
		mockFunc    func(ctx context.Context, symbol string) (*entity.StockSymbol, error)
		expectError bool
	}{
		{
			name:   "成功驗證股票代號",
			symbol: "2330",
			mockFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return &entity.StockSymbol{
					Symbol: "2330",
					Name:   "台積電",
				}, nil
			},
			expectError: false,
		},
		{
			name:   "驗證股票代號失敗",
			symbol: "9999",
			mockFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return nil, fmt.Errorf("查詢失敗")
			},
			expectError: true,
		},
		{
			name:   "查無此股票代號",
			symbol: "9999",
			mockFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return nil, nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockValidation := &mockValidationPort{
				ValidateSymbolFunc: tt.mockFunc,
			}
			mockLog := &mockLogger{}
			validation := NewStockValidationUsecase(mockValidation, mockLog)
			stockSymbol, err := validation.ValidateSymbol(context.Background(), tt.symbol)
			if tt.expectError {
				if err == nil {
					t.Errorf("期望錯誤但沒有發生錯誤")
				}
			} else {
				if err != nil {
					t.Errorf("期望沒有錯誤但發生錯誤: %v", err)
				}

				if stockSymbol == nil {
					t.Errorf("期望有股票代號但為 nil")
				}
			}
		})
	}
}
