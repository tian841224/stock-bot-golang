package formatter

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/tian841224/stock-bot/internal/application/dto"
)

type LineFormatter interface {
	FormatStockNews(news []dto.StockNews, stockName, symbol string) *dto.LineStockNewsMessage
}

// TelegramFormatter Telegram 訊息格式化器
type lineFormatter struct {
}

// NewTelegramFormatter 建立新的 Telegram 格式化器
func NewLineFormatter() *lineFormatter {
	return &lineFormatter{}
}

// FormatStockNews 格式化 Line 股票新聞訊息（使用 Flex Message）
func (f *lineFormatter) FormatStockNews(news []dto.StockNews, stockName, symbol string) *dto.LineStockNewsMessage {
	if len(news) == 0 {
		return &dto.LineStockNewsMessage{
			Text:           fmt.Sprintf("⚡️%s(%s)-即時新聞\n\n暫無新聞資料", stockName, symbol),
			UseFlexMessage: false,
		}
	}

	// Flex Message 建議最多顯示 10 則（避免訊息過大）
	maxItems := 10
	if len(news) > maxItems {
		news = news[:maxItems]
	}

	// 建立 Flex Message Bubble
	flexContainer := f.createNewsFlexMessage(news, stockName, symbol)

	return &dto.LineStockNewsMessage{
		Text:           fmt.Sprintf("⚡️%s(%s)-即時新聞", stockName, symbol),
		FlexContainer:  flexContainer,
		UseFlexMessage: true,
	}
}

// createNewsFlexMessage 建立新聞列表的 Flex Message
func (f *lineFormatter) createNewsFlexMessage(news []dto.StockNews, stockName, symbol string) *linebot.BubbleContainer {
	// 建立標題區塊
	header := &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   fmt.Sprintf("⚡️ %s (%s)", stockName, symbol),
				Weight: linebot.FlexTextWeightTypeBold,
				Size:   linebot.FlexTextSizeTypeLg,
				Color:  "#1DB446",
			},
			&linebot.TextComponent{
				Type:  linebot.FlexComponentTypeText,
				Text:  "即時新聞",
				Size:  linebot.FlexTextSizeTypeSm,
				Color: "#999999",
			},
		},
		PaddingAll: "15px",
	}

	// 建立新聞列表
	newsItems := make([]linebot.FlexComponent, 0, len(news))
	for i, n := range news {
		// 標題限制 100 字元
		title := n.Title
		if len([]rune(title)) > 100 {
			title = string([]rune(title)[:97]) + "..."
		}

		// 新聞項目
		newsItem := &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				// 新聞編號與標題
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeBaseline,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   fmt.Sprintf("%d.", i+1),
							Size:   linebot.FlexTextSizeTypeSm,
							Color:  "#1DB446",
							Weight: linebot.FlexTextWeightTypeBold,
							Flex:   linebot.IntPtr(0),
						},
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   title,
							Size:   linebot.FlexTextSizeTypeSm,
							Wrap:   true,
							Color:  "#111111",
							Flex:   linebot.IntPtr(1),
							Action: &linebot.URIAction{
								Label: "查看",
								URI:   n.Link,
							},
						},
					},
					Spacing: linebot.FlexComponentSpacingTypeSm,
				},
				// 日期與來源
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeBaseline,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "  ",
							Flex:  linebot.IntPtr(0),
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  fmt.Sprintf("📅 %s  |  📰 %s", n.Date, n.Source),
							Size:  linebot.FlexTextSizeTypeXs,
							Color: "#999999",
							Flex:  linebot.IntPtr(1),
						},
					},
				},
			},
			PaddingAll: "10px",
			Margin:     linebot.FlexComponentMarginTypeMd,
		}

		// 新增分隔線（除了最後一個）
		if i < len(news)-1 {
			newsItem.Contents = append(newsItem.Contents, &linebot.SeparatorComponent{
				Type:   linebot.FlexComponentTypeSeparator,
				Margin: linebot.FlexComponentMarginTypeMd,
			})
		}

		newsItems = append(newsItems, newsItem)
	}

	// 主體內容
	body := &linebot.BoxComponent{
		Type:     linebot.FlexComponentTypeBox,
		Layout:   linebot.FlexBoxLayoutTypeVertical,
		Contents: newsItems,
		Spacing:  linebot.FlexComponentSpacingTypeNone,
	}

	// 組合成 Bubble
	bubble := &linebot.BubbleContainer{
		Type:   linebot.FlexContainerTypeBubble,
		Header: header,
		Body:   body,
	}

	return bubble
}
