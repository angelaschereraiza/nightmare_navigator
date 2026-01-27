package main

import (
	"log"

	"nightmare_navigator/internal/config"
	"nightmare_navigator/internal/imdb"
	"nightmare_navigator/internal/telegram_bot"
)

func main() {
	// Load application configuration from config file
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize IMDb manager and update latest movie ratings from IMDb
	log.Println("Loading latest movie data from IMDb...")
	imdbManager := imdb.NewIMDbManager(*cfg)
	imdbManager.SaveLatestIMDbRatings()
	log.Println("Movie data successfully loaded and saved")

	// Start the Nightmare Navigator Telegram bot
	log.Println("Starting Telegram bot...")
	telegram_bot.RunTelegramBot(cfg, imdbManager)
}

