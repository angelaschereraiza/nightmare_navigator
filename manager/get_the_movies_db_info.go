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
	GenreIDs    []int  `json:"genre_ids"`
	Overview    string `json:"overview"`
	Runtime     int
}

type TheMovieDbInfos struct {
	Results []TheMovieDbInfo `json:"results"`
}

type MovieDetailsResponse struct {
	Runtime int `json:"runtime"`
}

func getTheMovieDbInfoByTitle(title string) *TheMovieDbInfo {
	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("query", title)

	searchURL := fmt.Sprintf("%s/search/movie?%s", baseURL, params.Encode())

	res, err := http.Get(searchURL)
	if err != nil {
		log.Println("Error making request:", err)
		return nil
	}
	defer res.Body.Close()

	var theMovieDbInfos TheMovieDbInfos
	err = json.NewDecoder(res.Body).Decode(&theMovieDbInfos)
	if err != nil {
		log.Println("Error decoding JSON response:", err)
		return nil
	}

	if len(theMovieDbInfos.Results) > 0 {
		theMovieDbInfo := theMovieDbInfos.Results[0]

		if theMovieDbInfo.ReleaseDate != "" {
			releaseDate, err := time.Parse("2006-01-02", theMovieDbInfo.ReleaseDate)
			if err != nil {
				log.Println("Error parsing release date:", err)
			}

			theMovieDbInfo.ReleaseDate = releaseDate.Format("02.01.06")
		}

		// Additional API request to get movie runtime
		movieDetailsURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?api_key=%s", theMovieDbInfo.ID, apiKey)
		detailsRes, err := http.Get(movieDetailsURL)
		if err != nil {
			log.Println("Error fetching movie details:", err)
		}
		defer detailsRes.Body.Close()

		var movieDetails MovieDetailsResponse
		err = json.NewDecoder(detailsRes.Body).Decode(&movieDetails)
		if err != nil {
			log.Println("Error decoding JSON response for movie details:", err)
		}

		theMovieDbInfo.Runtime = movieDetails.Runtime

		return &theMovieDbInfo
	}

	return nil
}

func getGenres() *map[int]string {
	genreURL := fmt.Sprintf("%s/genre/movie/list?api_key=%s", baseURL, apiKey)
	genresRes, err := http.Get(genreURL)
	if err != nil {
		log.Println("Error fetching genre names:", err)
		return nil
	}
	defer genresRes.Body.Close()

	var genresResponse GenresResponse
	err = json.NewDecoder(genresRes.Body).Decode(&genresResponse)
	if err != nil {
		log.Println("Error decoding JSON response for genres:", err)
		return nil
	}

	genreMap := make(map[int]string)
	for _, genre := range genresResponse.Genres {
		genreMap[genre.ID] = genre.Name
	}

	return &genreMap
}
