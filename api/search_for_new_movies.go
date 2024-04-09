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

var alreadyReturnedMovies = make(map[string]bool)
var alreadyReturnedMoviesFile = filepath.Join(".", "already_returned_movies.json")

func SearchForNewMovies() *[]string {
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
	err := saveReturnedMovies()
	if err != nil {
		log.Println("Error saving returned movies:", err)
	}

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
func loadAlreadyReturnedMovies() error {
	jsonData, err := os.ReadFile(alreadyReturnedMoviesFile)
	if err != nil {
		return err
	}
	var movieTitles []string
	err = json.Unmarshal(jsonData, &movieTitles)
	if err != nil {
		return err
	}
	for _, title := range movieTitles {
		alreadyReturnedMovies[title] = true
	}
	return nil
}

func init() {
	// Get the absolute path of the directory containing the executable
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)

	// Set the filepath for storing the returnedMovies list
	alreadyReturnedMoviesFile = filepath.Join(exPath, "already_returned_movies.json")

	// Load alreadyReturnedMovies list from JSON file
	err = loadAlreadyReturnedMovies()
	if err != nil {
		log.Println("Error loading already returned movies:", err)
	}
}
