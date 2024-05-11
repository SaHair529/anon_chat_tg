package handlers

import (
	"anon_chat_tg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CallbackHandler struct {
	bot *tgbotapi.BotAPI
	db *db.DB
}

func NewCallbackHandler(bot *tgbotapi.BotAPI, db *db.DB) *CallbackHandler {
	return &CallbackHandler{
		bot: bot,
		db: db,
	}
}

func (cbh *CallbackHandler) HandleCallback(callback tgbotapi.Update) {}
