package imdb

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type IMDbMovieInfo struct {
	Title         string `json:"primaryTitle"`
	OriginalTitle string `json:"originalTitle"`
	TitleId       string `json:"tconst"`
	IMDb          string `json:"averageRating"`
	IMDbVotes     string `json:"numVotes"`
}

func GetIMDbInfoByTitle(title string) *IMDbMovieInfo {
	file, err := os.Open(downloadDir + "/" + jsonFilename)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()

	content, err := io.ReadAll(io.Reader(file))
	if err != nil {
		log.Println(err)
		return nil
	}

	var movies []IMDbMovieInfo
	err = json.Unmarshal(content, &movies)
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, movie := range movies {
		if movie.Title == title || movie.OriginalTitle == title {
			return &movie
		}
	}

	return nil
}
