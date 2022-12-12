package telegramBot

import (
	"context"
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/serverStandMonitor/internal/services"
)

type Bot struct {
	tgbot   *tgbotapi.BotAPI
	service services.DevicesRepoService
}

func NewBot(service services.DevicesRepoService) *Bot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_APITOKEN"))
	if err != nil {
		log.Error().Msgf("Error auth telegram bot: %s", err)
	}
	//TODO: add to env BOTDEBUG
	bot.Debug = false

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		tgbot:   bot,
		service: service,
	}
}

func (bot *Bot) Update(ctx context.Context, updateTimeout int) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = updateTimeout

	updates := bot.tgbot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.Message == nil && update.CallbackQuery == nil { // ignore any non-Message updates
			log.Info().Msgf("tgbot warn: Not message: %s", update.Message)
			continue
		}

		if update.Message == nil && update.CallbackQuery != nil {
			err := bot.callbackQueryHandle(ctx, update.CallbackQuery)
			if err != nil {
				log.Error().Msgf("Error tgbot handle message: %s", err)
			}
			log.Info().Msgf("CallbackQuery from user: %s, data: %s", update.CallbackQuery.From, update.CallbackQuery.Data)
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			log.Info().Msgf("tgbot warn: Not command: %s", update.Message)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is not command")
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.tgbot.Send(msg); err != nil {
				log.Error().Msgf("Error tgbot send message: %s", err)
			}
			continue
		}

		log.Info().Msgf("tgbot receive command: %s", update.Message.Command())

		if err := bot.commandHandle(update.Message); err != nil {
			log.Error().Msgf("Error tgbot handle message: %s", err)
		}

	}
	log.Info().Msgf("exit telegram bot routine")
}

func (bot *Bot) commandHandle(msg *tgbotapi.Message) error {

	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "")
	replyMsg.ReplyToMessageID = msg.MessageID

	// Extract the command from the Message.

	switch msg.Command() {
	case "help":
		replyMsg.Text = "I understand /list"
	case "list":
		err := bot.list(&replyMsg)
		if err != nil {
			return err
		}
	case "start":
		replyMsg.Text = fmt.Sprintf("Hello, %s! I'm stand device monitor.\nEnter /list command and select device", msg.Chat.UserName)
	default:
		replyMsg.Text = "I don't know this command"
	}

	_, err := bot.tgbot.Send(replyMsg)
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Error shutdown telegram bot: %s", ctx.Err())
		default:
			bot.tgbot.StopReceivingUpdates()
		}
	}
}

func (bot *Bot) list(msg *tgbotapi.MessageConfig) error {

	devices, err := bot.service.GetDevices(context.Background())
	if err != nil {
		return err
	}

	var inlineKeyboardRow []tgbotapi.InlineKeyboardButton
	var inlineNumericKeyboard tgbotapi.InlineKeyboardMarkup

	for i, device := range devices {
		buttonText := fmt.Sprintf("%s %s %s:%s", device.DeviceVendor, device.DeviceName, device.DeviceIpAddress, device.DevicePort)
		buttonId := fmt.Sprintf("%d", device.ID)

		inlineKeyboardRow = append(inlineKeyboardRow, tgbotapi.InlineKeyboardButton{Text: buttonText, CallbackData: &buttonId})

		if ((i + 1) % 2) == 0 {
			inlineNumericKeyboard.InlineKeyboard = append(inlineNumericKeyboard.InlineKeyboard, inlineKeyboardRow)
			inlineKeyboardRow = nil
		}

	}

	if inlineKeyboardRow != nil {
		inlineNumericKeyboard.InlineKeyboard = append(inlineNumericKeyboard.InlineKeyboard, inlineKeyboardRow)
		inlineKeyboardRow = nil
	}

	msg.Text = "Select device:"
	msg.ReplyMarkup = inlineNumericKeyboard
	return nil
}

func (bot *Bot) callbackQueryHandle(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {

	id, err := strconv.Atoi(callbackQuery.Data)
	if err != nil {
		callback := tgbotapi.NewCallback(callbackQuery.ID, "Internal error")
		bot.tgbot.Request(callback)
		return err
	}

	deviceEntity, err := bot.service.GetDeviceById(ctx, uint(id))
	if err != nil {
		callback := tgbotapi.NewCallback(callbackQuery.ID, "Internal error")
		bot.tgbot.Request(callback)
		return err
	}

	var status string
	if deviceEntity.DeviceStatus == true {
		status = "ONLINE"
	} else {
		status = "OFFLINE"
	}

	callbackData := fmt.Sprintf(
		"Device %s %s is %s",
		deviceEntity.DeviceVendor,
		deviceEntity.DeviceName,
		status,
	)
	callback := tgbotapi.NewCallback(callbackQuery.ID, callbackData)
	_, err = bot.tgbot.Request(callback)
	if err != nil {
		return err
	}

	replyMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, callbackData)
	_, err = bot.tgbot.Send(replyMsg)
	if err != nil {
		return err
	}
	return nil
}
