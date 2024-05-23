package get_filtered_movie_info

import (
	"time"

	movieinfo "nightmare_navigator/internal/movie_info"
	omdb "nightmare_navigator/internal/movie_info/get_omdb_info"
	"nightmare_navigator/internal/config"
)

type GetIMDbInfosByDateAndGenreFunc func(int, []string, time.Time, func(string) *movieinfo.MovieInfo) *[]movieinfo.MovieInfo
type BuildMovieInfoStringsFunc func([]movieinfo.MovieInfo) *[]string

func GetFilteredMovieInfos(count int, genres []string, date time.Time, getIMDbInfosByDateAndGenre GetIMDbInfosByDateAndGenreFunc, buildMovieInfoStrings BuildMovieInfoStringsFunc, cfg config.Config) *[]string {
	if count <= 0 {
		return &[]string{}
	}

	getOMDbInfoByTitle := func(title string) *movieinfo.MovieInfo {
		manager := omdb.NewOMDbManager(cfg)
		omdbInfo := manager.GetOMDbInfoByTitle(title)
		if omdbInfo == nil {
			return nil
		}
		return &movieinfo.MovieInfo{
			Rated:   omdbInfo.Rated,
			Country: omdbInfo.Country,
		}
	}

	movies := getIMDbInfosByDateAndGenre(count, genres, date, getOMDbInfoByTitle)

	if movies == nil || len(*movies) == 0 {
		return &[]string{}
	}

	return buildMovieInfoStrings(*movies)
}
