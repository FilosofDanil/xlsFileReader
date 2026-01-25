package main

import (
	"example/hello/bot"
	"example/hello/handler"
	"log"
	"time"
)

const (
	TelegramBotToken = "8390270520:AAFe6be2LWMg60yPyxsJEwwZFBH5rNlbDA"
	MaxRetries       = 3
	RetryDelay       = 15 * time.Second
)

func main() {
	log.Println("Application starting...")

	statusChan := make(chan bot.BotStatus, 10)
	failureCount := 0

	for {
		log.Printf("Launch attempt %d/%d", failureCount+1, MaxRetries)

		botService := bot.NewService(TelegramBotToken, statusChan)
		messageHandler := handler.NewHandler()

		go func() {
			err := botService.Start(messageHandler.HandleUpdate)
			if err != nil {
				log.Printf("Bot crashed with error: %v", err)
				statusChan <- bot.BotStatus{
					Status:  "crashed",
					Error:   err,
					Message: "Bot encountered an error",
				}
			}
		}()

		crashed := false

		select {
		case status := <-statusChan:
			switch status.Status {
			case "started":
				log.Printf("Bot status: %s - %s", status.Status, status.Message)
				failureCount = 0
			case "failed", "crashed":
				log.Printf("Bot status: %s - %s (Error: %v)", status.Status, status.Message, status.Error)
				crashed = true
			case "stopped":
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
