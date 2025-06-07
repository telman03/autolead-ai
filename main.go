package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/telman03/autolead-ai/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := bot.InitSupabase(); err != nil {
		log.Fatal("❌ Failed to initialize Supabase:", err)
	}


	bot.StartBot(os.Getenv("TELEGRAM_BOT_TOKEN"))
}
