package api

import (
	"fmt"
	"log"
	"nightmare_navigator/imdb"
	"regexp"
	"strings"
	"time"
)

func GetAdditionalMovieInfo(imdbMovieInfos []imdb.IMDbMovieInfo) *[]string {
	var results []string

	for _, imdbMovieInfo := range imdbMovieInfos {
		// Gets additional movie infos
		omdbMovieInfo := GetOMDbInfoByTitle(imdbMovieInfo.Title)
		theMovieDbInfo := GetTheMovieDbInfoByTitle(imdbMovieInfo.Title)
		genreMap := GetGenres()

		// Checks if release date is newer than current and skip if true
		releaseDate, err := time.Parse("02.01.06", theMovieDbInfo.ReleaseDate)
		if err != nil {
			log.Println(err)
			continue
		}
		if releaseDate.After(time.Now()) {
			continue
		}

		// If the title is not written in Latin characters, the movie is skipped
		if !containsLatinChars(imdbMovieInfo.Title) {
			continue
		}

		results = append(results, buildMovieInfoString(imdbMovieInfo, omdbMovieInfo, theMovieDbInfo, *genreMap))
	}

	return &results
}

func buildMovieInfoString(imdbMovieInfo imdb.IMDbMovieInfo, omdbMovieInfo *OMDbMovieInfo, theMovieDbInfo *TheMovieDbInfo, genreMap map[int]string) string {
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
