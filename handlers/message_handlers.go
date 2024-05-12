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
	conversation, err := h.db.GetUserConversation(tgUpdate.Message.Chat.ID)
	if err != nil {
		log.Printf("Failed to get conversation: %v", err)
	}

	if conversation.ID != 0 {
		msg := tgbotapi.NewMessage(conversation.OtherUserChatId, tgUpdate.Message.Text)
		_, err := h.bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send message %v: ", err)
		}
	}
}
