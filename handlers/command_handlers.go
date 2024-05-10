package handlers

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

type Message struct {
	Message string `json:"message"`
}

type CommandHandler struct {
	bot *tgbotapi.BotAPI
	messages map[string]Message
}

func NewCommandHandler(bot *tgbotapi.BotAPI) *CommandHandler {
	messagesFile, err := os.Open("handlers/messages.json")
	onFail("Failed to open file %v", err)
	defer messagesFile.Close()
	var messages map[string]Message
	err = json.NewDecoder(messagesFile).Decode(&messages)
	onFail("Failed to decode messages %v", err)

	return &CommandHandler{
		bot: bot,
		messages: messages,
	}
}

func (h *CommandHandler) HandleCommand(tgUpdate tgbotapi.Update) {
	switch tgUpdate.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, h.messages["start_without_city"].Message)
		_, err := h.bot.Send(msg)
		onFail("Failed to send message %v", err)
	}
}

func onFail(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
