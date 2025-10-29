package twstock

import (
	"fmt"

	twseDto "stock-bot/internal/infrastructure/twse/dto"
	"stock-bot/pkg/logger"

	"go.uber.org/zap"
)

// ========== 大盤資訊相關方法 ==========

// GetDailyMarketInfo 取得大盤資訊
func (s *StockService) GetDailyMarketInfo(count int) (twseDto.DailyMarketInfoResponseDto, error) {
	logger.Log.Info("取得大盤資訊", zap.Int("count", count))

	response, err := s.twseAPI.GetDailyMarketInfo()
	if err != nil {
		logger.Log.Error("呼叫 TWSE API 失敗", zap.Error(err))
		return twseDto.DailyMarketInfoResponseDto{}, err
	}

	if len(response.Data) == 0 {
		return twseDto.DailyMarketInfoResponseDto{}, fmt.Errorf("查無市場資料")
	}

	// 如果指定了筆數且小於總資料數，則從最後開始取指定筆數
	if count > 0 && count < len(response.Data) {
		originalCount := len(response.Data)
		// 取最後的 count 筆資料（從陣列末尾開始）
		startIndex := len(response.Data) - count
		response.Data = response.Data[startIndex:]
		logger.Log.Info("篩選最後資料", zap.Int("original", originalCount), zap.Int("filtered", count))
	}

	return response, nil
}
