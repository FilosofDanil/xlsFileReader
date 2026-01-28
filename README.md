# XLS Reader Telegram Bot

## English

### What is this?
This project is a **Telegram bot** that accepts **Excel files (`.xlsx` / `.xls`)**, extracts contract numbers from rows that contain **`ББС ІНШУРАНС`**, and generates a **SQL script (`script.txt`)** with a `SELECT ... FROM (VALUES ...)` block, then sends it back to the user.

### Tech stack / libraries
- **Go**
- **Telegram Bot API**: `github.com/go-telegram-bot-api/telegram-bot-api/v5`
- **XLSX reader**: `github.com/xuri/excelize/v2`
- **XLS reader**: `github.com/extrame/xls`

### How to use
- Open the bot in Telegram (optional: `/start`).
- Send an **`.xlsx` or `.xls`** file.
- The bot replies with **`script.txt`** (SQL query containing your extracted contracts).

### Examples (screenshots)
Add your screenshots here and keep these paths:
- `docs/screenshots/telegram-chat.png`
- `docs/screenshots/script-result.png`

### Run locally
1. Put your Telegram token in `cmd/main.go` (`TelegramBotToken` constant).
2. Install dependencies:

```bash
go mod tidy
```

3. Run:

```bash
go run ./cmd/main.go
```

### Run with Docker
1. Put your Telegram token in `cmd/main.go`.
2. Start:

```bash
docker-compose up -d --build
```

3. Logs:

```bash
docker-compose logs -f
```


## Deutsch

### Was ist das?
Dieses Projekt ist ein **Telegram-Bot**, der **Excel-Dateien (`.xlsx` / `.xls`)** entgegennimmt, Vertragsnummern aus Zeilen mit **`ББС ІНШУРАНС`** extrahiert und ein **SQL-Skript (`script.txt`)** generiert und zurücksendet.

### Tech-Stack / Bibliotheken
- **Go**
- **Telegram Bot API**: `github.com/go-telegram-bot-api/telegram-bot-api/v5`
- **XLSX**: `github.com/xuri/excelize/v2`
- **XLS**: `github.com/extrame/xls`

### Nutzung
- Bot in Telegram öffnen (optional: `/start`).
- Eine **`.xlsx`- oder `.xls`-Datei** senden.
- Der Bot antwortet mit **`script.txt`** (SQL-Abfrage mit den extrahierten Verträgen).

### Beispiele (Screenshots)
Empfohlene Pfade:
- `docs/screenshots/telegram-chat.png`
- `docs/screenshots/script-result.png`

### Lokal starten
1. Telegram-Token in `cmd/main.go` setzen (`TelegramBotToken`).
2. Abhängigkeiten:

```bash
go mod tidy
```

3. Starten:

```bash
go run ./cmd/main.go
```

### Mit Docker starten
1. Telegram-Token in `cmd/main.go` setzen.
2. Start:

```bash
docker-compose up -d --build
```


## Українська

### Що це?
Це **Telegram-бот**, який приймає **Excel-файли (`.xlsx` / `.xls`)**, знаходить рядки з **`ББС ІНШУРАНС`**, витягує дані договорів і формує **SQL-скрипт (`script.txt`)**, після чого надсилає його користувачу.

### Стек / бібліотеки
- **Go**
- **Telegram Bot API**: `github.com/go-telegram-bot-api/telegram-bot-api/v5`
- **XLSX**: `github.com/xuri/excelize/v2`
- **XLS**: `github.com/extrame/xls`

### Як користуватись
- Відкрийте бота в Telegram (опційно: `/start`).
- Надішліть файл **`.xlsx` або `.xls`**.
- Бот поверне **`script.txt`** (SQL-запит з витягнутими договорами).

### Приклади (скріншоти)
Додайте скріншоти у:
- `docs/screenshots/telegram-chat.png`
- `docs/screenshots/script-result.png`

### Запуск локально
1. Вкажіть Telegram token у `cmd/main.go` (константа `TelegramBotToken`).
2. Встановіть залежності:

```bash
go mod tidy
```

3. Запустіть:

```bash
go run ./cmd/main.go
```

### Запуск через Docker
1. Вкажіть Telegram token у `cmd/main.go`.
2. Запуск:

```bash
docker-compose up -d --build
```
# Telegram Bot Application

A Go-based Telegram bot that handles various types of messages including text, images, audio, files, and more.

## Project Structure

```
.
├── cmd/
│   ├── main.go              # Main entry point (production)
│   └── prototype/
│       └── main.go          # Test entry point with test data
├── handler/
│   └── handler.go           # Message handlers for all message types
├── env/
│   └── env.go              # Environment configuration management
├── go.mod                   # Go module dependencies
├── .env.example            # Example environment configuration
└── README.md               # This file
```

## Features

### Message Handlers
- **Text Messages**: Processes and analyzes text messages
- **Commands**: Handles `/start`, `/help`, `/status` commands
- **Images**: Receives and processes photos with metadata
- **Audio**: Handles audio files with metadata
- **Documents**: Processes any file/document uploads
- **Video**: Handles video files
- **Voice**: Processes voice messages
- **Callback Queries**: Handles inline keyboard button clicks

## Setup

### Prerequisites
- Go 1.24.3 or higher
- A Telegram Bot Token (get one from [@BotFather](https://t.me/BotFather))

### Installation

1. **Clone or navigate to the project directory**

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   ```bash
   # Copy the example env file
   copy .env.example .env
   
   # Edit .env and add your bot token
   # BOT_TOKEN=your_actual_bot_token_here
   ```

## Running the Application

### Production Mode (using .env file)

1. Make sure your `.env` file is configured with your bot token
2. Run the main application:
   ```bash
   go run cmd/main.go
   ```

### Prototype/Test Mode (using hardcoded test data)

For quick testing without setting up environment variables:

1. Open `cmd/prototype/main.go`
2. Replace `YOUR_TEST_BOT_TOKEN_HERE` with your test bot token
3. Run the prototype:
   ```bash
   go run cmd/prototype/main.go
   ```

Alternatively, use an environment variable:
```bash
# PowerShell
$env:TEST_BOT_TOKEN="your_test_token_here"
go run cmd/prototype/main.go

# Command Prompt
set TEST_BOT_TOKEN=your_test_token_here
go run cmd/prototype/main.go
```

## Usage

Once the bot is running:

1. Open Telegram and search for your bot by username
2. Start a conversation with `/start`
3. Try different message types:
   - Send text messages
   - Upload images
   - Send audio files
   - Upload documents
   - Send videos
   - Record voice messages
   - Click inline buttons

## Configuration

Edit the `.env` file to configure:

- `BOT_TOKEN`: Your Telegram bot token (required)
- `PORT`: Server port (default: 8080)
- `DEBUG`: Enable debug logging (default: false)

## Dependencies

- [telebot.v3](https://github.com/tucnak/telebot) - Telegram Bot API framework
- [godotenv](https://github.com/joho/godotenv) - Environment variable loader

## Development

### Adding New Handlers

1. Open `handler/handler.go`
2. Create a new handler function following the pattern:
   ```go
   func (h *Handler) HandleNewType(c tele.Context) error {
       // Your handler logic here
       return c.Send("Response message")
   }
   ```
3. Register it in `RegisterHandlers()`:
   ```go
   h.bot.Handle(tele.OnNewType, h.HandleNewType)
   ```

## License

This project is open source and available for educational purposes.

