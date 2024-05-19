package manager

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

const (
	omdbApiURL = "http://www.omdbapi.com/"
	omdbApiKey = "d0bd48a2"
)

type OMDbMovieInfo struct {
	Title   string `json:"Title"`
	Rated   string `json:"Rated"`
	Country string `json:"Country"`
}

func getOMDbInfoByTitle(title string) *OMDbMovieInfo {
	params := url.Values{}
	params.Add("apikey", omdbApiKey)
	params.Add("t", title)

	res, err := http.Get(omdbApiURL + "?" + params.Encode())
	if err != nil {
		log.Println("HTTP request failed: ", err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("HTTP request failed with status code: ", res.StatusCode)
		return nil
	}

	var movie OMDbMovieInfo
	if err := json.NewDecoder(res.Body).Decode(&movie); err != nil {
		log.Println("Failed to decode response body:", err)
		return nil
	}

	return &movie
}