package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type MessageHandler struct {
	bot *tgbotapi.BotAPI
}

func NewMessageHandler(bot *tgbotapi.BotAPI) *MessageHandler {
	return &MessageHandler{
		bot: bot,
	}
}

func (h *MessageHandler) HandleMessage(tgUpdate tgbotapi.Update) {
	newMsg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, "HandleMessage")
	_, err := h.bot.Send(newMsg)
	if err != nil {
		log.Panic(err)
	}
}
