package handler

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		h.handleCommandMessages(update, bot)
		return
	}

	log.Printf("Received non-command message: %s", update.Message.Text)
}

func (h *Handler) handleCommandMessages(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	command := update.Message.Command()

	switch command {
	case "start":
		h.handleStartCommand(update, bot)
	default:
		log.Printf("Unknown command: %s", command)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Use /start to begin.")
		bot.Send(msg)
	}
}

func (h *Handler) handleStartCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	username := update.Message.From.FirstName
	if username == "" {
		username = update.Message.From.UserName
	}

	welcomeText := "Hello, " + username + "! 👋\n\n"
	welcomeText += "Welcome to the XLS File Reader Bot!\n"
	welcomeText += "I'm here to help you process Excel files."

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeText)
	msg.ParseMode = "HTML"

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Printf("Successfully responded to /start command from user: %s", username)
	}
}
