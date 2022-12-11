package telegramBot

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func TgBotInit() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_APITOKEN"))
	if err != nil {
		log.Error().Msgf("Error creating telegram bot: %s", err)
	}

	bot.Debug = true
}
