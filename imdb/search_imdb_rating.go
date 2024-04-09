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

func GetIMDbInfoByTitle(title string, year string) *IMDbMovieInfo {
	file, err := os.Open(downloadDir + "/" + jsonFilename)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()

	var movies []IMDbMovieInfo
	if err := json.NewDecoder(file).Decode(&movies); err != nil {
		log.Println(err)
		return nil
	}

	for _, movie := range movies {
		if year == movie.Year && (movie.Title == title || movie.OriginalTitle == title) {
			return &movie
		}
	}

	return nil
}
