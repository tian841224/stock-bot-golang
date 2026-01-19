package stock

import (
	"context"
	"fmt"
	"time"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/application/port"
	fugle "github.com/tian841224/stock-bot/internal/infrastructure/external/stock/fugle"
	fugleDto "github.com/tian841224/stock-bot/internal/infrastructure/external/stock/fugle/dto"
	"github.com/tian841224/stock-bot/pkg/formatter"
	"github.com/tian841224/stock-bot/pkg/imageutil"
)

type marketChartGateway struct {
	marketDataPort port.MarketDataPort
	validationPort port.ValidationPort
	fugleAPI       *fugle.FugleAPI
}

func NewMarketChartGateway(marketDataPort port.MarketDataPort, validationPort port.ValidationPort, fugleAPI *fugle.FugleAPI) *marketChartGateway {
	return &marketChartGateway{
		marketDataPort: marketDataPort,
		validationPort: validationPort,
		fugleAPI:       fugleAPI,
	}
}

var _ port.MarketChartPort = (*marketChartGateway)(nil)

func (g *marketChartGateway) GetRevenueChart(ctx context.Context, symbol string) ([]byte, error) {

	data, err := g.marketDataPort.GetStockRevenue(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("取得營收圖表失敗: %w", err)
	}
	// 轉換為圖表資料格式
	chartData := g.convertStockRevenueToChartData(data)

	// 產生圖表
	chartBytes, err := imageutil.GenerateRevenueChart(chartData, data.StockName, data.StockSymbol)
	if err != nil {
		return nil, fmt.Errorf("產生營收圖表失敗: %v", err)
	}
	return chartBytes, nil
}

func (g *marketChartGateway) GetHistoricalCandlesChart(ctx context.Context, symbol string) ([]byte, string, error) {
	// 取得股票名稱
	stock, err := g.validationPort.ValidateSymbol(ctx, symbol)
	if err != nil {
		return nil, "", fmt.Errorf("查無此股票代號，請重新確認")
	}
	stockName := stock.Name

	response, err := g.fugleAPI.GetStockHistoricalCandles(fugleDto.FugleCandlesRequestDto{
		Symbol:    stock.Symbol,
		From:      time.Now().AddDate(-1, 0, 1).Format("2006-01-02"),
		Timeframe: "D",
		Fields:    "open,high,low,close,volume",
		Sort:      "desc",
	})

	if err != nil {
		return nil, "", err
	}

	if len(response.Data) == 0 {
		return nil, "", fmt.Errorf("查無K線資料")
	}

	// 轉換資料
	chartData := make([]imageutil.CandlestickData, len(response.Data))
	for i, d := range response.Data {
		chartData[i] = imageutil.CandlestickData{
			Date:   d.Date,
			Open:   d.Open,
			High:   d.High,
			Low:    d.Low,
			Close:  d.Close,
			Volume: d.Volume,
		}
	}

	// 產生圖表
	chartBytes, err := imageutil.GenerateCandlestickChart(chartData, stockName, stock.Symbol)
	if err != nil {
		return nil, stockName, fmt.Errorf("產生K線圖失敗: %v", err)
	}

	return chartBytes, stockName, nil
}

func (g *marketChartGateway) GetPerformanceChart(ctx context.Context, symbol string) (*dto.StockPerformanceChart, error) {
	stock, err := g.validationPort.ValidateSymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}

	stockPerformance, err := g.marketDataPort.GetStockPerformance(ctx, stock.Symbol)
	if err != nil {
		return nil, fmt.Errorf("取得股票績效失敗: %w", err)
	}

	// 轉換資料格式以供圖表使用
	chartData := make([]imageutil.PerformanceData, len(stockPerformance))
	for i, data := range stockPerformance {
		chartData[i] = imageutil.PerformanceData{
			Period:      data.Period,
			PeriodName:  data.PeriodName,
			Performance: data.Performance,
		}
	}

	// 生成圖表
	title := fmt.Sprintf("%s (%s) 績效表現", stock.Name, stock.Symbol)
	var chartBytes []byte

	// 只支援折線圖
	chartBytes, err = imageutil.GeneratePerformanceLineChart(chartData, title)

	if err != nil {
		return nil, fmt.Errorf("生成圖表失敗: %w", err)
	}

	return &dto.StockPerformanceChart{
		Symbol:    stock.Symbol,
		StockName: stock.Name,
		Data:      stockPerformance,
		ChartData: chartBytes,
	}, nil
}

// convertToChartData 轉換營收資料為圖表格式
func (g *marketChartGateway) convertStockRevenueToChartData(revenueData *dto.StockRevenue) []imageutil.RevenueChartData {
	if revenueData == nil || len(revenueData.Time) == 0 {
		return []imageutil.RevenueChartData{}
	}

	chartData := make([]imageutil.RevenueChartData, len(revenueData.Time))

	// 取得最新的營收和年增率
	latestRevenue := int64(0)
	latestYoY := 0.0
	if len(revenueData.SaleMonth) > 0 {
		latestRevenue = revenueData.SaleMonth[len(revenueData.SaleMonth)-1]
	}
	if len(revenueData.YoY) > 0 {
		latestYoY = revenueData.YoY[len(revenueData.YoY)-1]
	}

	for i, timestamp := range revenueData.Time {
		// 轉換時間戳記為日期格式

		period := formatter.FormatTimeFromTimestamp(timestamp)

		// 取得對應的資料
		revenue := int64(0)
		yoy := 0.0
		stockPrice := 0.0

		if i < len(revenueData.SaleMonth) {
			revenue = revenueData.SaleMonth[i]
		}
		if i < len(revenueData.YoY) {
			yoy = revenueData.YoY[i]
		}
		if i < len(revenueData.StockPrice) {
			stockPrice = revenueData.StockPrice[i]
		}

		chartData[i] = imageutil.RevenueChartData{
			Period:        period,
			PeriodName:    period,
			Revenue:       revenue,
			YoY:           yoy,
			StockPrice:    stockPrice,
			LatestRevenue: latestRevenue,
			LatestYoY:     latestYoY,
		}
	}

	return chartData
}
