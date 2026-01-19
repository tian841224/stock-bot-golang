package formatter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	"github.com/tian841224/stock-bot/pkg/formatter"
	"github.com/tian841224/stock-bot/pkg/utils"
)

type formatterAdapter struct {
	marketChartPort   port.MarketChartPort
	validationPort    port.ValidationPort
	telegramFormatter TelegramFormatter
	lineFormatter     LineFormatter
}

func NewFormatterAdapter(marketChartPort port.MarketChartPort, validationPort port.ValidationPort, telegramFormatter TelegramFormatter, lineFormatter LineFormatter) *formatterAdapter {
	return &formatterAdapter{
		marketChartPort:   marketChartPort,
		validationPort:    validationPort,
		telegramFormatter: telegramFormatter,
		lineFormatter:     lineFormatter,
	}
}

var _ port.FormatterPort = (*formatterAdapter)(nil)

func (f *formatterAdapter) FormatDailyMarketInfo(data *[]dto.DailyMarketInfo, userType valueobject.UserType) string {
	var messageText strings.Builder

	if userType == valueobject.UserTypeTelegram {
		messageText.WriteString("<b>台灣股市大盤資訊</b>\n\n")
	} else {
		messageText.WriteString("台灣股市大盤資訊\n\n")
	}

	for _, row := range *data {
		date := row.Date
		volume := row.Volume
		amount := row.Amount
		transaction := row.Transaction
		index := row.Index
		change := row.Change

		if userType == valueobject.UserTypeTelegram {
			messageText.WriteString(fmt.Sprintf("<b>%s</b>\n", date))
			messageText.WriteString("<code>")
			messageText.WriteString(fmt.Sprintf("成交股數：%s\n", volume))
			messageText.WriteString(fmt.Sprintf("成交金額：%s\n", formatter.FormatAmountInt(utils.ToInt64(amount))))
			messageText.WriteString(fmt.Sprintf("成交筆數：%s\n", formatter.FormatAmountInt(utils.ToInt64(transaction))))
			messageText.WriteString(fmt.Sprintf("發行量加權股價指數：%s\n", index))
			messageText.WriteString(fmt.Sprintf("漲跌點數：%s\n", change))
			messageText.WriteString("</code>\n")
		} else {
			messageText.WriteString(fmt.Sprintf("%s\n", date))
			messageText.WriteString(fmt.Sprintf("成交股數：%s\n", volume))
			messageText.WriteString(fmt.Sprintf("成交金額：%s\n", formatter.FormatAmountInt(utils.ToInt64(amount))))
			messageText.WriteString(fmt.Sprintf("成交筆數：%s\n", formatter.FormatAmountInt(utils.ToInt64(transaction))))
			messageText.WriteString(fmt.Sprintf("發行量加權股價指數：%s\n", index))
			messageText.WriteString(fmt.Sprintf("漲跌點數：%s\n", change))
			messageText.WriteString("\n")
		}
	}
	return messageText.String()
}

func (f *formatterAdapter) FormatStockPerformance(stockName, symbol string, data *[]dto.StockPerformanceData, userType valueobject.UserType) string {
	var result strings.Builder

	if userType == valueobject.UserTypeTelegram {
		result.WriteString("<pre>")
		result.WriteString(fmt.Sprintf("📊 <b>%s (%s) 績效表現</b>\n\n", stockName, symbol))
	} else {
		result.WriteString(fmt.Sprintf("📊 %s (%s) 績效表現\n\n", stockName, symbol))
	}

	// 為每個績效期間添加表情符號和格式化
	for _, data := range *data {
		// 解析績效數值來決定表情符號
		performanceStr := strings.TrimSuffix(data.Performance, "%")
		performance, err := strconv.ParseFloat(performanceStr, 64)
		var emoji string
		if err == nil {
			if performance >= 0 {
				emoji = "📈"
			} else {
				emoji = "📉"
			}
		} else {
			emoji = "📊"
		}

		// 格式化顯示
		if userType == valueobject.UserTypeTelegram {
			result.WriteString(fmt.Sprintf("%s <b>%s</b>: %s\n", emoji, data.Period, data.Performance))
		} else {
			result.WriteString(fmt.Sprintf("%s %s: %s\n", emoji, data.Period, data.Performance))
		}
	}

	if userType == valueobject.UserTypeTelegram {
		result.WriteString("</pre>")
	}

	return result.String()
}

// FormatStockInfoMessage 格式化股票詳細資訊
func (f *formatterAdapter) FormatStockCompanyInfo(data *dto.StockCompanyInfo, userType valueobject.UserType) string {
	var message strings.Builder

	if userType == valueobject.UserTypeTelegram {
		message.WriteString("<pre>")
	}

	// 股票基本資訊
	message.WriteString("🏢" + data.Name)
	message.WriteString(" (")
	message.WriteString(data.Symbol)
	message.WriteString(")")
	message.WriteString(" | ")
	message.WriteString(data.Industry)
	message.WriteString(" | ")
	message.WriteString(data.Market)
	message.WriteString("\n\n")

	// 財務指標
	message.WriteString("💼財務指標:\n")
	message.WriteString("本益比: ")
	message.WriteString(fmt.Sprintf("%.2f", data.PE))
	message.WriteString("\n本淨比: ")
	message.WriteString(fmt.Sprintf("%.2f", data.PB))
	message.WriteString("\n市值: ")
	marketCapStr := fmt.Sprintf("%.2f", data.MarketCap/1000000000000)
	message.WriteString(marketCapStr)
	message.WriteString(" 兆\n每股淨值: ")
	message.WriteString(fmt.Sprintf("%.2f", data.BookValue))
	message.WriteString("\n近四季EPS: ")
	message.WriteString(fmt.Sprintf("%.2f", data.EPS))
	message.WriteString("\n營季EPS: ")
	message.WriteString(fmt.Sprintf("%.2f", data.QuarterEPS))
	message.WriteString("\n年股利: ")
	message.WriteString(fmt.Sprintf("%.2f", data.Dividend))
	message.WriteString("\n殖利率: ")
	message.WriteString(fmt.Sprintf("%.2f%%", data.DividendRate))
	message.WriteString("\n\n")

	// 獲利能力
	message.WriteString("💡獲利能力:\n")
	message.WriteString("毛利率: ")
	message.WriteString(fmt.Sprintf("%.2f%%", data.GrossMargin))
	message.WriteString("\n營益率: ")
	message.WriteString(fmt.Sprintf("%.2f%%", data.OperMargin))
	message.WriteString("\n淨利率: ")
	message.WriteString(fmt.Sprintf("%.2f%%", data.NetMargin))

	if userType == valueobject.UserTypeTelegram {
		message.WriteString("</pre>")
	}

	return message.String()
}

func (f *formatterAdapter) FormatTopVolumeStock(data *[]dto.TopVolume, userType valueobject.UserType) string {
	var messageText strings.Builder

	if userType == valueobject.UserTypeTelegram {
		messageText.WriteString("🔝<b>今日交易量前二十</b>\n\n")
	} else {
		messageText.WriteString("🔝今日交易量前二十\n\n")
	}

	for _, item := range *data {
		emoji := ""
		switch item.UpDownSign {
		case "+":
			emoji = "📈"
		case "-":
			emoji = "📉"
		default:
			emoji = ""
		}

		if userType == valueobject.UserTypeTelegram {
			messageText.WriteString(fmt.Sprintf("%s<b>%s (%s)</b>\n<code>", emoji, item.StockName, item.StockSymbol))
			messageText.WriteString(fmt.Sprintf("成交股數：%s\n", item.Volume))
			messageText.WriteString(fmt.Sprintf("成交筆數：%s\n", formatter.FormatAmountInt(utils.ToInt64(item.Transaction))))
			messageText.WriteString(fmt.Sprintf("開盤價：%.2f\n", item.OpenPrice))
			messageText.WriteString(fmt.Sprintf("收盤價：%.2f\n", item.ClosePrice))
			messageText.WriteString(fmt.Sprintf("漲跌幅：%s%.2f (%s)\n", item.UpDownSign, item.ChangeAmount, item.PercentageChange))
			messageText.WriteString(fmt.Sprintf("最高價：%.2f\n", item.HighPrice))
			messageText.WriteString(fmt.Sprintf("最低價：%.2f\n", item.LowPrice))
			messageText.WriteString("</code>\n")
		} else {
			messageText.WriteString(fmt.Sprintf("%s%s (%s)\n", emoji, item.StockName, item.StockSymbol))
			messageText.WriteString(fmt.Sprintf("成交股數：%s\n", item.Volume))
			messageText.WriteString(fmt.Sprintf("成交筆數：%s\n", formatter.FormatAmountInt(utils.ToInt64(item.Transaction))))
			messageText.WriteString(fmt.Sprintf("開盤價：%.2f\n", item.OpenPrice))
			messageText.WriteString(fmt.Sprintf("收盤價：%.2f\n", item.ClosePrice))
			messageText.WriteString(fmt.Sprintf("漲跌幅：%s%.2f (%s)\n", item.UpDownSign, item.ChangeAmount, item.PercentageChange))
			messageText.WriteString(fmt.Sprintf("最高價：%.2f\n", item.HighPrice))
			messageText.WriteString(fmt.Sprintf("最低價：%.2f\n", item.LowPrice))
			messageText.WriteString("\n")
		}
	}

	return messageText.String()
}

func (f *formatterAdapter) FormatStockPrice(data *dto.StockPrice, userType valueobject.UserType) string {
	displayDate := data.Date.Format("2006/01/02")

	emoji := ""
	switch data.UpDownSign {
	case "+":
		emoji = "📈"
	case "-":
		emoji = "📉"
	default:
		emoji = ""
	}

	if userType == valueobject.UserTypeTelegram {
		return fmt.Sprintf(`<b>%s</b>
			<b>─── %s (%s) %s ───</b><code>
開盤價：%.2f
收盤價：%.2f
漲跌幅：%.2f (%.2f%%)
最高價：%.2f
最低價：%.2f
交易量：%s
成交筆數：%s 張
		</code>`,
			displayDate,
			data.Name, data.Symbol, emoji,
			data.OpenPrice,
			data.ClosePrice,
			data.ChangeAmount, data.ChangeRate,
			data.HighPrice,
			data.LowPrice,
			formatter.FormatAmountInt(data.Volume),
			strconv.FormatInt(data.Transactions, 10))
	} else {
		return fmt.Sprintf(`%s
─── %s (%s) %s ───
開盤價：%.2f
收盤價：%.2f
漲跌幅：%.2f (%.2f%%)
最高價：%.2f
最低價：%.2f
交易量：%s
成交筆數：%s 張`,
			displayDate,
			data.Name, data.Symbol, emoji,
			data.OpenPrice,
			data.ClosePrice,
			data.ChangeAmount, data.ChangeRate,
			data.HighPrice,
			data.LowPrice,
			formatter.FormatAmountInt(data.Volume),
			strconv.FormatInt(data.Transactions, 10))
	}
}

func (f *formatterAdapter) FormatStockRevenue(data *dto.StockRevenue, userType valueobject.UserType) string {
	var message strings.Builder

	if userType == valueobject.UserTypeTelegram {
		message.WriteString(fmt.Sprintf("<b>📊 %s(%s) 月營收</b>\n\n", data.StockName, data.StockSymbol))
	} else {
		message.WriteString(fmt.Sprintf("📊 %s(%s) 月營收\n\n", data.StockName, data.StockSymbol))
	}

	// 檢查是否有資料
	if len(data.SaleMonth) == 0 || len(data.YoY) == 0 {
		message.WriteString("❌ 暫無營收資料")
		return message.String()
	}

	if userType == valueobject.UserTypeTelegram {
		message.WriteString("<pre>")
	}
	// 顯示所有資料
	for i := 0; i < len(data.Time); i++ {
		timeStr := formatter.FormatTimeFromTimestamp(data.Time[i])

		// 營收(千元) -> 億元
		monthRevenueE := float64(data.SaleMonth[i]) / 100000.0

		// 年增率
		yoy := data.YoY[i]

		// 累計營收(千元) -> 億元
		accumulatedRevenueE := float64(data.SaleAccumulated[i]) / 100000.0

		// 累計年增率
		accumulatedYoY := data.YoYAccumulated[i]

		message.WriteString(fmt.Sprintf("---%s---\n", timeStr))
		message.WriteString(fmt.Sprintf("營收(億元): %.2f\n", monthRevenueE))
		message.WriteString(fmt.Sprintf("年增率: %.2f%%\n", yoy))
		message.WriteString(fmt.Sprintf("累計營收(億元): %.2f\n", accumulatedRevenueE))
		message.WriteString(fmt.Sprintf("累計年增率: %.2f%%\n\n", accumulatedYoY))
	}
	if userType == valueobject.UserTypeTelegram {
		message.WriteString("</pre>")
	}

	return message.String()
}

// FormatTelegramNewsMessage 格式化 Telegram 股票新聞訊息
func (f *formatterAdapter) FormatTelegramNewsMessage(news []dto.StockNews, stockName, symbol string) *dto.TgStockNewsMessage {
	return f.telegramFormatter.FormatStockNews(news, stockName, symbol)
}

// FormatLineNewsMessage 格式化 Line 股票新聞訊息
func (f *formatterAdapter) FormatLineNewsMessage(news []dto.StockNews, stockName, symbol string) *dto.LineStockNewsMessage {
	return f.lineFormatter.FormatStockNews(news, stockName, symbol)
}

// FormatChartCaption 格式化圖表標題
func (f *formatterAdapter) FormatChartCaption(stockName, symbol, chartType string) string {
	return fmt.Sprintf("⚡️%s(%s)-%s", stockName, symbol, chartType)
}

// FormatSubscribed 格式化訂閱股票和項目
func (f *formatterAdapter) FormatSubscribed(stocks []*dto.UserSubscriptionStock, items []*dto.UserSubscriptionItem) string {
	// 組合訊息
	messageText := "📋 <b>您目前的訂閱項目</b>\n\n"

	// 訂閱功能清單
	messageText += "🔔 <b>已訂閱功能：</b>\n"
	if len(items) > 0 {
		for _, sub := range items {
			messageText += fmt.Sprintf("• %s\n", sub.Item.GetName())
		}
	} else {
		messageText += "• 尚未訂閱任何功能\n"
	}

	// 訂閱股票清單
	messageText += "\n📈 <b>已訂閱股票：</b>\n"
	if len(stocks) > 0 {
		for _, stock := range stocks {
			if stock.Status {
				messageText += fmt.Sprintf("• %s (%s)\n", stock.Name, stock.Symbol)
			}
		}
	} else {
		messageText += "• 尚未訂閱任何股票\n"
	}

	return messageText
}
