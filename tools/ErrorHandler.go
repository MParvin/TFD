package tools

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleError logs errors to stderr without exiting
func HandleError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

// HandleFatalError logs errors and exits the program
func HandleFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// SendErrorToAdmin sends error messages to the admin via Telegram
func SendErrorToAdmin(bot *tgbotapi.BotAPI, adminChatID int64, err error) {
	if err == nil || adminChatID == 0 {
		return
	}

	msg := tgbotapi.NewMessage(adminChatID, "⚠️ TFD Error: "+err.Error())
	_, sendErr := bot.Send(msg)
	if sendErr != nil {
		log.Printf("Failed to send error to admin: %v", sendErr)
	}
}
