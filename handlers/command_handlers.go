package handlers

import (
	"anon_chat_tg/db"
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
	db *db.DB
}

func NewCommandHandler(bot *tgbotapi.BotAPI, db *db.DB) *CommandHandler {
	messagesFile, err := os.Open("handlers/messages.json")
	onFail("Failed to open file %v", err)
	defer messagesFile.Close()
	var messages map[string]Message
	err = json.NewDecoder(messagesFile).Decode(&messages)
	onFail("Failed to decode messages %v", err)

	return &CommandHandler{
		bot: bot,
		messages: messages,
		db: db,
	}
}

func (h *CommandHandler) HandleCommand(tgUpdate tgbotapi.Update) {
	switch tgUpdate.Message.Command() {
	case "start":
		commandArgs := tgUpdate.Message.CommandArguments()
		if len(commandArgs) == 0 {
			msg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, h.messages["start_without_city"].Message)
			_, err := h.bot.Send(msg)
			onFail("Failed to send message %v", err)
		}	else {
			if h.db.IsUserAlreadyInQueue(tgUpdate.Message.Chat.ID) {
				msg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, h.messages["already_on_queue"].Message)
				_, err := h.bot.Send(msg)
				onFail("Failed to send message %v", err)
				return
			}

			users, err := h.db.GetUsersFromQueueByCity(commandArgs)
			onFail("Failed to get users by city from db %v", err)
			if len(users) == 0 {
				h.db.AddUserToQueue(tgUpdate.Message.Chat.ID, commandArgs)
				msg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, h.messages["added_to_queue"].Message)
				_, err := h.bot.Send(msg)
				onFail("Failed to send message %v", err)
			} else {
				h.db.BeginConversation(tgUpdate.Message.Chat.ID, users[0].ChatId)

				msg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, h.messages["conversation_begin"].Message)
				_, err := h.bot.Send(msg)
				onFail("Failed to send message %v", err)

				msg.ChatID = users[0].ChatId
				_, err = h.bot.Send(msg)
				onFail("Failed to send message %v", err)
			}
		}
	case "stop":
		if h.db.IsUserAlreadyInQueue(tgUpdate.Message.Chat.ID) {
			h.db.DeleteUserFromQueue(tgUpdate.Message.Chat.ID)
			msg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, h.messages["deleted_from_queue"].Message)
			_, err := h.bot.Send(msg)
			onFail("Failed to send message %v", err)
		} else if h.db.IsUserHasConversation(tgUpdate.Message.Chat.ID) {
			conversation, err := h.db.GetUserConversation(tgUpdate.Message.Chat.ID)
			if err != nil {
				log.Printf("Failed to get conversation: %v", err)
			}

			h.db.DeleteUserConversation(tgUpdate.Message.Chat.ID)

			msg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, h.messages["deleted_conversation"].Message)
			_, err = h.bot.Send(msg)
			onFail("Failed to send message %v", err)

			msg = tgbotapi.NewMessage(conversation.OtherUserChatId, h.messages["companion_stopped_conversation"].Message)
			_, err = h.bot.Send(msg)
			onFail("Failed to send message %v", err)
		} else {
			msg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, h.messages["nothing_to_stop"].Message)
			_, err := h.bot.Send(msg)
			onFail("Failed to send message %v", err)
		}
	}
}

func onFail(message string, err error) {
	if err != nil {
		log.Printf(message, err)
	}
}
