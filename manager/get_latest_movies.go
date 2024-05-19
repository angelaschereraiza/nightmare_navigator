package manager

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

const dataDir = "data"

var alreadyReturnedMovies = make(map[string]bool)
var alreadyReturnedMoviesFile = filepath.Join(dataDir, "already_returned_movies.json")

func GetLatestMovies() *[]string {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	if err := loadAlreadyReturnedMovies(); err != nil {
		log.Println("Error loading already returned movies:", err)
	}

	currentYear := strconv.Itoa(time.Now().Year())
	movieInfos := getIMDbInfosByYear(currentYear)
	newMovies := filterAlreadyReturnedMovies(movieInfos)

	sortMoviesByReleaseDate(newMovies)

	for _, movie := range newMovies {
		alreadyReturnedMovies[movie.Title] = true
	}

	if err := saveReturnedMovies(); err != nil {
		log.Println("Error saving returned movies:", err)
	}

	movieStrings := buildMovieInfoStrings(newMovies)
	
	return movieStrings
}

func filterAlreadyReturnedMovies(imdbRatingsMovies []MovieInfo) []MovieInfo {
	filteredMovies := make([]MovieInfo, 0, len(imdbRatingsMovies))
	for _, movie := range imdbRatingsMovies {
		if !alreadyReturnedMovies[movie.Title] {
			filteredMovies = append(filteredMovies, movie)
		}
	}
	return filteredMovies
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

func saveReturnedMovies() error {
	movieTitles := make([]string, 0, len(alreadyReturnedMovies))
	for title := range alreadyReturnedMovies {
		movieTitles = append(movieTitles, title)
	}

	jsonData, err := json.MarshalIndent(movieTitles, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(alreadyReturnedMoviesFile, jsonData, 0644)
}

func loadAlreadyReturnedMovies() error {
	jsonData, err := os.ReadFile(alreadyReturnedMoviesFile)
	if err != nil {
		if os.IsNotExist(err) {
			if err := createEmptyJsonFile(alreadyReturnedMoviesFile); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	var movieTitles []string
	if err := json.Unmarshal(jsonData, &movieTitles); err != nil {
		return err
	}

	for _, title := range movieTitles {
		alreadyReturnedMovies[title] = true
	}
	return nil
}

func createEmptyJsonFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	emptyJson := "[]"
	if _, err := file.WriteString(emptyJson); err != nil {
		return err
	}
	return nil
}