package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	downloadDir  = "data"
	jsonFilename = "already_returned_movies.json"
)

var alreadyReturnedMovies = make(map[string]bool)
var alreadyReturnedMoviesFile = filepath.Join(downloadDir, jsonFilename)

func SearchForNewMovies() *[]string {
	// Create the directory if it does not exist
	err := os.MkdirAll(downloadDir, 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	// Tries to read json file
	jsonData, err := os.ReadFile(alreadyReturnedMoviesFile)
	if err != nil {
		// If there was a problem, maybe the json file does not exist, so a new json file is created
		jsonData, err := os.Create(downloadDir + "/" + jsonFilename)
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

	genreMap := getGenres()
	if genreMap == nil {
		return nil
	}

	var results []string

	startDateString := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	endDateString := time.Now().Format("2006-01-02")

	page := 1
	hasNewMovies := true

	for hasNewMovies {
		params := url.Values{}
		params.Set("api_key", apiKey)
		params.Set("without_genres", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(excludeGenres)), ","), "[]"))
		params.Set("with_runtime.gte", "85")
		params.Set("sort_by", "primary_release_date.desc")
		params.Set("primary_release_date.gte", startDateString)
		params.Set("primary_release_date.lte", endDateString)
		params.Set("page", fmt.Sprintf("%d", page))

		searchURL := fmt.Sprintf("%s/discover/movie?%s", baseURL, params.Encode())

		res, err := http.Get(searchURL)
		if err != nil {
			log.Println("Error making request:", err)
			return nil
		}
		defer res.Body.Close()

		var movieResponse MovieResponse
		err = json.NewDecoder(res.Body).Decode(&movieResponse)
		if err != nil {
			log.Println("Error decoding JSON response:", err)
			return nil
		}

		if len(movieResponse.Results) == 0 {
			break
		}

		for _, movie := range movieResponse.Results {
			if !alreadyReturnedMovies[movie.Title] {
				result := getAdditionalMovieInfo(movie, *genreMap)
				if result != "" {
					results = append(results, result)
					alreadyReturnedMovies[movie.Title] = true
				}
			}
		}

		page++
	}

	// Save alreadyReturnedMovies list to JSON file
	err = saveReturnedMovies()
	if err != nil {
		log.Println("Error saving returned movies:", err)
	}

	results = reverseStringArray(results)

	return &results
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

func reverseStringArray(arr []string) []string {
	// Erstelle ein neues String-Array mit der gleichen Länge wie das ursprüngliche Array
	reversed := make([]string, len(arr))

	// Iteriere rückwärts durch das ursprüngliche Array und kopiere die Elemente in das neue Array
	for i := 0; i < len(arr); i++ {
		reversed[i] = arr[len(arr)-1-i]
	}

	return reversed
}
