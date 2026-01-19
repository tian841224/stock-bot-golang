package bot

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	linebotInfra "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/line"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

// LineMessageProcessor 處理 LINE 訊息的路由和編排
type LineMessageProcessor struct {
	lineCommandUsecase LineCommandUsecase
	userAccountPort    port.UserAccountPort
	lineBotClient      *linebotInfra.LineBotClient
	logger             logger.Logger
}

func NewLineMessageProcessor(
	lineCommandUsecase LineCommandUsecase,
	userAccountPort port.UserAccountPort,
	lineBotClient *linebotInfra.LineBotClient,
	log logger.Logger,
) *LineMessageProcessor {
	return &LineMessageProcessor{
		lineCommandUsecase: lineCommandUsecase,
		userAccountPort:    userAccountPort,
		lineBotClient:      lineBotClient,
		logger:             log,
	}
}

// ProcessTextMessage 處理 LINE 文字訊息
func (p *LineMessageProcessor) ProcessTextMessage(ctx context.Context, event *linebot.Event, message *linebot.TextMessage) error {
	if message.Text == "" {
		return nil
	}

	userID := event.Source.UserID
	replyToken := event.ReplyToken
	messageText := message.Text

	p.logger.Info("收到 LINE 訊息",
		logger.String("user_id", userID),
		logger.String("message", messageText))

	// 確保使用者存在
	if err := p.ensureUser(ctx, userID); err != nil {
		p.logger.Error("確保使用者存在失敗", logger.Error(err))
	}

	// 解析命令和參數
	command, arg1, arg2 := p.parseMessageArgs(messageText)
	if command == "" {
		return p.lineBotClient.ReplyMessage(replyToken, "你說了: "+messageText)
	}

	// 路由到對應的命令處理器
	return p.routeCommand(ctx, command, arg1, arg2, replyToken)
}

// ensureUser 確保使用者存在，不存在則建立
func (p *LineMessageProcessor) ensureUser(ctx context.Context, userID string) error {
	_, err := p.userAccountPort.GetOrCreate(ctx, userID, valueobject.UserTypeLine)
	return err
}

// routeCommand 路由命令到對應的處理器
func (p *LineMessageProcessor) routeCommand(ctx context.Context, command, arg1, arg2, replyToken string) error {
	switch command {
	case "/start":
		return p.lineCommandUsecase.GetUseGuideMessage(replyToken)
	case "/k":
		return p.handleHistoricalCandles(ctx, replyToken, arg1)
	case "/p":
		return p.handlePerformanceChart(ctx, replyToken, arg1)
	case "/d":
		return p.handleStockPrice(ctx, replyToken, arg1, arg2)
	case "/t":
		return p.lineCommandUsecase.GetTopVolumeStock(ctx, replyToken)
	case "/i":
		return p.lineCommandUsecase.GetStockCompanyInfo(ctx, arg1, replyToken)
	case "/r":
		return p.handleRevenueChart(ctx, replyToken, arg1)
	case "/m":
		return p.handleDailyMarket(ctx, replyToken, arg1)
	case "/n":
		return p.lineCommandUsecase.GetStockNews(ctx, arg1, replyToken)
	default:
		return p.handleUnknownCommand(replyToken)
	}
}

// 各個命令的具體處理邏輯

func (p *LineMessageProcessor) handleHistoricalCandles(ctx context.Context, replyToken, symbol string) error {
	if symbol == "" {
		return p.sendError(replyToken, "請輸入股票代號\n\n使用方式：\n/k 股票代號 - 查詢K線圖")
	}
	return p.lineCommandUsecase.GetHistoricalCandlesChart(ctx, symbol, replyToken)
}

func (p *LineMessageProcessor) handlePerformanceChart(ctx context.Context, replyToken, symbol string) error {
	if symbol == "" {
		return p.sendError(replyToken, "請輸入股票代號\n\n使用方式：\n/p 股票代號 - 查詢績效圖表")
	}
	return p.lineCommandUsecase.GetStockPerformanceChart(ctx, symbol, replyToken)
}

func (p *LineMessageProcessor) handleStockPrice(ctx context.Context, replyToken, symbol, rawDate string) error {
	if symbol == "" {
		return p.sendError(replyToken, "請輸入股票代號\n\n使用方式：\n/d 股票代號 - 查詢今日股價\n/d 股票代號 2025-12-09 - 查詢指定日期股價")
	}

	var datePtr *time.Time
	if rawDate != "" {
		parsed, err := p.parseDate(rawDate)
		if err != nil {
			return p.sendError(replyToken, "日期格式錯誤，請使用 YYYY-MM-DD 格式\n例如：2025-12-09")
		}
		datePtr = &parsed
	} else {
		// 如果為空則取今天
		now := time.Now()
		datePtr = &now
		if now.Hour() < 14 {
			now = now.AddDate(0, 0, -1)
			datePtr = &now
		}
	}

	return p.lineCommandUsecase.GetStockPrice(ctx, symbol, datePtr, replyToken)
}

func (p *LineMessageProcessor) handleRevenueChart(ctx context.Context, replyToken, symbol string) error {
	if symbol == "" {
		return p.sendError(replyToken, "請輸入股票代號\n\n使用方式：\n/r 股票代號 - 查詢月營收圖表")
	}
	return p.lineCommandUsecase.GetStockRevenueChart(ctx, symbol, replyToken)
}

func (p *LineMessageProcessor) handleDailyMarket(ctx context.Context, replyToken, countStr string) error {
	count := 1
	if countStr != "" {
		countInt, err := strconv.Atoi(countStr)
		if countInt <= 0 || err != nil {
			return p.sendError(replyToken, "請輸入有效的數字，且大於0\n\n使用方式：\n/m [數量] - 查詢指定筆數的大盤資訊")
		}
		count = countInt
	}
	return p.lineCommandUsecase.GetDailyMarketInfo(ctx, replyToken, count)
}

func (p *LineMessageProcessor) handleUnknownCommand(replyToken string) error {
	return p.sendError(replyToken, "指令不存在，輸入 /start 查看說明")
}

// 輔助方法

func (p *LineMessageProcessor) sendError(replyToken, message string) error {
	p.logger.Warn("發送錯誤訊息", logger.String("message", message))
	return p.lineBotClient.ReplyMessage(replyToken, message)
}

func (p *LineMessageProcessor) parseMessageArgs(messageText string) (command, arg1, arg2 string) {
	parts := strings.Fields(messageText)
	if len(parts) == 0 {
		return "", "", ""
	}

	command = parts[0]
	if len(parts) > 1 {
		arg1 = parts[1]
	}
	if len(parts) > 2 {
		arg2 = parts[2]
	}
	return command, arg1, arg2
}

func (p *LineMessageProcessor) parseDate(value string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", value, time.Local)
}
