package manager

import (
	"time"
)

func GetFilteredMovies(count int, genres []string, date time.Time) *[]string {
	if count <= 0 {
		return &[]string{}
	}

	movies := getIMDbInfosByDateAndGenre(count, genres, date)

	if movies == nil || len(*movies) == 0 {
		return &[]string{}
	}

	return buildMovieInfoStrings(*movies)
}
