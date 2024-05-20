package manager

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MovieInfo struct {
	Country       string
	IMDb          string
	IMDbVotes     string
	Genres        string
	OriginalTitle string
	Description   string
	Rated         string
	ReleaseDate   string
	Runtime       int
	Title         string
	TitleId       string
	Year          string
}

type IMDbJsonData struct {
	Data []struct {
		AverageRating string `json:"averageRating"`
		Description   string `json:"description"`
		Genres        string `json:"genres"`
		NumVotes      string `json:"numVotes"`
		OriginalTitle string `json:"originalTitle"`
		PrimaryTitle  string `json:"primaryTitle"`
		ReleaseDate   string `json:"releaseDate"`
		Runtime       int    `json:"runtime"`
		Tconst        string `json:"tconst"`
	} `json:"data"`
	StartYear string `json:"startYear"`
}

func loadIMDbData() ([]MovieInfo, error) {
	file, err := os.Open(filepath.Join(downloadDir, jsonFilename))
	if err != nil {
		log.Println("Error opening IMDb JSON file:", err)
		return nil, err
	}
	defer file.Close()

	var movieDatas []IMDbJsonData
	if err := json.NewDecoder(file).Decode(&movieDatas); err != nil {
		log.Println("Error decoding IMDb JSON file:", err)
		return nil, err
	}

	var movieInfos []MovieInfo
	for _, yearDatas := range movieDatas {
		for _, movie := range yearDatas.Data {
			movieInfos = append(movieInfos, MovieInfo{
				Description:   movie.Description,
				IMDb:          movie.AverageRating,
				IMDbVotes:     movie.NumVotes,
				Genres:        movie.Genres,
				OriginalTitle: movie.OriginalTitle,
				ReleaseDate:   movie.ReleaseDate,
				Runtime:       movie.Runtime,
				Title:         movie.PrimaryTitle,
				TitleId:       movie.Tconst,
				Year:          yearDatas.StartYear,
			})
		}
	}

	return movieInfos, nil
}

func getIMDbInfosByYear(year string) []MovieInfo {
	imdbMovieInfos, err := loadIMDbData()
	if err != nil {
		return nil
	}

	var moviesByYear []MovieInfo
	for _, movie := range imdbMovieInfos {
		omdbMovieDbInfo := getOMDbInfoByTitle(movie.Title)
		if omdbMovieDbInfo != nil {
			movie.Rated = omdbMovieDbInfo.Rated
			movie.Country = omdbMovieDbInfo.Country
		}

		if movie.Year == year {
			moviesByYear = append(moviesByYear, movie)
		}
	}

	return moviesByYear
}

func getIMDbInfosByDateAndGenre(count int, genres []string, date time.Time) *[]MovieInfo {
	imdbMovieInfos, err := loadIMDbData()
	if err != nil {
		return nil
	}

	var result []MovieInfo
	collectedCount := 0

	for i := 0; collectedCount < count; i++ {
		year := strconv.Itoa(date.Year() - i)
		yearMovies := filterMoviesByYear(imdbMovieInfos, year, genres, count-collectedCount)

		for _, movie := range yearMovies {
			omdbMovieDbInfo := getOMDbInfoByTitle(movie.Title)
			if omdbMovieDbInfo != nil {
				movie.Rated = omdbMovieDbInfo.Rated
				movie.Country = omdbMovieDbInfo.Country
			}
			result = append(result, movie)
		}

		collectedCount += len(yearMovies)
		if len(yearMovies) == 0 {
			break
		}
	}

	return &result
}

func sortMoviesByReleaseDate(movies []MovieInfo) {
	sort.Slice(movies, func(i, j int) bool {
		releaseDateI, errI := time.Parse("02.01.06", movies[i].ReleaseDate)
		releaseDateJ, errJ := time.Parse("02.01.06", movies[j].ReleaseDate)
		if errI != nil || errJ != nil {
			return movies[i].ReleaseDate < movies[j].ReleaseDate
		}
		return releaseDateI.Before(releaseDateJ)
	})
}

func filterMoviesByYear(movies []MovieInfo, year string, genres []string, count int) []MovieInfo {
	var filteredMovies []MovieInfo
	sortMoviesByReleaseDate(movies)

	for _, movie := range movies {
		if movie.Year == year && movieMatchesGenres(movie, genres) {
			filteredMovies = append(filteredMovies, movie)
			if len(filteredMovies) >= count {
				break
			}
		}
	}

	return filteredMovies
}

func movieMatchesGenres(movie MovieInfo, genres []string) bool {
	movieGenres := strings.Split(movie.Genres, ",")
	for _, genre := range genres {
		found := false
		for _, movieGenre := range movieGenres {
			if strings.TrimSpace(movieGenre) == genre {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
