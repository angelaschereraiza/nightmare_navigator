package manager

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type IMDbMovieInfo struct {
	Title       string
	IMDb        string
	IMDbVotes   string
	Year        string
	TitleId     string
	ReleaseDate string
}

type IMDbJsonData struct {
	Data []struct {
		Tconst        string `json:"tconst"`
		PrimaryTitle  string `json:"primaryTitle"`
		OriginalTitle string `json:"originalTitle"`
		Genres        string `json:"genres"`
		AverageRating string `json:"averageRating"`
		NumVotes      string `json:"numVotes"`
		ReleaseDate   string `json:"releaseDate"`
	} `json:"data"`
	StartYear string `json:"startYear"`
}

func loadIMDbData() *[]IMDbMovieInfo {
	file, err := os.Open(downloadDir + "/" + jsonFilename)
	if err != nil {
		log.Println("Error opening IMDb JSON file:", err)
		return nil
	}
	defer file.Close()

	var movieDatas []IMDbJsonData
	if err := json.NewDecoder(file).Decode(&movieDatas); err != nil {
		log.Println("Error decoding IMDb JSON file:", err)
		return nil
	}

	var imdbMovieInfos []IMDbMovieInfo
	for _, yearDatas := range movieDatas {
		for _, movie := range yearDatas.Data {
			imdbMovieInfo := IMDbMovieInfo{
				Title:       movie.PrimaryTitle,
				IMDb:        movie.AverageRating,
				IMDbVotes:   movie.NumVotes,
				Year:        yearDatas.StartYear,
				TitleId:     movie.Tconst,
				ReleaseDate: movie.ReleaseDate,
			}

			imdbMovieInfos = append(imdbMovieInfos, imdbMovieInfo)
		}
	}

	return &imdbMovieInfos
}

func getIMDbMovieInfosByYear(year string) []IMDbMovieInfo {
	imdbMovieInfos := loadIMDbData()

	var moviesByYear []IMDbMovieInfo

	for _, movie := range *imdbMovieInfos {
		if movie.Year == year {
			moviesByYear = append(moviesByYear, movie)
		}
	}

	return moviesByYear
}

func getIMDbMovieInfosByDateAndGenre(count *int, genres *[]string, date *time.Time) []IMDbMovieInfo {
	imdbMovieInfos := loadIMDbData()

	var moviesByYear []IMDbMovieInfo
	collectedCount := 0

	if count == nil {
		*count = 10
	}

	for _, movie := range *imdbMovieInfos {
		for collectedCount < *count {
			releaseDate, err := time.Parse("02.01.06", movie.ReleaseDate)
			if err != nil {
				log.Println(err)
			}

			if date.After(releaseDate) {
				moviesByYear = append(moviesByYear, movie)
				collectedCount++
			}
		}
	}

	return moviesByYear
}
