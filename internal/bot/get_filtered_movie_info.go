package bot

import (
	"time"

	"nightmare_navigator/internal/config"
	"nightmare_navigator/pkg/omdb"
	movieinfo "nightmare_navigator/pkg/movie_info"
)

type GetIMDbInfosByDateAndGenreFunc func(config.Config, int, []string, time.Time, func(string) *movieinfo.MovieInfo) *[]movieinfo.MovieInfo
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

	movies := getIMDbInfosByDateAndGenre(cfg, count, genres, date, getOMDbInfoByTitle)

	if movies == nil || len(*movies) == 0 {
		return &[]string{}
	}

	return buildMovieInfoStrings(*movies)
}
