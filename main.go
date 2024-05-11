package main

import (
	"anon_chat_tg/config"
	db2 "anon_chat_tg/db"
	"anon_chat_tg/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config/config.json")
	onFail("Failed to load config %v", err)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	onFail("Failed to create bot %v", err)

	db, err := db2.NewDB()
	onFail("Failed to connect db %v", err)

	messageHandler := handlers.NewMessageHandler(bot, db)
	commandHandler := handlers.NewCommandHandler(bot)
	callbackHandler := handlers.NewCallbackHandler(bot)

	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = 60

	updatesChan, err := bot.GetUpdatesChan(updates)
	onFail("Failed to create update channel %v", err)

	for {
		select {
		case update := <-updatesChan:
			if update.Message != nil {
				if update.Message.IsCommand() {
					commandHandler.HandleCommand(update)
				} else {
					messageHandler.HandleMessage(update)
				}
			} else if update.CallbackQuery != nil {
				callbackHandler.HandleCallback(update)
			}
		}
	}
}

func onFail(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
