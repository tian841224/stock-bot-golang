package presenter

import shared "github.com/tian841224/stock-bot/internal/application/dto"

// 對外輸出的股票相關回應模型，透過 type alias 共用核心 DTO。
type (
	StockPriceResponse         = shared.StockPrice
	StockPerformanceResponse   = shared.StockPerformance
	StockPerformanceChartAsset = shared.StockPerformanceChart
	StockCompanyInfoResponse   = shared.StockCompanyInfo
	TopVolumeResponse          = shared.TopVolume
	KlineCandlesResponse       = shared.KlineCandles
)
