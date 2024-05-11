package handlers

import (
	"anon_chat_tg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type MessageHandler struct {
	bot *tgbotapi.BotAPI
	db *db.DB
}

func NewMessageHandler(bot *tgbotapi.BotAPI, db *db.DB) *MessageHandler {
	return &MessageHandler{
		bot: bot,
		db: db,
	}
}

func (h *MessageHandler) HandleMessage(tgUpdate tgbotapi.Update) {
	newMsg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, "HandleMessage")
	_, err := h.bot.Send(newMsg)
	if err != nil {
		log.Panic(err)
	}
}
