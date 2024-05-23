package telegram_bot

import (
	"log"
	"strings"
	"time"

	"nightmare_navigator/internal/config"
	buildstring "nightmare_navigator/internal/movie_info/build_movie_info_string"
	filteredinfo "nightmare_navigator/internal/movie_info/get_filtered_movie_info"
	imdbinfo "nightmare_navigator/internal/movie_info/get_imdb_info"
	latestinfo "nightmare_navigator/internal/movie_info/get_latest_movie_info"
	imdb "nightmare_navigator/internal/movie_info/save_imdb_info"
	"nightmare_navigator/internal/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func RunTelegramBot() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	imdbManager := imdb.NewSaveIMDbInfoManager(*cfg)
	// imdbManager.SaveLatestIMDbRatings()

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBot.Token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Bot started as %s", bot.Self.UserName)

	timer := time.NewTimer(durationUntilNextExecution())

	go func() {
		for {
			<-timer.C
			executeAt0300AM(bot, cfg.TelegramBot.ChannelName, *imdbManager, *cfg)
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
			movieInfos := filteredinfo.GetFilteredMovieInfos(
				util.ExtractCount(update.Message.Text),
				util.ExtractGenres(update.Message.Text),
				util.ExtractDate(update.Message.Text),
				imdbinfo.GetIMDbInfosByDateAndGenre,
				buildstring.BuildMovieInfoStrings,
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

func executeAt0300AM(bot *tgbotapi.BotAPI, channelName string, imdbManager imdb.SaveIMDbInfoManager, cfg config.Config) {
	imdbManager.SaveLatestIMDbRatings()
	latestMoviesManager := latestinfo.NewLatestMoviesManager(cfg)
	newMovies := latestMoviesManager.GetLatestMovieInfos(imdbinfo.GetIMDbInfosByYear, buildstring.BuildMovieInfoStrings)
	if newMovies != nil {
		for _, movie := range *newMovies {
			msg := tgbotapi.NewMessageToChannel(channelName, movie)
			_, err := bot.Send(msg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}