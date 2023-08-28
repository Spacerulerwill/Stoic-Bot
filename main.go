package main

import (
	"log"
	"os"
	bot "stoicbot/src"

	"github.com/joho/godotenv"
)

func main() {
	// Load bot token
	godotenv.Load(".env")

	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("Environment variable BOT_TOKEN not found!")
		return
	}

	bot.Run(botToken, "1145325203291910235")
}
