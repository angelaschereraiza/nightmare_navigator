package movie_infos

import (
	"time"
)

func GetFilteredMovieInfos(count int, genres []string, date time.Time) *[]string {
	if count <= 0 {
		return &[]string{}
	}

	movies := GetIMDbInfosByDateAndGenre(count, genres, date)

	if movies == nil || len(*movies) == 0 {
		return &[]string{}
	}

	return BuildMovieInfoStrings(*movies)
}
