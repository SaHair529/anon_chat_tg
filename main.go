package anon_chat_tg

import (
	"anon_chat_tg/config"
	"anon_chat_tg/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig()
	onFail("Failed to load config %v", err)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	onFail("Failed to create bot %v", err)

	messageHandler := handlers.NewMessageHandler(bot)
	commandHandler := handlers.NewCommandHandler(bot)
	callbackHandler := handlers.NewCallbackHandler(bot)

	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = 60

	updatesChan, err := bot.GetUpdatesChan(updates)
	onFail("Failed to create update channel %v", err)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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
