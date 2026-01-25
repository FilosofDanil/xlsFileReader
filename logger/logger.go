package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	LogFilePath        = "logs/app.log"
	DefaultLogFilePath = "log/log.txt"
)

func Init() {
	logPath := LogFilePath

	if err := ensureLogFile(logPath); err != nil {
		log.Printf("Failed to create log file at %s: %v. Using default path.", logPath, err)
		logPath = DefaultLogFilePath

		if err := ensureLogFile(logPath); err != nil {
			log.Printf("Failed to create default log file at %s: %v. Using stdout only.", logPath, err)
			return
		}
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v. Using stdout only.", err)
		return
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("Logger initialized. Writing to: %s", logPath)
}

func ensureLogFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	file.Close()

	return nil
}
