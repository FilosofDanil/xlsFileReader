package main

import (
	"example/hello/bot"
	"example/hello/handler"
	"example/hello/logger"
	"log"
	"time"
)

const (
	MaxRetries = 3
	RetryDelay = 15 * time.Second
)

func main() {
	logger.Init()
	log.Println("Application starting...")

	// Load environment variables from .env file
	bot.LoadEnv()

	// Get bot token from environment
	telegramBotToken := bot.GetBotToken()
	if telegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required but not set")
	}

	stopChan := make(chan struct{})

	go supervisor(stopChan, telegramBotToken)

	select {
	case <-stopChan:
		log.Println("Received stop signal. Shutting down application.")
	}
}

func supervisor(stopChan chan<- struct{}, token string) {
	statusChan := make(chan bot.BotStatus, 10)
	failureCount := 0

	for {
		log.Printf("Launch attempt %d/%d", failureCount+1, MaxRetries)

		botService := bot.NewService(token, statusChan)
		messageHandler := handler.NewHandler()

		go func() {
			err := botService.Start(messageHandler.HandleUpdate)
			if err != nil {
				log.Printf("Bot crashed with error: %v", err)
				statusChan <- bot.BotStatus{
					Status:  bot.StatusCrashed,
					Error:   err,
					Message: "Bot encountered an error",
				}
			}
		}()

		crashed := false

		select {
		case status := <-statusChan:
			switch status.Status {
			case bot.StatusStarted:
				log.Printf("Bot status: %s - %s", status.Status, status.Message)
				failureCount = 0
			case bot.StatusFailed, bot.StatusCrashed:
				log.Printf("Bot status: %s - %s (Error: %v)", status.Status, status.Message, status.Error)
				crashed = true
			case bot.StatusStopped:
				log.Printf("Bot status: %s - %s", status.Status, status.Message)
				crashed = true
			}
		case <-time.After(30 * time.Second):
			log.Println("Bot initialization timeout - considering as failure")
			crashed = true
		}

		if !crashed {
			log.Println("Bot is running, waiting for updates...")
			for status := range statusChan {
				log.Printf("Bot status update: %s - %s", status.Status, status.Message)
				if status.Error != nil {
					log.Printf("Error: %v", status.Error)
					crashed = true
					break
				}
			}
		}

		if crashed {
			failureCount++

			if failureCount >= MaxRetries {
				log.Printf("Bot failed %d times in a row. Maximum retries reached.", failureCount)
				panic("Bot failed to start after maximum retries. Application terminated.")
			}

			log.Printf("Bot will restart in %v... (Attempt %d/%d)", RetryDelay, failureCount+1, MaxRetries)
			time.Sleep(RetryDelay)
		}
	}
}
