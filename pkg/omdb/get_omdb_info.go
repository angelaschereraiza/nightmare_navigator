package omdb

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"nightmare_navigator/internal/config"
)

type OMDbMovieInfo struct {
	Title    string `json:"Title"`
	Rated    string `json:"Rated"`
	Country  string `json:"Country"`
	Response string `json:"Response"`
	Error    string `json:"Error"`
}

type OMDbManager struct {
	cfg config.Config
}

func NewOMDbManager(cfg config.Config) *OMDbManager {
	return &OMDbManager{cfg: cfg}
}

func (mgr *OMDbManager) GetOMDbInfoByTitle(title string) *OMDbMovieInfo {
	params := url.Values{}
	params.Add("apikey", mgr.cfg.OMDb.ApiKey)
	params.Add("t", title)

	res, err := http.Get(mgr.cfg.OMDb.ApiURL + "?" + params.Encode())
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

	if movie.Response == "False" {
		log.Println("OMDb API response error:", movie.Error)
		return nil
	}

	return &movie
}
