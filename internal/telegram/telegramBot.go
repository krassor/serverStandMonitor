package telegramBot

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type TelegramBot interface {
	Update(updateTimeout int)
}

type telegramBotImpl struct {
	tgbot *tgbotapi.BotAPI
}

func NewTgBotApi() TelegramBot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_APITOKEN"))
	if err != nil {
		log.Error().Msgf("Error auth telegram bot %s: %s", bot.Self.UserName, err)
	}

	bot.Debug = true

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	return &telegramBotImpl{
		tgbot: bot,
	}
}

func (tgBotImpl *telegramBotImpl) Update(updateTimeout int) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = updateTimeout

	updates := tgBotImpl.tgbot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			log.Info().Msgf("tgbot warn: Not command: %s", update.Message.Command())
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /list and /status."
		case "list":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := tgBotImpl.tgbot.Send(msg); err != nil {
			log.Error().Msgf("Error tgbot send message: %s", err)
		}
	}
}
