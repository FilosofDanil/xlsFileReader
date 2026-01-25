package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotStatus struct {
	Status  string
	Error   error
	Message string
}

type Service struct {
	token        string
	statusChan   chan<- BotStatus
	messagesChan chan tgbotapi.Update
}

func NewService(token string, statusChan chan<- BotStatus) *Service {
	return &Service{
		token:        token,
		statusChan:   statusChan,
		messagesChan: make(chan tgbotapi.Update, 100),
	}
}

func (s *Service) Start(handleMessage func(tgbotapi.Update, *tgbotapi.BotAPI)) error {
	log.Println("Starting Telegram bot service...")

	bot, err := tgbotapi.NewBotAPI(s.token)
	if err != nil {
		s.notifyStatus("failed", err, "Failed to create bot instance")
		return fmt.Errorf("failed to create bot: %w", err)
	}

	log.Printf("Bot authorized on account: %s", bot.Self.UserName)
	s.notifyStatus("started", nil, fmt.Sprintf("Bot started successfully as @%s", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	log.Println("Bot is now listening for updates...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("Received message from user %s: %s", update.Message.From.UserName, update.Message.Text)
		handleMessage(update, bot)
	}

	s.notifyStatus("stopped", nil, "Bot stopped receiving updates")
	return nil
}

func (s *Service) notifyStatus(status string, err error, message string) {
	if s.statusChan != nil {
		select {
		case s.statusChan <- BotStatus{
			Status:  status,
			Error:   err,
			Message: message,
		}:
		default:
			log.Printf("Status channel is full, dropping status: %s", status)
		}
	}
}
