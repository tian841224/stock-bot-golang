package logic_test

import (
	"testing"

	logic "github.com/tian841224/stock-bot/internal/domain/service"
)

func TestCalculateStockPerformance(t *testing.T) {
	tests := []struct {
		name             string
		todayPrice       float64
		prevPrice        float64
		wantChangeAmount float64
		wantChangeRate   float64
		wantUpDownSign   string
	}{
		{
			name:             "漲",
			todayPrice:       110.0,
			prevPrice:        100.0,
			wantChangeAmount: 10.0,
			wantChangeRate:   10.0,
			wantUpDownSign:   "+",
		},
		{
			name:             "跌",
			todayPrice:       90.0,
			prevPrice:        100.0,
			wantChangeAmount: -10.0,
			wantChangeRate:   -10.0,
			wantUpDownSign:   "-",
		},
		{
			name:             "平",
			todayPrice:       100.0,
			prevPrice:        100.0,
			wantChangeAmount: 0.0,
			wantChangeRate:   0.0,
			wantUpDownSign:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChangeAmount, gotChangeRate, gotUpDownSign := logic.CalculateStockPerformance(tt.todayPrice, tt.prevPrice)
			if gotChangeAmount != tt.wantChangeAmount {
				t.Errorf("CalculateStockPerformance() gotChangeAmount = %v, want %v", gotChangeAmount, tt.wantChangeAmount)
			}
			if gotChangeRate != tt.wantChangeRate {
				t.Errorf("CalculateStockPerformance() gotChangeRate = %v, want %v", gotChangeRate, tt.wantChangeRate)
			}
			if gotUpDownSign != tt.wantUpDownSign {
				t.Errorf("CalculateStockPerformance() gotUpDownSign = %v, want %v", gotUpDownSign, tt.wantUpDownSign)
			}
		})
	}
}
