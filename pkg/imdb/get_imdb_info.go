package imdb

import (
	"encoding/json"
	"log"
	"nightmare_navigator/internal/config"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	movieinfo "nightmare_navigator/pkg/movie_info"
)

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

func loadIMDbData(cfg config.Config) ([]movieinfo.MovieInfo, error) {
	file, err := os.Open(filepath.Join("data", cfg.IMDb.JSONFilename))
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

	var movieInfos []movieinfo.MovieInfo
	for _, yearDatas := range movieDatas {
		for _, movie := range yearDatas.Data {
			movieInfos = append(movieInfos, movieinfo.MovieInfo{
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

func GetIMDbInfosByYear(cfg config.Config, year string, getOMDbInfoByTitle func(string) *movieinfo.MovieInfo) []movieinfo.MovieInfo {
	imdbMovieInfos, err := loadIMDbData(cfg)
	if err != nil {
		return nil
	}

	var moviesByYear []movieinfo.MovieInfo

	for _, movie := range imdbMovieInfos {
		if movie.Year == year {
			omdbMovieDbInfo := getOMDbInfoByTitle(movie.Title)
			if omdbMovieDbInfo != nil {
				movie.Rated = omdbMovieDbInfo.Rated
				movie.Country = omdbMovieDbInfo.Country
			}

			moviesByYear = append(moviesByYear, movie)
		}
	}

	return moviesByYear
}

func GetIMDbInfosByDateAndGenre(cfg config.Config, count int, genres []string, date time.Time, getOMDbInfoByTitle func(string) *movieinfo.MovieInfo) *[]movieinfo.MovieInfo {
	imdbMovieInfos, err := loadIMDbData(cfg)
	if err != nil {
		return nil
	}

	var result []movieinfo.MovieInfo
	collectedCount := 0

	for i := 0; collectedCount < count; i++ {
		year := strconv.Itoa(date.Year() - i)
		filteredMovies := filterMovies(imdbMovieInfos, year, date, genres, count-collectedCount)

		for _, movie := range filteredMovies {
			omdbMovieDbInfo := getOMDbInfoByTitle(movie.Title)

			if omdbMovieDbInfo != nil {
				movie.Rated = omdbMovieDbInfo.Rated
				movie.Country = omdbMovieDbInfo.Country
			}
			result = append(result, movie)

			collectedCount += len(filteredMovies)
			if len(filteredMovies) == 0 {
				break
			}
		}
	}

	return &result
}

func sortMoviesByReleaseDate(movies []movieinfo.MovieInfo) {
	sort.Slice(movies, func(i, j int) bool {
		releaseDateI, errI := time.Parse("02.01.06", movies[i].ReleaseDate)
		releaseDateJ, errJ := time.Parse("02.01.06", movies[j].ReleaseDate)
		if errI != nil || errJ != nil {
			return movies[i].ReleaseDate < movies[j].ReleaseDate
		}
		return releaseDateI.After(releaseDateJ)
	})
}

func filterMovies(movies []movieinfo.MovieInfo, year string, date time.Time, genres []string, count int) []movieinfo.MovieInfo {
	var filteredMovies []movieinfo.MovieInfo
	sortMoviesByReleaseDate(movies)

	for _, movie := range movies {
		if movie.Year == year && movieMatchesGenres(movie, genres) && movieMatchDate(movie, date) && movie.Country != "India" {
			filteredMovies = append(filteredMovies, movie)
			if len(filteredMovies) >= count {
				break
			}
		}
	}

	return filteredMovies
}

func movieMatchesGenres(movie movieinfo.MovieInfo, genres []string) bool {
	movieGenres := splitGenres(movie.Genres)
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

func splitGenres(genres string) []string {
	parts := strings.Split(genres, " and ")
	result := []string{}

	for _, part := range parts {
		subParts := strings.Split(part, ", ")
		for _, subPart := range subParts {
			result = append(result, strings.TrimSpace(subPart))
		}
	}

	return result
}

func movieMatchDate(movie movieinfo.MovieInfo, date time.Time) bool {
	releaseDate, err := time.Parse("02.01.06", movie.ReleaseDate)

	if err != nil {
		log.Println("Error parse release date", err)
		return false
	}

	return !releaseDate.After(date)
}
