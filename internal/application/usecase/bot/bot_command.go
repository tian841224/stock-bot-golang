package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/application/usecase/stock"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

type BotCommandUsecase interface {
	GetUseGuideMessage() string
	GetDailyMarketInfo(ctx context.Context, userType valueobject.UserType, count int) (string, error)
	GetStockPerformance(ctx context.Context, userType valueobject.UserType, symbol string) (string, error)
	GetStockPerformanceChart(ctx context.Context, symbol string) (*dto.ChartAsset, error)
	GetTopVolumeStock(ctx context.Context, userType valueobject.UserType) (string, error)
	GetStockPrice(ctx context.Context, userType valueobject.UserType, symbol string, date *time.Time) (string, error)
	GetStockRevenueChart(ctx context.Context, symbol string) (*dto.ChartAsset, error)
	GetHistoricalCandlesChart(ctx context.Context, symbol string) (*dto.ChartAsset, error)
	GetStockCompanyInfo(ctx context.Context, userType valueobject.UserType, symbol string) (string, error)
	GetStockNewsForLine(ctx context.Context, symbol string) (*dto.LineStockNewsMessage, error)
	GetStockNewsForTelegram(ctx context.Context, symbol string) (*dto.TgStockNewsMessage, error)
}

var _ port.BotCommandPort = (*botCommandUsecase)(nil)

type botCommandUsecase struct {
	botCommandPort     port.BotCommandPort
	marketDataUsecase  stock.MarketDataUsecase
	marketChartUsecase stock.MarketChartUsecase
	formatterPort      port.FormatterPort
}

func NewBotCommandUsecase(
	botCommandPort port.BotCommandPort,
	formatterPort port.FormatterPort,
	marketDataUsecase stock.MarketDataUsecase,
	marketChartUsecase stock.MarketChartUsecase,
) BotCommandUsecase {
	return &botCommandUsecase{
		botCommandPort:     botCommandPort,
		formatterPort:      formatterPort,
		marketDataUsecase:  marketDataUsecase,
		marketChartUsecase: marketChartUsecase,
	}
}

func (u *botCommandUsecase) GetUseGuideMessage() string {
	text := `台股機器人指令指南🤖

	📊 圖表指令
	- /k [股票代碼] - K線圖 (含月均價、最高最低價標示、成交量)
	- /p [股票代碼] - 股票績效圖表 (折線圖)
	- /r [股票代碼] - 月營收圖表 (柱狀圖+年增率折線)
	
	📈 股票資訊指令
	- /d [股票代碼] - 查詢當日收盤資訊 (可指定日期)
	- /d [股票代碼] [日期] - 查詢指定日期股價 (格式: YYYY-MM-DD)
	- /i [股票代碼] - 查詢公司資訊
	- /n [股票代碼] - 查詢股票新聞
	
	📊 市場總覽指令/
	- /m - 查詢最新大盤資訊 (預設1筆)
	- /m [數量] - 查詢指定筆數的大盤資訊
	- /t - 查詢當日交易量前20名
	
	🔔 訂閱管理
	- /add [股票代碼] - 新增訂閱股票
	- /del [股票代碼] - 刪除訂閱股票
	- /sub [項目] - 訂閱功能
	- /unsub [項目] - 取消訂閱功能
	- /list - 查詢已訂閱功能及股票
	
	💡 使用範例：
	/k 2330 - 台積電K線圖
	/p 0050 - 元大台灣50績效圖表
	/r 2330 - 台積電月營收圖表
	/d 2330 2025-01-15 - 查詢台積電指定日期股價
	/m 3 - 查詢最新3筆大盤資訊`
	return text
}

func (u *botCommandUsecase) GetDailyMarketInfo(ctx context.Context, userType valueobject.UserType, count int) (string, error) {
	marketData, err := u.marketDataUsecase.GetDailyMarketInfo(ctx, count)
	if err != nil {
		return "", err
	}

	if marketData == nil {
		return "", nil
	}

	return u.formatterPort.FormatDailyMarketInfo(marketData, userType), nil
}

func (u *botCommandUsecase) GetStockPerformance(ctx context.Context, userType valueobject.UserType, symbol string) (string, error) {
	stockPerformance, err := u.marketDataUsecase.GetStockPerformance(ctx, symbol)
	if err != nil {
		return "", err
	}

	if stockPerformance == nil {
		return "", nil
	}

	return u.formatterPort.FormatStockPerformance(stockPerformance.Name, stockPerformance.Symbol, &stockPerformance.Data, userType), nil
}

func (u *botCommandUsecase) GetStockPerformanceChart(ctx context.Context, symbol string) (*dto.ChartAsset, error) {
	chart, err := u.marketChartUsecase.GetPerformanceChart(ctx, symbol)
	if err != nil {
		return nil, err
	}

	if chart == nil {
		return nil, nil
	}

	return &dto.ChartAsset{
		Data:     chart.ChartData,
		FileName: fmt.Sprintf("⚡️%s(%s)-績效圖表", chart.StockName, symbol),
	}, nil
}

func (u *botCommandUsecase) GetTopVolumeStock(ctx context.Context, userType valueobject.UserType) (string, error) {
	items, err := u.marketDataUsecase.GetTopVolumeStock(ctx)
	if err != nil {
		return "", err
	}

	if items == nil {
		return "", nil
	}

	return u.formatterPort.FormatTopVolumeStock(items, userType), nil
}

func (u *botCommandUsecase) GetStockPrice(ctx context.Context, userType valueobject.UserType, symbol string, date *time.Time) (string, error) {
	price, err := u.marketDataUsecase.GetStockPrice(ctx, symbol, date)
	if err != nil {
		return "", err
	}

	if price == nil {
		return "", nil
	}

	return u.formatterPort.FormatStockPrice(price, userType), nil
}

func (u *botCommandUsecase) GetStockRevenueChart(ctx context.Context, symbol string) (*dto.ChartAsset, error) {
	chart, err := u.marketChartUsecase.GetRevenueChart(ctx, symbol)
	if err != nil {
		return nil, err
	}

	if chart == nil {
		return nil, nil
	}

	return &dto.ChartAsset{
		Data:     chart.ChartData,
		FileName: fmt.Sprintf("⚡️%s(%s)-月營收圖表", chart.StockName, symbol),
	}, nil
}

func (u *botCommandUsecase) GetHistoricalCandlesChart(ctx context.Context, symbol string) (*dto.ChartAsset, error) {
	chart, err := u.marketChartUsecase.GetHistoricalCandlesChart(ctx, symbol)
	if err != nil {
		return nil, err
	}

	if chart == nil {
		return nil, nil
	}

	return &dto.ChartAsset{
		Data:     chart.ChartData,
		FileName: fmt.Sprintf("⚡️%s(%s)-歷史K線圖", chart.StockName, symbol),
	}, nil
}

func (u *botCommandUsecase) GetStockCompanyInfo(ctx context.Context, userType valueobject.UserType, symbol string) (string, error) {
	companyInfo, err := u.marketDataUsecase.GetStockCompanyInfo(ctx, symbol)
	if err != nil {
		return "", err
	}

	if companyInfo == nil {
		return "", nil
	}

	return u.formatterPort.FormatStockCompanyInfo(companyInfo, userType), nil
}

func (u *botCommandUsecase) GetStockNewsForLine(ctx context.Context, symbol string) (*dto.LineStockNewsMessage, error) {
	news, err := u.marketDataUsecase.GetStockNews(ctx, symbol, 10)
	if err != nil {
		return nil, err
	}

	if news == nil || len(*news) == 0 {
		return &dto.LineStockNewsMessage{
			Text: fmt.Sprintf("⚡️%s 暫無新聞資料", symbol),
		}, nil
	}

	stockName := (*news)[0].StockName
	return u.formatterPort.FormatLineNewsMessage(*news, stockName, symbol), nil
}

func (u *botCommandUsecase) GetStockNewsForTelegram(ctx context.Context, symbol string) (*dto.TgStockNewsMessage, error) {
	news, err := u.marketDataUsecase.GetStockNews(ctx, symbol, 10)
	if err != nil {
		return nil, err
	}

	if news == nil || len(*news) == 0 {
		return &dto.TgStockNewsMessage{
			Text: fmt.Sprintf("⚡️%s 暫無新聞資料", symbol),
		}, nil
	}

	stockName := (*news)[0].StockName
	return u.formatterPort.FormatTelegramNewsMessage(*news, stockName, symbol), nil
}
