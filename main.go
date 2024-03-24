package main

import (
	"log"
	"nightmare_navigator/themoviedb"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

		if strings.Contains(update.Message.Text, "movies") {
			count := extractCount(update.Message.Text)
			date := extractDate(update.Message.Text)
			for _, movie := range *themoviedb.GetFilteredLatestMovies(count, date) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, movie)
				_, err = bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func extractCount(text string) int {
	count := 20
	numStr := ""
	for _, char := range text {
		if unicode.IsDigit(char) {
			numStr += string(char)
		} else if numStr != "" {
			break
		}
	}

	if numStr != "" {
		count, _ = strconv.Atoi(numStr)
	}

	return count
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
