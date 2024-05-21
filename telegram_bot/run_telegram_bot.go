package telegram_bot

import (
	"log"
	"nightmare_navigator/manager"
	"nightmare_navigator/utils"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func RunTelegramBot() {
	// Initial imdb data saving
	manager.SaveLatestIMDbRatings()

	// Starts telegram bot
	bot, err := tgbotapi.NewBotAPI("6860257928:AAG8dygOS9j4rFl6x5oyrWxx8LIHbWsZATc")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Bot started as %s", bot.Self.UserName)

	// Function to calculate the duration until the next 03:00 AM
	durationUntilNextExecution := func() time.Duration {
		now := time.Now()
		nextExecution := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
		if now.After(nextExecution) {
			nextExecution = nextExecution.Add(24 * time.Hour)
		}
		return nextExecution.Sub(now)
	}

	// Function to be executed at 03:00 AM
	executeAt0300AM := func() {
		// Updates imdb_movies.json from the public IMDb dataset
		manager.SaveLatestIMDbRatings()

		// Checks if there are new movies this year and sends the movie information to all bot channel users
		newMovies := manager.GetLatestMovies()
		if newMovies != nil {
			for _, movie := range *newMovies {
				msg := tgbotapi.NewMessageToChannel("@nightmare_navigator", movie)
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	// Initial execution
	timer := time.NewTimer(durationUntilNextExecution())

	go func() {
		for {
			<-timer.C
			executeAt0300AM()
			timer.Reset(durationUntilNextExecution())
		}
	}()

	// Keeps the bot running by waiting for messages
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 6000

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if strings.Contains(update.Message.Text, "movie") {
			for _, movie := range *manager.GetFilteredMovies(utils.ExtractCount(update.Message.Text), utils.ExtractGenres(update.Message.Text), utils.ExtractDate(update.Message.Text)) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, movie)
				_, err = bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
