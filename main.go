package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/telman03/autolead-ai/bot"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	// Get token from .env
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN not found in .env")
	}

	// Start the Telegram bot
	bot.StartBot(token)
}
