package movie_infos

import (
	"fmt"
	"regexp"
	"strings"
)

type MovieInfo struct {
	Country       string
	IMDb          string
	IMDbVotes     string
	Genres        string
	OriginalTitle string
	Description   string
	Rated         string
	ReleaseDate   string
	Runtime       int
	Title         string
	TitleId       string
	Year          string
}

func BuildMovieInfoStrings(movieInfos []MovieInfo) *[]string {
	results := make([]string, len(movieInfos))

	for i, movieInfo := range movieInfos {
		results[i] = BuildMovieInfoString(movieInfo)
	}
	
	return &results
}

func BuildMovieInfoString(movieInfo MovieInfo) string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("Title: %s\n", movieInfo.Title))

	if movieInfo.OriginalTitle != "" && !containsLatinChars(movieInfo.OriginalTitle) {
		result.WriteString(fmt.Sprintf("Original Title: %s\n", movieInfo.OriginalTitle))
	}

	result.WriteString(fmt.Sprintf("IMDb Rating: %s\n", movieInfo.IMDb))
	result.WriteString(fmt.Sprintf("IMDb Votes: %s\n", movieInfo.IMDbVotes))
	result.WriteString(fmt.Sprintf("IMDb Link: https://www.imdb.com/title/%s\n", movieInfo.TitleId))

	if movieInfo.Country != "" && movieInfo.Country != "N/A" {
		result.WriteString(fmt.Sprintf("Country: %s\n", movieInfo.Country))
	}

	if movieInfo.Rated != "" && movieInfo.Rated != "N/A" {
		result.WriteString(fmt.Sprintf("Rated: %s\n", movieInfo.Rated))
	}

	result.WriteString(fmt.Sprintf("Genres: %s\n", movieInfo.Genres))
	result.WriteString(fmt.Sprintf("Released: %s\n", movieInfo.ReleaseDate))

	if movieInfo.Runtime != 0 {
		result.WriteString(fmt.Sprintf("Runtime: %d minutes\n", movieInfo.Runtime))
	}

	if movieInfo.Description != "" {
		result.WriteString(fmt.Sprintf("Description: %s\n", movieInfo.Description))
	}

	return result.String()
}

func containsLatinChars(s string) bool {
	return regexp.MustCompile("[a-zA-Z]").MatchString(s)
}