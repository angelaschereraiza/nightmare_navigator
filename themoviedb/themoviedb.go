package themoviedb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"nightmare_navigator/omdb"
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
	ID            int    `json:"id"`
	Title         string `json:"title"`
	OriginalTitle string `json:"original_title"`
	ReleaseDate   string `json:"release_date"`
	GenreIDs      []int  `json:"genre_ids"`
	Overview      string `json:"overview"`
}

type MovieResponse struct {
	Results []Movie `json:"results"`
}

type MovieDetailsResponse struct {
	Runtime int `json:"runtime"`
}

func GetFilteredLatestMovies(count int, date time.Time) *[]string {
	baseURL := "https://api.themoviedb.org/3"
	apiKey := "6882b8441ce200fda300c1e46eeb3e64"

	horrorGenreID := 27
	excludeGenres := []int{12, 14, 10751, 99, 10402}

	// API request for genre names
	genreURL := fmt.Sprintf("%s/genre/movie/list?api_key=%s", baseURL, apiKey)
	genresRes, err := http.Get(genreURL)
	if err != nil {
		fmt.Println("Error fetching genre names:", err)
		return nil
	}
	defer genresRes.Body.Close()

	var genresResponse GenresResponse
	err = json.NewDecoder(genresRes.Body).Decode(&genresResponse)
	if err != nil {
		fmt.Println("Error decoding JSON response for genres:", err)
		return nil
	}

	// Create a map for mapping genre IDs to genre names
	genreMap := make(map[int]string)
	for _, genre := range genresResponse.Genres {
		genreMap[genre.ID] = genre.Name
	}

	var results []string
	totalPages := (count + 19) / 20 // Calculate number of pages needed
	for page := 1; page <= totalPages; page++ {
		// API request for movies for current page
		params := url.Values{}
		params.Set("api_key", apiKey)
		params.Set("with_genres", fmt.Sprintf("%d", horrorGenreID))
		params.Set("without_genres", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(excludeGenres)), ","), "[]"))
		params.Set("primary_release_date.lte", date.Format("2006-01-02"))
		params.Set("with_runtime.gte", "85")
		params.Set("sort_by", "primary_release_date.desc")
		params.Set("page", fmt.Sprintf("%d", page))

		searchURL := fmt.Sprintf("%s/discover/movie?%s", baseURL, params.Encode())

		res, err := http.Get(searchURL)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}
		defer res.Body.Close()

		var movieResponse MovieResponse
		err = json.NewDecoder(res.Body).Decode(&movieResponse)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)
			continue
		}

		for i, movie := range movieResponse.Results {
			if i >= count {
				break
			}

			var result strings.Builder
			// Additional API request to get movie details
			movieDetailsURL := fmt.Sprintf("%s/movie/%d?api_key=%s", baseURL, movie.ID, apiKey)
			detailsRes, err := http.Get(movieDetailsURL)
			if err != nil {
				fmt.Println("Error fetching movie details:", err)
				continue
			}
			defer detailsRes.Body.Close()

			var movieDetails MovieDetailsResponse
			err = json.NewDecoder(detailsRes.Body).Decode(&movieDetails)
			if err != nil {
				fmt.Println("Error decoding JSON response for movie details:", err)
				continue
			}

			// Retrieve additional information via the omdb api
			omdbMovieInformation := omdb.GetMovieByName(movie.Title)

			if !containsLatinChars(movie.Title) {
				continue
			}

			// Title
			result.WriteString(fmt.Sprintf("Title: %s\n", movie.Title))
			if movie.Title != movie.OriginalTitle {
				result.WriteString(fmt.Sprintf("Original Title: %s\n", movie.OriginalTitle))
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
				fmt.Println("Error decoding JSON response:", err)
				return nil
			}
			result.WriteString(fmt.Sprintf("Released: %s\n", releaseDate.Format("02.01.06")))

			// Runtime
			if movieDetails.Runtime != 0 {
				result.WriteString(fmt.Sprintf("Runtime: %d minutes\n", movieDetails.Runtime))
			}

			if omdbMovieInformation != nil && omdbMovieInformation.Title != "" {
				if omdbMovieInformation.Rated != "N/A" {
					result.WriteString(fmt.Sprintf("Rated: %s\n", omdbMovieInformation.Rated))
				}
				if omdbMovieInformation.Country != "N/A" {
					result.WriteString(fmt.Sprintf("Country: %s\n", omdbMovieInformation.Country))
				}
				if omdbMovieInformation.IMDb != "N/A" {
					result.WriteString(fmt.Sprintf("IMDb: %s\n", omdbMovieInformation.IMDb))
				}
				if omdbMovieInformation.ImdbVotes != "N/A" {
					result.WriteString(fmt.Sprintf("Imdb Votes: %s\n", omdbMovieInformation.ImdbVotes))
				}
				if omdbMovieInformation.MetaScore != "N/A" {
					result.WriteString(fmt.Sprintf("MetaScore: %s\n", omdbMovieInformation.MetaScore))
				}
			}

			// Description
			if movie.Overview != "" {
				result.WriteString(fmt.Sprintf("Description: %s\n", movie.Overview))
			}

			if omdbMovieInformation != nil && omdbMovieInformation.Title != "" {
				result.WriteString(fmt.Sprintf("Poster: %s\n", omdbMovieInformation.Poster))
			}

			results = append(results, result.String())
		}
	}

	return &results
}

func containsLatinChars(s string) bool {
	match, _ := regexp.MatchString("[a-zA-Z]", s)
	return match
}
