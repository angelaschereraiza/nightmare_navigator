package omdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Movie struct {
	Title     string `json:"Title"`
	Rated     string `json:"Rated"`
	Released  string `json:"Released"`
	Runtime   string `json:"Runtime"`
	Genres    string `json:"Genre"`
	Country   string `json:"Country"`
	IMDb      string `json:"imdbRating"`
	MetaScore string `json:"Metascore"`
	Plot      string `json:"Plot"`
	Director  string `json:"Director"`
	Writer    string `json:"Writer"`
	Actors    string `json:"Actors"`
	Language  string `json:"Language"`
	Awards    string `json:"Awards"`
	Poster    string `json:"Poster"`
	ImdbVotes string `json:"imdbVotes"`
	ImdbID    string `json:"imdbID"`
}

func GetFilteredLatestMoviesFromOmDB(date time.Time) *[]string {
	apiURL := "http://www.omdbapi.com/"
	apiKey := "d0bd48a2"

	// Berechnen Sie das Jahr des aktuellen Datums
	year := date.Format("2006")

	// API-Anfrage für Filme, die nach dem Veröffentlichungsjahr sortiert sind
	params := map[string]string{
		"apikey": apiKey,
		"s":      "horror",
		"type":   "movie",
		"y":      year, // Nur Filme des aktuellen Jahres abrufen
		"plot":   "short",
		"r":      "json", // Ergebnisse im JSON-Format erhalten
		"sort":   "released", // Nach Veröffentlichungsdatum sortieren
	}

	var results []string

	// HTTP-Anfrage ausführen
	res, err := http.Get(apiURL + "?" + buildQueryString(params))
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	// Antwort dekodieren
	var searchResponse struct {
		Search       []Movie `json:"Search"`
		TotalResults string  `json:"totalResults"`
		Response     string  `json:"Response"`
	}
	err = json.NewDecoder(res.Body).Decode(&searchResponse)
	if err != nil {
		return nil
	}

	for _, movie := range searchResponse.Search {
		result := fmt.Sprintf("Title: %s\n", movie.Title)
		result += fmt.Sprintf("Rated: %s\n", movie.Rated)
		result += fmt.Sprintf("Released: %s\n", movie.Released)
		result += fmt.Sprintf("Runtime: %s\n", movie.Runtime)
		result += fmt.Sprintf("Genres: %s\n", movie.Genres)
		result += fmt.Sprintf("Country: %s\n", movie.Country)
		result += fmt.Sprintf("IMDb: %s\n", movie.IMDb)
		result += fmt.Sprintf("MetaScore: %s\n", movie.MetaScore)
		result += fmt.Sprintf("Poster: %s\n", movie.Poster)
		result += fmt.Sprintf("IMDb Votes: %s\n", movie.ImdbVotes)
		result += fmt.Sprintf("IMDb ID: %s\n", movie.ImdbID)
		result += fmt.Sprintf("Plot: %s\n", movie.Plot)

		results = append(results, result)
	}

	return &results
}

func buildQueryString(params map[string]string) string {
	var parts []string
	for key, value := range params {
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(parts, "&")
}