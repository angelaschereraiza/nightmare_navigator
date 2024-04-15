package imdb

import (
	"encoding/json"
	"log"
	"os"
)

type IMDbMovieInfo struct {
	Title         string `json:"primaryTitle"`
	OriginalTitle string `json:"originalTitle"`
	TitleId       string `json:"tconst"`
	IMDb          string `json:"averageRating"`
	IMDbVotes     string `json:"numVotes"`
	Year          string `json:"startYear"`
}

var imdbData map[string]map[string]*IMDbMovieInfo

func init() {
	imdbData = make(map[string]map[string]*IMDbMovieInfo)
	loadIMDbData()
}

func loadIMDbData() {
	file, err := os.Open(downloadDir + "/" + jsonFilename)
	if err != nil {
		log.Println("Error opening IMDb JSON file:", err)
		return
	}
	defer file.Close()

	var movies []IMDbMovieInfo
	if err := json.NewDecoder(file).Decode(&movies); err != nil {
		log.Println("Error decoding IMDb JSON file:", err)
		return
	}

	for _, movie := range movies {
		if _, ok := imdbData[movie.Year]; !ok {
			imdbData[movie.Year] = make(map[string]*IMDbMovieInfo)
		}
		imdbData[movie.Year][movie.Title] = &movie
		imdbData[movie.Year][movie.OriginalTitle] = &movie
	}
}

func GetIMDbInfoByTitle(title string, year string) *IMDbMovieInfo {
	if yearData, ok := imdbData[year]; ok {
		if movie, ok := yearData[title]; ok {
			return movie
		}
	}
	return nil
}
