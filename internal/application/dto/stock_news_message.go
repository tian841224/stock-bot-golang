package dto

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/v8/linebot"
)

type TgStockNewsMessage struct {
	Text                 string
	InlineKeyboardMarkup *tgbotapi.InlineKeyboardMarkup
}

type LineStockNewsMessage struct {
	Text            string
	CarouselColumns []*linebot.CarouselColumn
	FlexContainer   linebot.FlexContainer
	UseFlexMessage  bool
}
