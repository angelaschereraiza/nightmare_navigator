package telegram_bot

import (
	"log"
	"strings"
	"time"

	"nightmare_navigator/internal/config"
	movieinfo "nightmare_navigator/internal/movie_info"
	"nightmare_navigator/internal/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func RunTelegramBot(cfg *config.Config, imdbManager *movieinfo.SaveIMDbInfoManager) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBot.Token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Bot started as %s", bot.Self.UserName)

	timer := time.NewTimer(durationUntilNextExecution())

	go func() {
		for {
			<-timer.C
			executeAt0300AM(bot, cfg.TelegramBot.GroupId, *imdbManager, *cfg)
			timer.Reset(durationUntilNextExecution())
		}
	}()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if strings.Contains(update.Message.Text, "movie") {
			movieInfos := movieinfo.GetFilteredMovieInfos(
				util.ExtractCount(update.Message.Text),
				util.ExtractGenres(update.Message.Text),
				util.ExtractDate(update.Message.Text),
				movieinfo.GetIMDbInfosByDateAndGenre,
				movieinfo.BuildMovieInfoStrings,
				*cfg,
			)

			for _, movieInfo := range *movieInfos {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, movieInfo)
				_, err = bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}

		if strings.Contains(update.Message.Text, "/help") {
			helpInfo := "Usage:\n- Type 'movies' for the 10 latest horror movies \n- Specify a number for more \ne.g. '15 movies' \n- Add a date e.g. 'movies 14.06.19' \n- Specify genres: Sci-Fi, Fantasy, Thriller, Animation or Mystery e.g. 'sci-fi movies' \n- Combine options e.g. '20 sci-fi thriller movies 14.06.19'"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpInfo)
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func durationUntilNextExecution() time.Duration {
	now := time.Now()
	nextExecution := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
	if now.After(nextExecution) {
		nextExecution = nextExecution.Add(24 * time.Hour)
	}
	return nextExecution.Sub(now)
}

func executeAt0300AM(bot *tgbotapi.BotAPI, groupId int64, imdbManager movieinfo.SaveIMDbInfoManager, cfg config.Config) {
	imdbManager.SaveLatestIMDbRatings()
	latestMoviesManager := movieinfo.NewLatestMoviesManager(cfg)
	newMovies := latestMoviesManager.GetLatestMovieInfos(movieinfo.GetIMDbInfosByYear, movieinfo.BuildMovieInfoStrings)
	if newMovies != nil {
		for _, movie := range *newMovies {
			msg := tgbotapi.NewMessage(groupId, movie)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
