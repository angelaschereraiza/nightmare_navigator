package main

import (
	"log"
	"nightmare_navigator/themoviedb"
	"regexp"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6860257928:AAG8dygOS9j4rFl6x5oyrWxx8LIHbWsZATc")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Bot started as %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if strings.HasPrefix(update.Message.Text, "movies") {
			date := extractDate(update.Message.Text)
			for _, movie := range *themoviedb.GetFilteredLatestMovies(date) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, movie)
				_, err = bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func extractDate(text string) time.Time {
	now := time.Now()

	if text == "" {
		return now
	}

	re := regexp.MustCompile(`\d{2}\.\d{2}\.\d{2}`)
	dateStr := re.FindString(text)
	if dateStr == "" {
		return now
	}

	date, err := time.Parse("02.01.06", dateStr)
	if err != nil {
		log.Println(err)
		return now
	}

	return date
}