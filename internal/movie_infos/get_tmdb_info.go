package movie_infos

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var BaseURL = "https://api.themoviedb.org/3"

const apiKey = "6882b8441ce200fda300c1e46eeb3e64"

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GenresResponse struct {
	Genres []Genre `json:"genres"`
}

type TMDbInfo struct {
	ID          int    `json:"id"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
	Runtime     int
	Title       string `json:"title"`
}

type TMDbInfos struct {
	Results []TMDbInfo `json:"results"`
}

type MovieDetailsResponse struct {
	Runtime int `json:"runtime"`
}

func GetTMDbInfoByTitle(title string, year string) *TMDbInfo {
	// Fetch movie details from The Movie Database (TMDb) API
	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("query", title)
	searchURL := fmt.Sprintf("%s/search/movie?%s", BaseURL, params.Encode())

	tmdbInfos, err := fetchMovieSearchResults(searchURL)
	if err != nil || len(tmdbInfos.Results) == 0 {
		return nil
	}

	movieInfo := TMDbInfo{}

	for _, movie := range tmdbInfos.Results {
		releaseDate := setReleaseDate(movie.ReleaseDate, year)
		if releaseDate == nil {
			continue
		}

		movieInfo = movie
		movieInfo.ReleaseDate = *releaseDate

		// Fetch additional movie details (runtime) from TMDb API
		movieDetailsURL := fmt.Sprintf("%s/movie/%d?api_key=%s", BaseURL, movie.ID, apiKey)
		movieDetails, err := fetchMovieDetails(movieDetailsURL)
		if err != nil {
			return nil
		}

		movieInfo.Runtime = movieDetails.Runtime
	}

	if movieInfo.ReleaseDate == "" {
		return nil
	}

	return &movieInfo
}

func fetchMovieSearchResults(url string) (*TMDbInfos, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error making request:", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("Error: received non-200 response code: %d\n", res.StatusCode)
		return nil, fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	var tmdbInfos TMDbInfos
	if err := json.NewDecoder(res.Body).Decode(&tmdbInfos); err != nil {
		log.Println("Error decoding JSON response:", err)
		return nil, err
	}

	return &tmdbInfos, nil
}

func fetchMovieDetails(url string) (*MovieDetailsResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching movie details:", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("Error: received non-200 response code: %d\n", res.StatusCode)
		return nil, fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	var movieDetails MovieDetailsResponse
	if err := json.NewDecoder(res.Body).Decode(&movieDetails); err != nil {
		log.Println("Error decoding JSON response for movie details:", err)
		return nil, err
	}

	return &movieDetails, nil
}

func setReleaseDate(releaseDate string, year string) *string {
	if releaseDate == "" {
		return nil
	}

	parsedDate, err := time.Parse("2006-01-02", releaseDate)
	if err != nil {
		log.Println("Error parse release date", err)
		return nil
	}

	if strconv.Itoa(parsedDate.Year()) != year {
		log.Printf("release date is not in the same year. Release date: %s, year: %s", parsedDate.Format("02.01.06"), year)
		return nil
	}

	if parsedDate.After(time.Now()) {
		log.Println("release date is in the future")
		return nil
	}

	releaseDate = parsedDate.Format("02.01.06")

	return &releaseDate
}
