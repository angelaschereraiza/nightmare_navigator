package api

import (
	"encoding/json"
	"log"
	"nightmare_navigator/imdb"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var alreadyReturnedMovies = make(map[string]bool)
var alreadyReturnedMoviesFile = filepath.Join("data", "already_returned_movies.json")

func SearchForNewMovies() *[]string {
	// Create the directory if it does not exist
	err := os.MkdirAll("data", 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	// Tries to read alreadyReturnedMovies json file
	jsonData, err := os.ReadFile(alreadyReturnedMoviesFile)
	if err != nil {
		// If there was a problem, maybe the json file does not exist, so a new json file is created
		jsonData, err := os.Create(alreadyReturnedMoviesFile)
		if err != nil {
			log.Fatal("Error creating json file:", err)
		}
		defer jsonData.Close()
	} else {
		// Load alreadyReturnedMovies list from JSON file
		err = loadAlreadyReturnedMovies(jsonData)
		if err != nil {
			log.Println("Error loading already returned movies:", err)
		}
	}

	// Gets all movies from imdb ratings json where the year is the current year
	imdbRatingMovies := imdb.GetIMDbMovieInfosByYear(strconv.Itoa(time.Now().Year()))

	// Check if imdbRatinsMovie is contained in alreadyReturnedMovies and removes it from the list if true
	newMovies := filterAlreadyReturnedMovies(imdbRatingMovies)

	// Adds the title of the new movies into alreadyReturnedMovies list
	for _, movie := range newMovies {
		alreadyReturnedMovies[movie.Title] = true
	}

	// Save alreadyReturnedMovies list to JSON file
	err = saveReturnedMovies()
	if err != nil {
		log.Println("Error saving returned movies:", err)
	}

	// Gets additional movie information
	newMoviesWithAllInfo := GetAdditionalMovieInfo(newMovies)

	return newMoviesWithAllInfo
}

// Function to filter out movies that are already returned
func filterAlreadyReturnedMovies(imdbRatingsMovies []imdb.IMDbMovieInfo) []imdb.IMDbMovieInfo {
	var filteredMovies []imdb.IMDbMovieInfo

	for _, movie := range imdbRatingsMovies {
		if !alreadyReturnedMovies[movie.Title] {
			filteredMovies = append(filteredMovies, movie)
		}
	}

	return filteredMovies
}

// Function to save the alreadyReturnedMovies list to a JSON file
func saveReturnedMovies() error {
	var movieTitles []string
	for title := range alreadyReturnedMovies {
		movieTitles = append(movieTitles, title)
	}

	jsonData, err := json.MarshalIndent(movieTitles, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(alreadyReturnedMoviesFile, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Function to load the alreadyReturnedMovies list from a JSON file
func loadAlreadyReturnedMovies(jsonData []byte) error {
	var movieTitles []string
	err := json.Unmarshal(jsonData, &movieTitles)
	if err != nil {
		return err
	}
	for _, title := range movieTitles {
		alreadyReturnedMovies[title] = true
	}
	return nil
}
