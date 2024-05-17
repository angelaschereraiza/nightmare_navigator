package manager

import (
	"fmt"
	"regexp"
	"strings"
)

func getAdditionalMovieInfo(imdbMovieInfos []IMDbMovieInfo) *[]string {
	var results []string

	for _, imdbMovieInfo := range imdbMovieInfos {
		// Gets additional movie infos
		omdbMovieInfo := getOMDbInfoByTitle(imdbMovieInfo.Title)
		theMovieDbInfo := getTheMovieDbInfoByTitle(imdbMovieInfo.Title)
		genreMap := getGenres()

		// If the title is not written in Latin characters, the movie is skipped
		if !containsLatinChars(imdbMovieInfo.Title) {
			continue
		}

		results = append(results, buildMovieInfoString(imdbMovieInfo, omdbMovieInfo, theMovieDbInfo, *genreMap))
	}

	return &results
}

func buildMovieInfoString(imdbMovieInfo IMDbMovieInfo, omdbMovieInfo *OMDbMovieInfo, theMovieDbInfo *TheMovieDbInfo, genreMap map[int]string) string {
	var result strings.Builder

	// Title
	result.WriteString(fmt.Sprintf("Title: %s\n", imdbMovieInfo.Title))

	// IMDb
	result.WriteString(fmt.Sprintf("IMDb Rating: %s\n", imdbMovieInfo.IMDb))
	result.WriteString(fmt.Sprintf("IMDb Votes: %s\n", imdbMovieInfo.IMDbVotes))
	result.WriteString(fmt.Sprintf("IMDb Link: https://www.imdb.com/title/%s\n", imdbMovieInfo.TitleId))

	if omdbMovieInfo != nil {
		// Country
		if omdbMovieInfo.Country != "" && omdbMovieInfo.Country != "N/A" {
			result.WriteString(fmt.Sprintf("Country: %s\n", omdbMovieInfo.Country))
		}

		// Rated
		if omdbMovieInfo.Rated != "" && omdbMovieInfo.Rated != "N/A" {
			result.WriteString(fmt.Sprintf("Rated: %s\n", omdbMovieInfo.Rated))
		}
	}

	if omdbMovieInfo != nil {
		// Genres
		result.WriteString("Genres: ")
		for i, genreID := range theMovieDbInfo.GenreIDs {
			genreName, found := genreMap[genreID]
			if found {
				result.WriteString(genreName)
				if i < len(theMovieDbInfo.GenreIDs)-1 {
					result.WriteString(", ")
				} else {
					result.WriteString("\n")
				}
			}
		}

		// Release Date
		result.WriteString(fmt.Sprintf("Released: %s\n", theMovieDbInfo.ReleaseDate))

		// Runtime
		if theMovieDbInfo.Runtime != 0 {
			result.WriteString(fmt.Sprintf("Runtime: %d minutes\n", theMovieDbInfo.Runtime))
		}

		// Description
		if theMovieDbInfo.Overview != "" {
			result.WriteString(fmt.Sprintf("Description: %s\n", theMovieDbInfo.Overview))
		}
	}

	return result.String()
}

func containsLatinChars(s string) bool {
	match, _ := regexp.MatchString("[a-zA-Z]", s)
	return match
}
