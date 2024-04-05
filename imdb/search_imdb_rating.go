package imdb

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type IMDbMovie struct {
	Title         string `json:"primaryTitle"`
	OriginalTitle string `json:"originalTitle"`
	TitleId       string `json:"tconst"`
	IMDb          string `json:"averageRating"`
	IMDbVotes     string `json:"numVotes"`
}

func GetIMDbInfoByTitle(title string) *IMDbMovie {
	ratingsFilePath := "imdb/imdb_ratings.json"
	file, err := os.Open(ratingsFilePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var movie IMDbMovie
		err := json.Unmarshal(scanner.Bytes(), &movie)
		if err != nil {
			log.Println(err)
			continue
		}
		
		if movie.Title == title || movie.OriginalTitle == title {
			return &movie
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	return nil
}
