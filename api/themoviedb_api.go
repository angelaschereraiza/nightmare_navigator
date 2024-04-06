package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"nightmare_navigator/imdb"
	"regexp"
	"strings"
	"time"
)

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GenresResponse struct {
	Genres []Genre `json:"genres"`
}

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	GenreIDs    []int  `json:"genre_ids"`
	Overview    string `json:"overview"`
}

type MovieResponse struct {
	Results []Movie `json:"results"`
}

type MovieDetailsResponse struct {
	Runtime int `json:"runtime"`
}

func GetFilteredLatestMovies(count int, genres []int, date time.Time) *[]string {
	baseURL := "https://api.themoviedb.org/3"
	apiKey := "6882b8441ce200fda300c1e46eeb3e64"

	excludeGenres := []int{12, 36, 10751, 99, 10402}

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

	var results []string
	addedMovies := make(map[string]bool)

	// Number of movies already collected
	collected := 0
	page := 1

	for collected < count {
		page ++

		params := url.Values{}
		params.Set("api_key", apiKey)
		params.Set("with_genres", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(genres)), ","), "[]"))
		params.Set("without_genres", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(excludeGenres)), ","), "[]"))
		params.Set("primary_release_date.lte", date.Format("2006-01-02"))
		params.Set("with_runtime.gte", "85")
		params.Set("sort_by", "primary_release_date.desc")
		params.Set("page", fmt.Sprintf("%d", page))
	
		searchURL := fmt.Sprintf("%s/discover/movie?%s", baseURL, params.Encode())

		res, err := http.Get(searchURL)
		if err != nil {
			log.Println("Error making request:", err)
			continue
		}
		defer res.Body.Close()

		var movieResponse MovieResponse
		err = json.NewDecoder(res.Body).Decode(&movieResponse)
		if err != nil {
			log.Println("Error decoding JSON response:", err)
			continue
		}

		for _, movie := range movieResponse.Results {
			if collected >= count {
				break
			}

			if addedMovies[movie.Title] {
                continue
            }

			var result strings.Builder

			// Additional API request to get movie details
			movieDetailsURL := fmt.Sprintf("%s/movie/%d?api_key=%s", baseURL, movie.ID, apiKey)
			detailsRes, err := http.Get(movieDetailsURL)
			if err != nil {
				log.Println("Error fetching movie details:", err)
				continue
			}
			defer detailsRes.Body.Close()

			var movieDetails MovieDetailsResponse
			err = json.NewDecoder(detailsRes.Body).Decode(&movieDetails)
			if err != nil {
				log.Println("Error decoding JSON response for movie details:", err)
				continue
			}

			// Retrieve additional information via the imdb rating list
			imdbMovieInfo := imdb.GetIMDbInfoByTitle(movie.Title)
			if imdbMovieInfo == nil {
				continue
			}

			// Retrieve additional information via the omdb api
			omdbMovieInfo := GetOMDbInfoByTitle(movie.Title)

			// If the tile is not written in Latin characters, the movie is skipped
			if !containsLatinChars(movie.Title) {
				continue
			}

			// Title
			result.WriteString(fmt.Sprintf("Title: %s\n", movie.Title))

			// IMDb
			result.WriteString(fmt.Sprintf("IMDb Rating: %s\n", imdbMovieInfo.IMDb))
			result.WriteString(fmt.Sprintf("IMDb Votes: %s\n", imdbMovieInfo.IMDbVotes))
			result.WriteString(fmt.Sprintf("IMDb Link: https://www.imdb.com/title/%s\n", imdbMovieInfo.TitleId))

			// OMDb
			if omdbMovieInfo != nil && omdbMovieInfo.Title != "" {
				if omdbMovieInfo.Rated != "N/A" {
					result.WriteString(fmt.Sprintf("Rated: %s\n", omdbMovieInfo.Rated))
				}
				if omdbMovieInfo.Country != "N/A" {
					result.WriteString(fmt.Sprintf("Country: %s\n", omdbMovieInfo.Country))
				}
			}

			// Genres
			result.WriteString("Genres: ")
			for i, genreID := range movie.GenreIDs {
				genreName, found := genreMap[genreID]
				if found {
					result.WriteString(genreName)
					if i < len(movie.GenreIDs)-1 {
						result.WriteString(", ")
					} else {
						result.WriteString("\n")
					}
				}
			}

			// Release Date
			releaseDate, err := time.Parse("2006-01-02", movie.ReleaseDate)
			if err != nil {
				log.Println("Error decoding JSON response:", err)
				return nil
			}
			result.WriteString(fmt.Sprintf("Released: %s\n", releaseDate.Format("02.01.06")))

			// Runtime
			if movieDetails.Runtime != 0 {
				result.WriteString(fmt.Sprintf("Runtime: %d minutes\n", movieDetails.Runtime))
			}

			// Description
			if movie.Overview != "" {
				result.WriteString(fmt.Sprintf("Description: %s\n", movie.Overview))
			}

			addedMovies[movie.Title] = true
			results = append(results, result.String())

			collected++
		}
	}

	return &results
}

func containsLatinChars(s string) bool {
	match, _ := regexp.MatchString("[a-zA-Z]", s)
	return match
}
