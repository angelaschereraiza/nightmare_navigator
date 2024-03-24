package omdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type OMDBMovie struct {
	Title     string `json:"Title"`
	Rated     string `json:"Rated"`
	Country   string `json:"Country"`
	IMDb      string `json:"imdbRating"`
	MetaScore string `json:"Metascore"`
	ImdbVotes string `json:"imdbVotes"`
	Poster    string `json:"Poster"`
}

func GetMovieByName(name string) *OMDBMovie {
	apiURL := "http://www.omdbapi.com/"
	apiKey := "d0bd48a2"

	params := map[string]string{
		"apikey": apiKey,
		"t":      name,
	}

	res, err := http.Get(apiURL + "?" + buildQueryString(params))
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	// Antwort dekodieren
	var movie OMDBMovie
	err = json.NewDecoder(res.Body).Decode(&movie)
	if err != nil {
		return nil
	}

	return &movie
}

func buildQueryString(params map[string]string) string {
	var parts []string
	for key, value := range params {
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(parts, "&")
}
