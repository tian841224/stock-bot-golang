package twstock

import (
	"fmt"

	stockDto "stock-bot/internal/service/twstock/dto"
	"stock-bot/pkg/logger"

	"go.uber.org/zap"
)

// ========== 財報相關方法 ==========

// GetStockRevenue 取得股票財報
func (s *StockService) GetStockRevenue(stockID string) (*stockDto.RevenueDto, error) {
	logger.Log.Info("取得股票財報", zap.String("stockID", stockID))

	// 取得近12個月財報
	response, err := s.cnyesAPI.GetRevenue(stockID, 12)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("API回應錯誤: %s", response.Message)
	}

	return s.formatRevenue(response.Data), nil
}
