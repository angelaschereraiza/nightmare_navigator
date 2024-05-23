package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TelegramBot TelegramBot `yaml:"telegram_bot"`
	General     General     `yaml:"latest_movies"`
	IMDb        IMDb        `yaml:"imdb"`
	TMDb        TMDb        `yaml:"tmdb"`
	OMDb        OMDb        `yaml:"omdb"`
}

type TelegramBot struct {
	Token       string `yaml:"token"`
	ChannelName string `yaml:"channel_name"`
}

type General struct {
	DataDir                   string `yaml:"data_dir"`
	AlreadyReturnedMoviesJSON string `yaml:"already_returned_movies_json"`
}

type IMDb struct {
	BasicsFilename  string  `yaml:"basics_filename"`
	JSONFilename    string  `yaml:"json_filename"`
	MinRating       float64 `yaml:"min_rating"`
	MinVotes        int     `yaml:"min_votes"`
	RatingsFilename string  `yaml:"ratings_filename"`
	DownloadDir     string  `yaml:"download_dir"`
	IMDbBaseUrl     string  `yaml:"imdb_base_url"`
}

type TMDb struct {
	ApiKey string `yaml:"api_key"`
	ApiURL string `yaml:"api_url"`
}

type OMDb struct {
	ApiKey string `yaml:"api_key"`
	ApiURL string `yaml:"api_url"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
