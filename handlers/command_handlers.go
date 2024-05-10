package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type CommandHandler struct {
	bot *tgbotapi.BotAPI
}

func NewCommandHandler(bot *tgbotapi.BotAPI) *CommandHandler {
	return &CommandHandler{
		bot: bot,
	}
}

func (h *CommandHandler) HandleCommand(tgUpdate tgbotapi.Update) {
	newMsg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, "HandleCommand")
	_, err := h.bot.Send(newMsg)
	if err != nil {
		log.Panic(err)
	}
}
