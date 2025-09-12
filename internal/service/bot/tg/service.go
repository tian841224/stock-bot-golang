package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgService struct {
	botClient *tgbotapi.BotAPI
}

func NewTgService(botClient *tgbotapi.BotAPI) *TgService {
	return &TgService{botClient: botClient}
}

func (s *TgService) HandleUpdate(update *tgbotapi.Update) error {
	if update.Message != nil {
		_, err := s.botClient.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "你好！我收到了你的訊息："+update.Message.Text))
		if err != nil {
			return err
		}
	}
	return nil
}
