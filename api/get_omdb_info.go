package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type OMDbMovieInfo struct {
	Title   string `json:"Title"`
	Rated   string `json:"Rated"`
	Country string `json:"Country"`
	Poster  string `json:"Poster"`
}

func GetOMDbInfoByTitle(name string) *OMDbMovieInfo {
	apiURL := "http://www.omdbapi.com/"
	apiKey := "d0bd48a2"

	params := map[string]string{
		"apikey": apiKey,
		"t":      name,
	}

	res, err := http.Get(apiURL + "?" + buildQueryString(params))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	var movie OMDbMovieInfo
	err = json.NewDecoder(res.Body).Decode(&movie)
	if err != nil {
		log.Println(err)
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
