package telegramBot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
