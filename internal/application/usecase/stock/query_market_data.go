package stock

import (
	"context"
	"fmt"
	"time"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/application/port"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type MarketDataUsecase interface {
	GetDailyMarketInfo(ctx context.Context, count int) (*[]dto.DailyMarketInfo, error)
	GetStockPerformance(ctx context.Context, symbol string) (*dto.StockPerformance, error)
	GetTopVolumeStock(ctx context.Context) (*[]dto.TopVolume, error)
	GetStockPrice(ctx context.Context, symbol string, date *time.Time) (*dto.StockPrice, error)
	GetLatestTradeDate(ctx context.Context) (time.Time, error)
	GetLatestTradeDateByDateRange(ctx context.Context, startDate time.Time, endDate time.Time) ([]time.Time, error) // 修正：改為 []time.Time
	GetStockNews(ctx context.Context, symbol string, limit int) (*[]dto.StockNews, error)
	GetStockCompanyInfo(ctx context.Context, symbol string) (*dto.StockCompanyInfo, error)
}

type marketDataUsecase struct {
	market        port.MarketDataPort
	validation    port.ValidationPort
	tradeDateRepo port.TradeDateRepository
	logger        logger.Logger
}

func NewMarketDataUsecase(
	market port.MarketDataPort,
	validation port.ValidationPort,
	tradeDateRepo port.TradeDateRepository,
	logger logger.Logger,
) *marketDataUsecase {
	return &marketDataUsecase{market: market, validation: validation, tradeDateRepo: tradeDateRepo, logger: logger}
}

func (uc *marketDataUsecase) GetDailyMarketInfo(ctx context.Context, count int) (*[]dto.DailyMarketInfo, error) {
	result, err := uc.market.GetDailyMarketInfo(ctx, count)
	if err != nil {
		uc.logger.Error("取得大盤快照失敗", logger.Error(err))
		return nil, fmt.Errorf("查無資料，請確認後再試")
	}

	responseDto := make([]dto.DailyMarketInfo, len(*result))
	for i, data := range *result {
		responseDto[i] = dto.DailyMarketInfo{
			Date:        data.Date,
			Volume:      data.Volume,
			Amount:      data.Amount,
			Transaction: data.Transaction,
			Index:       data.Index,
			Change:      data.Change,
		}
	}
	return &responseDto, nil
}

func (uc *marketDataUsecase) GetStockPerformance(ctx context.Context, symbol string) (*dto.StockPerformance, error) {
	stock, err := uc.validation.ValidateSymbol(ctx, symbol)
	if err != nil || stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}
	result, err := uc.market.GetStockPerformance(ctx, stock.Symbol)
	if err != nil {
		uc.logger.Error("取得股票績效失敗", logger.Error(err))
		return nil, fmt.Errorf("取得績效資料失敗，請稍後再試")
	}
	if result == nil {
		return nil, fmt.Errorf("查無資料，請確認後再試")
	}
	return &dto.StockPerformance{
		Symbol: stock.Symbol,
		Name:   stock.Name,
		Data:   make([]dto.StockPerformanceData, len(result)),
	}, nil
}

func (uc *marketDataUsecase) GetTopVolumeStock(ctx context.Context) (*[]dto.TopVolume, error) {
	result, err := uc.market.GetTopVolumeStock(ctx)
	if err != nil {
		uc.logger.Error("取得交易量排行失敗", logger.Error(err))
		return nil, fmt.Errorf("查無資料，請確認後再試")
	}

	items := make([]dto.TopVolume, len(result))
	for i, item := range result {
		items[i] = dto.TopVolume{
			StockSymbol:      item.StockSymbol,
			StockName:        item.StockName,
			OpenPrice:        item.OpenPrice,
			ClosePrice:       item.ClosePrice,
			HighPrice:        item.HighPrice,
			LowPrice:         item.LowPrice,
			Volume:           item.Volume,
			Transaction:      item.Transaction,
			Amount:           item.Amount,
			ChangeAmount:     item.ChangeAmount,
			PercentageChange: item.PercentageChange,
			UpDownSign:       item.UpDownSign,
		}
	}
	return &items, nil
}

func (uc *marketDataUsecase) GetStockPrice(ctx context.Context, symbol string, date *time.Time) (*dto.StockPrice, error) {
	stock, err := uc.validation.ValidateSymbol(ctx, symbol)
	if err != nil || stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}

	// 如果 date 為 nil，使用當前時間
	if date == nil {
		now := time.Now()
		date = &now
	}

	// 判斷時間是否為兩點前，若是則取前一天
	if date.Hour() < 14 {
		prevDate := date.AddDate(0, 0, -1)
		date = &prevDate
	}

	// 取得最近30天交易日
	tradeDates, err := uc.market.GetLatestTradeDateByDateRange(ctx, date.AddDate(0, 0, -30), *date)

	if err != nil {
		uc.logger.Error("取得最近交易日失敗", logger.Error(err))
		return nil, fmt.Errorf("無法取得最近交易日")
	}

	tradeDate := tradeDates[len(tradeDates)-1]
	uc.logger.Info("取得股價資訊", logger.String("symbol", symbol), logger.Time("tradeDate", tradeDate))
	result, err := uc.market.GetStockPrice(ctx, symbol, &tradeDate)
	if err != nil {
		uc.logger.Error("取得股價資訊失敗", logger.Error(err))
		return nil, fmt.Errorf("查無資料，請確認後再試")
	}

	if len(*result) == 0 {
		return nil, fmt.Errorf("查無資料，請確認後再試")
	}

	tradeDateResult := (*result)[0]

	// 取得前一天收盤價
	prevTradeDate := tradeDates[len(tradeDates)-2]
	uc.logger.Info("取得股價資訊", logger.String("symbol", symbol), logger.Time("prevTradeDate", prevTradeDate))
	prevTradeDateResults, err := uc.market.GetStockPrice(ctx, symbol, &prevTradeDate)

	if err != nil {
		uc.logger.Error("取得股價資訊失敗", logger.Error(err))
		return nil, fmt.Errorf("查無資料，請確認後再試")
	}

	if len(*prevTradeDateResults) == 0 {
		return nil, fmt.Errorf("查無資料，請確認後再試")
	}

	prevTradeDateResult := (*prevTradeDateResults)[0]

	// 計算漲跌幅
	changeAmount := tradeDateResult.ClosePrice - prevTradeDateResult.ClosePrice
	changeRate := (changeAmount / prevTradeDateResult.ClosePrice) * 100
	upDownSign := ""
	if changeAmount > 0 {
		upDownSign = "+"
	} else {
		upDownSign = "-"
	}

	return &dto.StockPrice{
		Symbol:         tradeDateResult.Symbol,
		Name:           tradeDateResult.Name,
		Date:           tradeDateResult.Date,
		OpenPrice:      tradeDateResult.OpenPrice,
		ClosePrice:     tradeDateResult.ClosePrice,
		HighPrice:      tradeDateResult.HighPrice,
		LowPrice:       tradeDateResult.LowPrice,
		Volume:         tradeDateResult.Volume,
		Transactions:   tradeDateResult.Transactions,
		ChangeAmount:   changeAmount,
		ChangeRate:     changeRate,
		UpDownSign:     upDownSign,
		PrevClosePrice: prevTradeDateResult.ClosePrice,
	}, nil
}

func (uc *marketDataUsecase) GetLatestTradeDate(ctx context.Context) (time.Time, error) {
	tradeDate, err := uc.tradeDateRepo.GetByDateRange(ctx, time.Now().AddDate(0, 0, -30), time.Now())
	if err != nil {
		return time.Time{}, err
	}

	if len(tradeDate) == 0 {
		return time.Time{}, fmt.Errorf("找不到交易日資料")
	}

	now := time.Now()
	// 從最新的日期開始往回找
	for i := len(tradeDate) - 1; i >= 0; i-- {
		date := tradeDate[i]
		// 只比較日期部分（忽略時間）
		if date.Date.Year() == now.Year() &&
			date.Date.Month() == now.Month() &&
			date.Date.Day() == now.Day() {
			return date.Date, nil
		}
		// 找到最近一個過去的交易日
		if date.Date.Before(now) || date.Date.Equal(now.Truncate(24*time.Hour)) {
			return date.Date, nil
		}
	}

	return time.Time{}, fmt.Errorf("找不到交易日資料")
}

func (uc *marketDataUsecase) GetLatestTradeDateByDateRange(ctx context.Context, startDate time.Time, endDate time.Time) ([]time.Time, error) {
	tradeDate, err := uc.tradeDateRepo.GetByDateRange(ctx, startDate, endDate)
	if err != nil {
		return []time.Time{}, err
	}

	if len(tradeDate) == 0 {
		return []time.Time{}, fmt.Errorf("找不到交易日資料")
	}

	tradeDates := make([]time.Time, len(tradeDate))
	for i, date := range tradeDate {
		tradeDates[i] = date.Date
	}
	return tradeDates, nil
}

func (uc *marketDataUsecase) GetStockNews(ctx context.Context, symbol string, limit int) (*[]dto.StockNews, error) {
	stock, err := uc.validation.ValidateSymbol(ctx, symbol)
	if err != nil || stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}
	if limit <= 0 {
		limit = 10
	}
	articles, err := uc.market.GetStockNews(ctx, stock.Symbol)
	if err != nil {
		uc.logger.Error("取得股票新聞失敗", logger.Error(err))
		return nil, fmt.Errorf("取得新聞失敗，請稍後再試")
	}
	if len(articles) > limit {
		articles = articles[:limit]
	}

	items := make([]dto.StockNews, len(articles))
	for i, item := range articles {
		items[i] = dto.StockNews{
			Title:       item.Title,
			Date:        item.Date,
			StockSymbol: item.StockSymbol,
			StockName:   stock.Name,
			Link:        item.Link,
			Source:      item.Source,
		}
	}
	return &items, nil
}

func (uc *marketDataUsecase) GetStockCompanyInfo(ctx context.Context, symbol string) (*dto.StockCompanyInfo, error) {
	stock, err := uc.validation.ValidateSymbol(ctx, symbol)
	if err != nil || stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}
	companyInfo, err := uc.market.GetStockCompanyInfo(ctx, stock.Symbol)
	if err != nil {
		return nil, fmt.Errorf("取得公司資訊失敗，請稍後再試")
	}
	return companyInfo, nil
}
