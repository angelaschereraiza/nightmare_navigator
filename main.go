package main

import (
	"log"
	"nightmare_navigator/api"
	"nightmare_navigator/imdb"
	"nightmare_navigator/utils"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// Starts telegram bot
	bot, err := tgbotapi.NewBotAPI("6860257928:AAG8dygOS9j4rFl6x5oyrWxx8LIHbWsZATc")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Bot started as %s", bot.Self.UserName)

	// Function to calculate the duration until the next 01:00
	durationUntilNextExecution := func() time.Duration {
		now := time.Now()
		nextExecution := time.Date(now.Year(), now.Month(), now.Day(), 01, 00, 0, 0, now.Location())
		if now.After(nextExecution) {
			nextExecution = nextExecution.Add(24 * time.Hour)
		}
		return nextExecution.Sub(now)
	}

	// Creates a timer that triggers at 01:00 AM
	timer := time.NewTimer(durationUntilNextExecution())

	go func() {
		<-timer.C
		imdb.SaveLatestIMDbRatings()
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
			for _, movie := range *api.GetFilteredLatestMovies(
				utils.ExtractCount(update.Message.Text),
				utils.ExtractGenres(update.Message.Text),
				utils.ExtractDate(update.Message.Text)) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, movie)
				_, err = bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
