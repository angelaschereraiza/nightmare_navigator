package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://api.themoviedb.org/3"
	apiKey  = "6882b8441ce200fda300c1e46eeb3e64"
)

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GenresResponse struct {
	Genres []Genre `json:"genres"`
}

type TheMovieDbInfo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Overview    string `json:"overview"`
}

type TheMovieDbInfos struct {
	Results []TheMovieDbInfo `json:"results"`
}

type MovieDetailsResponse struct {
	Runtime int `json:"runtime"`
}

func getTheMovieDbInfoByTitle(title string) *MovieInfo {
	// Fetch movie details from The Movie Database (TMDb) API
	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("query", title)
	searchURL := fmt.Sprintf("%s/search/movie?%s", baseURL, params.Encode())

	theMovieDbInfos, err := fetchMovieSearchResults(searchURL)
	if err != nil || len(theMovieDbInfos.Results) == 0 {
		return nil
	}

	theMovieDbInfo := theMovieDbInfos.Results[0]
	movieInfo := &MovieInfo{
		Description: theMovieDbInfo.Overview,
	}

	if err := setReleaseDate(movieInfo, theMovieDbInfo.ReleaseDate); err != nil {
		log.Println("Error parsing release date:", err)
		return nil
	}

	// Fetch additional movie details (runtime) from TMDb API
	movieDetailsURL := fmt.Sprintf("%s/movie/%d?api_key=%s", baseURL, theMovieDbInfo.ID, apiKey)
	movieDetails, err := fetchMovieDetails(movieDetailsURL)
	if err != nil {
		return nil
	}

	movieInfo.Runtime = movieDetails.Runtime
	return movieInfo
}

func fetchMovieSearchResults(url string) (*TheMovieDbInfos, error) {
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

	var theMovieDbInfos TheMovieDbInfos
	if err := json.NewDecoder(res.Body).Decode(&theMovieDbInfos); err != nil {
		log.Println("Error decoding JSON response:", err)
		return nil, err
	}

	return &theMovieDbInfos, nil
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

func setReleaseDate(movieInfo *MovieInfo, releaseDate string) error {
	if releaseDate == "" {
		return nil
	}

	parsedDate, err := time.Parse("2006-01-02", releaseDate)
	if err != nil {
		return err
	}

	if parsedDate.After(time.Now()) {
		return fmt.Errorf("release date is in the future")
	}

	movieInfo.ReleaseDate = parsedDate.Format("02.01.06")
	return nil
}
