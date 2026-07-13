package stock

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/entity"
)

func TestMarketDataUsecase_GetDailyMarketInfo(t *testing.T) {
	testCount := 4
	tests := []struct {
		name          string
		count         int
		mockFunc      func(ctx context.Context, count int) (*[]dto.DailyMarketInfo, error)
		expectError   bool
		errorContains string
	}{
		{
			name:  "成功取得大盤快照",
			count: testCount,
			mockFunc: func(ctx context.Context, count int) (*[]dto.DailyMarketInfo, error) {
				result := []dto.DailyMarketInfo{
					{Date: "2024-01-01"},
					{Date: "2024-01-02"},
					{Date: "2024-01-03"},
					{Date: "2024-01-04"},
				}
				return &result, nil
			},
			expectError: false,
		},
		{
			name:  "取得大盤快照失敗",
			count: testCount,
			mockFunc: func(ctx context.Context, count int) (*[]dto.DailyMarketInfo, error) {
				return nil, fmt.Errorf("取得大盤快照失敗")
			},
			expectError:   true,
			errorContains: "查無資料",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMarket := &mockMarketDataPort{
				GetDailyMarketInfoFunc: tt.mockFunc,
				GetLatestTradeDateByDateRangeFunc: func(ctx context.Context, startDate time.Time, endDate time.Time) ([]time.Time, error) {
					result := []time.Time{time.Now(), time.Now().AddDate(0, 0, -1)}
					return result, nil
				},
				GetStockNewsFunc: func(ctx context.Context, symbol string) ([]dto.StockNews, error) {
					result := []dto.StockNews{
						{Title: "測試新聞", Date: time.Now().Format("2006-01-02"), StockSymbol: "2330", Link: "https://www.google.com", Source: "測試來源"},
					}
					return result, nil
				},
				GetStockPriceFunc: func(ctx context.Context, symbol string, dates ...*time.Time) (*[]dto.StockPrice, error) {
					result := []dto.StockPrice{
						{Symbol: "2330", Name: "台積電", Date: time.Now(), OpenPrice: 100, ClosePrice: 101, HighPrice: 102, LowPrice: 99, Volume: 1000, Transactions: 1000, Amount: 1000000},
					}
					return &result, nil
				},
			}

			uc := NewMarketDataUsecase(mockMarket, nil, nil, &mockLogger{})

			result, err := uc.GetDailyMarketInfo(context.Background(), tt.count)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望錯誤但沒有發生錯誤")
				}
				if tt.errorContains != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errorContains) {
						t.Errorf("錯誤訊息不符合期望，期望包含: %s, 實際: %s", tt.errorContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("不期望錯誤但發生錯誤: %v", err)
				}
				if result == nil {
					t.Errorf("期望有結果但為 nil")
				}
			}
		})
	}
}

// TestMarketDataUsecase_GetStockPerformance 驗證績效資料會實際被帶入回傳結果，
// 防止先前「回傳長度正確但內容全為零值的空切片」的回歸。
func TestMarketDataUsecase_GetStockPerformance(t *testing.T) {
	expected := []dto.StockPerformanceData{
		{Period: "1d", PeriodName: "一日", Performance: "1.23%"},
		{Period: "1w", PeriodName: "一週", Performance: "-4.56%"},
	}

	mockMarket := &mockMarketDataPort{
		GetStockPerformanceFunc: func(ctx context.Context, symbol string) ([]dto.StockPerformanceData, error) {
			return expected, nil
		},
	}
	mockValidation := &mockValidationPort{
		ValidateSymbolFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
			return &entity.StockSymbol{Symbol: "2330", Name: "台積電"}, nil
		},
	}

	uc := NewMarketDataUsecase(mockMarket, mockValidation, nil, &mockLogger{})

	result, err := uc.GetStockPerformance(context.Background(), "2330")
	if err != nil {
		t.Fatalf("不期望錯誤但發生錯誤: %v", err)
	}
	if result == nil {
		t.Fatal("期望有結果但為 nil")
	}
	if result.Symbol != "2330" || result.Name != "台積電" {
		t.Errorf("股票代號/名稱不符，實際: %s / %s", result.Symbol, result.Name)
	}
	if len(result.Data) != len(expected) {
		t.Fatalf("績效資料筆數不符，期望 %d，實際 %d", len(expected), len(result.Data))
	}
	for i, want := range expected {
		if result.Data[i] != want {
			t.Errorf("第 %d 筆績效資料不符，期望 %+v，實際 %+v", i, want, result.Data[i])
		}
	}
}
