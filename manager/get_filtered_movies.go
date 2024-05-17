package manager

import (
	"strconv"
	"time"
)

func GetFilteredMovies(count *int, genres *string, date *time.Time) *[]string {
	// Gets all movies from imdb ratings json where the year is the current year
	imdbRatingMovies := getIMDbMovieInfosByYear(strconv.Itoa(time.Now().Year()))

	// Check if imdbRatinsMovie is contained in alreadyReturnedMovies and removes it from the list if true
	newMovies := filterAlreadyReturnedMovies(imdbRatingMovies)

	// Adds the title of the new movies into alreadyReturnedMovies list
	for _, movie := range newMovies {
		alreadyReturnedMovies[movie.Title] = true
	}

	// Gets additional movie information
	newMoviesWithAllInfo := getAdditionalMovieInfo(newMovies)

	return newMoviesWithAllInfo
}
