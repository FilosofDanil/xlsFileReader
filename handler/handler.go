package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xuri/excelize/v2"
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

	if update.Message.Document != nil {
		h.handleFileMessages(update, bot)
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, TextUnknownCommand)
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

	welcomeText := GetWelcomeText(username)
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

	welcomeText := GetWelcomeText(username)
	msg := tgbotapi.NewMessage(chatID, welcomeText)

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (h *Handler) handleDefaultText(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID

	instructionsText := GetInstructionsText()
	msg := tgbotapi.NewMessage(chatID, instructionsText)

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (h *Handler) handleFileMessages(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID
	document := update.Message.Document
	state := h.getState(chatID)

	log.Printf("Received file from chat %d: %s", chatID, document.FileName)

	if state != StateStart {
		log.Printf("User not in START state, ignoring file")
		return
	}

	if !h.isValidExcelFile(document.FileName) {
		log.Printf("Invalid file type received: %s", document.FileName)
		msg := tgbotapi.NewMessage(chatID, TextFileInvalidType)
		bot.Send(msg)
		return
	}

	fileURL, err := bot.GetFileDirectURL(document.FileID)
	if err != nil {
		log.Printf("Error getting file URL: %v", err)
		msg := tgbotapi.NewMessage(chatID, TextFileDownloadError)
		bot.Send(msg)
		return
	}

	filePath, err := h.downloadAndSaveFile(fileURL, document.FileName)
	if err != nil {
		log.Printf("Error downloading/saving file: %v", err)
		msg := tgbotapi.NewMessage(chatID, TextFileSaveError)
		bot.Send(msg)
		return
	}

	log.Printf("File saved successfully: %s", filePath)

	fileSizeKB := float64(document.FileSize) / 1024
	responseText := TextFileReceived
	responseText += fmt.Sprintf(TextFileName, document.FileName)
	responseText += fmt.Sprintf(TextFileSize, fileSizeKB)
	responseText += TextFileProcessing

	msg := tgbotapi.NewMessage(chatID, responseText)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}

	textContent, err := h.readExcelFile(filePath)
	if err != nil {
		log.Printf("Error reading Excel file: %v", err)
		msg := tgbotapi.NewMessage(chatID, TextFileReadError)
		bot.Send(msg)
		return
	}

	err = h.sendTextFileToUser(bot, chatID, textContent)
	if err != nil {
		log.Printf("Error sending file to user: %v", err)
		return
	}

	log.Printf("Successfully sent script.txt to user %d", chatID)
}

func (h *Handler) isValidExcelFile(fileName string) bool {
	lowerName := strings.ToLower(fileName)
	return strings.HasSuffix(lowerName, ".xls") || strings.HasSuffix(lowerName, ".xlsx")
}

func (h *Handler) downloadAndSaveFile(fileURL, fileName string) (string, error) {
	filesDir := "files"
	if err := os.MkdirAll(filesDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create files directory: %w", err)
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	filePath := filepath.Join(filesDir, fileName)
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

func (h *Handler) readExcelFile(filePath string) (string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	var content strings.Builder

	sheetList := f.GetSheetList()
	for _, sheetName := range sheetList {
		content.WriteString(fmt.Sprintf("=== Sheet: %s ===\n\n", sheetName))

		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Printf("Error reading sheet %s: %v", sheetName, err)
			continue
		}

		for rowIndex, row := range rows {
			content.WriteString(fmt.Sprintf("Row %d: ", rowIndex+1))
			for colIndex, cell := range row {
				if colIndex > 0 {
					content.WriteString(" | ")
				}
				content.WriteString(cell)
			}
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}

	return content.String(), nil
}

func (h *Handler) sendTextFileToUser(bot *tgbotapi.BotAPI, chatID int64, content string) error {
	fileBytes := tgbotapi.FileBytes{
		Name:  "script.txt",
		Bytes: []byte(content),
	}

	doc := tgbotapi.NewDocument(chatID, fileBytes)
	doc.Caption = TextFileProcessed

	_, err := bot.Send(doc)
	if err != nil {
		return fmt.Errorf("failed to send document: %w", err)
	}

	return nil
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
