package movie_info

import (
	"time"

	"nightmare_navigator/internal/config"
	"nightmare_navigator/internal/omdb"
)

type GetIMDbInfosByDateAndGenreFunc func(config.Config, int, []string, time.Time, func(string) *MovieInfo) *[]MovieInfo
type BuildMovieInfoStringsFunc func([]MovieInfo) *[]string

func GetFilteredMovieInfos(count int, genres []string, date time.Time, getIMDbInfosByDateAndGenre GetIMDbInfosByDateAndGenreFunc, buildMovieInfoStrings BuildMovieInfoStringsFunc, cfg config.Config) *[]string {
	if count <= 0 {
		return &[]string{}
	}

	getOMDbInfoByTitle := func(title string) *MovieInfo {
		manager := omdb.NewOMDbManager(cfg)
		omdbInfo := manager.GetOMDbInfoByTitle(title)
		if omdbInfo == nil {
			return nil
		}
		return &MovieInfo{
			Rated:   omdbInfo.Rated,
			Country: omdbInfo.Country,
		}
	}

	movies := getIMDbInfosByDateAndGenre(cfg, count, genres, date, getOMDbInfoByTitle)

	if movies == nil || len(*movies) == 0 {
		return &[]string{}
	}

	return buildMovieInfoStrings(*movies)
}
