package handler

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	userStates map[int64]string
}

func NewHandler() *Handler {
	return &Handler{
		userStates: make(map[int64]string),
	}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		h.handleCommandMessages(update, bot)
		return
	}

	h.handleTextMessages(update, bot)
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
	chatID := update.Message.Chat.ID
	username := update.Message.From.FirstName
	if username == "" {
		username = update.Message.From.UserName
	}

	h.setState(chatID, StateStart)

	welcomeText := "Hello, " + username + "! 👋\n\n"
	welcomeText += "Welcome to the XLS File Reader Bot!\n"
	welcomeText += "I'm here to help you process Excel files.\n\n"
	welcomeText += "📋 Available functions:\n"
	welcomeText += "• Send me an Excel file (.xls, .xlsx) to read and process\n"
	welcomeText += "• I will extract and display the data for you\n"
	welcomeText += "• Use /start to see this message again"

	msg := tgbotapi.NewMessage(chatID, welcomeText)

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Printf("Successfully responded to /start command from user: %s", username)
	}
}

func (h *Handler) handleTextMessages(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	log.Printf("Received text message from chat %d: %s", chatID, text)

	state := h.getState(chatID)

	switch state {
	case StateStart:
		h.handleStartStateText(update, bot)
	default:
		h.handleDefaultText(update, bot)
	}
}

func (h *Handler) handleStartStateText(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.FirstName
	if username == "" {
		username = update.Message.From.UserName
	}

	welcomeText := "Hello, " + username + "! 👋\n\n"
	welcomeText += "Welcome to the XLS File Reader Bot!\n"
	welcomeText += "I'm here to help you process Excel files.\n\n"
	welcomeText += "📋 Available functions:\n"
	welcomeText += "• Send me an Excel file (.xls, .xlsx) to read and process\n"
	welcomeText += "• I will extract and display the data for you\n"
	welcomeText += "• Use /start to see this message again"

	msg := tgbotapi.NewMessage(chatID, welcomeText)

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (h *Handler) handleDefaultText(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID

	instructionsText := "📖 Bot Instructions\n\n"
	instructionsText += "This bot helps you read and process Excel files.\n\n"
	instructionsText += "📋 Functions:\n"
	instructionsText += "• Send Excel files (.xls, .xlsx) - I will read and display the data\n"
	instructionsText += "• File processing - Extract information from your spreadsheets\n"
	instructionsText += "• Data display - View your Excel data in a readable format\n\n"
	instructionsText += "💡 To get started, use /start command or simply send me an Excel file!"

	msg := tgbotapi.NewMessage(chatID, instructionsText)

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (h *Handler) setState(chatID int64, state string) {
	h.userStates[chatID] = state
	log.Printf("State changed for chat %d: %s", chatID, state)
}

func (h *Handler) getState(chatID int64) string {
	if state, exists := h.userStates[chatID]; exists {
		return state
	}
	return StateDefault
}
