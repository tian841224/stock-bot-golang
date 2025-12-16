package stock

import (
	"context"
	"fmt"
	"testing"

	"github.com/tian841224/stock-bot/internal/application/dto"
	usecase_dto "github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/entity"
)

type mockMarketChartPort struct {
	GetRevenueChartFunc           func(ctx context.Context, symbol string) ([]byte, error)
	GetHistoricalCandlesChartFunc func(ctx context.Context, symbol string) ([]byte, string, error)
	GetPerformanceChartFunc       func(ctx context.Context, symbol string) (*dto.StockPerformanceChart, error)
}

func (m *mockMarketChartPort) GetRevenueChart(ctx context.Context, symbol string) ([]byte, error) {
	if m.GetRevenueChartFunc != nil {
		return m.GetRevenueChartFunc(ctx, symbol)
	}
	return nil, nil
}

func (m *mockMarketChartPort) GetHistoricalCandlesChart(ctx context.Context, symbol string) ([]byte, string, error) {
	if m.GetHistoricalCandlesChartFunc != nil {
		return m.GetHistoricalCandlesChartFunc(ctx, symbol)
	}
	return nil, "", nil
}

func (m *mockMarketChartPort) GetPerformanceChart(ctx context.Context, symbol string) (*dto.StockPerformanceChart, error) {
	if m.GetPerformanceChartFunc != nil {
		return m.GetPerformanceChartFunc(ctx, symbol)
	}
	return nil, nil
}

func TestMarketChartUsecase_GetRevenueChart(t *testing.T) {
	symbol := "2330"
	validStock := &entity.StockSymbol{Symbol: symbol, Name: "台積電"}
	chartData := []byte("chart data")

	tests := []struct {
		name             string
		symbol           string
		mockValidateFunc func(ctx context.Context, symbol string) (*entity.StockSymbol, error)
		mockChartFunc    func(ctx context.Context, symbol string) ([]byte, error)
		expectError      bool
		errorContains    string
	}{
		{
			name:   "成功取得營收圖表",
			symbol: symbol,
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return validStock, nil
			},
			mockChartFunc: func(ctx context.Context, symbol string) ([]byte, error) {
				return chartData, nil
			},
			expectError: false,
		},
		{
			name:   "驗證股票代號失敗",
			symbol: "invalid",
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return nil, fmt.Errorf("invalid symbol")
			},
			expectError:   true,
			errorContains: "查無此股票代號",
		},
		{
			name:   "取得圖表失敗",
			symbol: symbol,
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return validStock, nil
			},
			mockChartFunc: func(ctx context.Context, symbol string) ([]byte, error) {
				return nil, fmt.Errorf("port error")
			},
			expectError:   true,
			errorContains: "取得營收圖表失敗",
		},
		{
			name:   "圖表資料為nil",
			symbol: symbol,
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return validStock, nil
			},
			mockChartFunc: func(ctx context.Context, symbol string) ([]byte, error) {
				return nil, nil
			},
			expectError:   true,
			errorContains: "查無資料",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockValidation := &mockValidationPort{
				ValidateSymbolFunc: tt.mockValidateFunc,
			}
			mockMarketChart := &mockMarketChartPort{
				GetRevenueChartFunc: tt.mockChartFunc,
			}

			uc := NewMarketDataChartUsecase(mockMarketChart, mockValidation, &mockLogger{})
			_, err := uc.GetRevenueChart(context.Background(), tt.symbol)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望錯誤但沒有發生錯誤")
				}
				if tt.errorContains != "" && err != nil {
					if !containsString(err.Error(), tt.errorContains) {
						t.Errorf("錯誤訊息不符合期望，期望包含: %s, 實際: %s", tt.errorContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("不期望錯誤但發生錯誤: %v", err)
				}
			}
		})
	}
}

func TestMarketChartUsecase_GetHistoricalCandlesChart(t *testing.T) {
	query := usecase_dto.CandleQuery{Symbol: "2330"}
	chartData := []byte("chart data")

	tests := []struct {
		name          string
		query         usecase_dto.CandleQuery
		mockChartFunc func(ctx context.Context, symbol string) ([]byte, string, error)
		expectError   bool
		errorContains string
	}{
		{
			name:  "成功取得歷史K線圖",
			query: query,
			mockChartFunc: func(ctx context.Context, symbol string) ([]byte, string, error) {
				return chartData, "台積電", nil
			},
			expectError: false,
		},
		{
			name:  "取得歷史K線圖失敗",
			query: query,
			mockChartFunc: func(ctx context.Context, symbol string) ([]byte, string, error) {
				return nil, "", fmt.Errorf("port error")
			},
			expectError:   true,
			errorContains: "取得歷史K線圖失敗",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMarketChart := &mockMarketChartPort{
				GetHistoricalCandlesChartFunc: tt.mockChartFunc,
			}

			uc := NewMarketDataChartUsecase(mockMarketChart, nil, &mockLogger{})
			_, err := uc.GetHistoricalCandlesChart(context.Background(), tt.query.Symbol)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望錯誤但沒有發生錯誤")
				}
				if tt.errorContains != "" && err != nil {
					if !containsString(err.Error(), tt.errorContains) {
						t.Errorf("錯誤訊息不符合期望，期望包含: %s, 實際: %s", tt.errorContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("不期望錯誤但發生錯誤: %v", err)
				}
			}
		})
	}
}

func TestMarketChartUsecase_GetPerformanceChart(t *testing.T) {
	symbol := "2330"
	validStock := &entity.StockSymbol{Symbol: symbol, Name: "台積電"}
	chartData := &dto.StockPerformanceChart{
		Data:      []dto.StockPerformanceData{{Period: "1M", PeriodName: "一個月", Performance: "10.5"}},
		ChartData: []byte("chart data"),
	}

	tests := []struct {
		name             string
		symbol           string
		mockValidateFunc func(ctx context.Context, symbol string) (*entity.StockSymbol, error)
		mockChartFunc    func(ctx context.Context, symbol string) (*dto.StockPerformanceChart, error)
		expectError      bool
		errorContains    string
	}{
		{
			name:   "成功取得績效圖表",
			symbol: symbol,
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return validStock, nil
			},
			mockChartFunc: func(ctx context.Context, symbol string) (*dto.StockPerformanceChart, error) {
				return chartData, nil
			},
			expectError: false,
		},
		{
			name:   "驗證股票代號失敗",
			symbol: "invalid",
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return nil, fmt.Errorf("invalid symbol")
			},
			expectError:   true,
			errorContains: "查無此股票代號",
		},
		{
			name:   "取得績效圖表失敗",
			symbol: symbol,
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return validStock, nil
			},
			mockChartFunc: func(ctx context.Context, symbol string) (*dto.StockPerformanceChart, error) {
				return nil, fmt.Errorf("port error")
			},
			expectError:   true,
			errorContains: "查無資料",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockValidation := &mockValidationPort{
				ValidateSymbolFunc: tt.mockValidateFunc,
			}
			mockMarketChart := &mockMarketChartPort{
				GetPerformanceChartFunc: tt.mockChartFunc,
			}

			uc := NewMarketDataChartUsecase(mockMarketChart, mockValidation, &mockLogger{})
			_, err := uc.GetPerformanceChart(context.Background(), tt.symbol)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望錯誤但沒有發生錯誤")
				}
				if tt.errorContains != "" && err != nil {
					if !containsString(err.Error(), tt.errorContains) {
						t.Errorf("錯誤訊息不符合期望，期望包含: %s, 實際: %s", tt.errorContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("不期望錯誤但發生錯誤: %v", err)
				}
			}
		})
	}
}

// containsString 檢查字串是否包含子字串
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
