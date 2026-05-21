package bot

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"nightmare_navigator/internal/config"
	movieinfo "nightmare_navigator/pkg/movie_info"
)

func TestGetLatestMovieInfos(t *testing.T) {
	alreadyReturnedMovies = make(map[string]bool)
	mockGetIMDbInfosByYear := func(cfg config.Config, year string, getOMDbInfoByTitle func(string) *movieinfo.MovieInfo) []movieinfo.MovieInfo {
		return []movieinfo.MovieInfo{
			{Title: "Movie1", Year: year, ReleaseDate: time.Now().AddDate(0, -2, 0).Format("02.01.06")},
			{Title: "Movie2", Year: year, ReleaseDate: time.Now().AddDate(0, -2, 0).Format("02.01.06")},
		}
	}

	mockBuildMovieInfoStrings := func(movies []movieinfo.MovieInfo) *[]string {
		movieStrings := []string{}
		for _, movie := range movies {
			movieStrings = append(movieStrings, fmt.Sprintf("%s (%s)", movie.Title, movie.Year))
		}
		return &movieStrings
	}

	cfg := config.Config{
		General: config.General{
			DataDir:                   "data",
			AlreadyReturnedMoviesJSON: "already_returned_movies.json",
		},
	}

	tmpDir, err := os.MkdirTemp("", "testdata")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg.General.DataDir = tmpDir
	alreadyReturnedMoviesFile = filepath.Join(cfg.General.DataDir, cfg.General.AlreadyReturnedMoviesJSON)

	if err := createEmptyJsonFile(alreadyReturnedMoviesFile); err != nil {
		t.Fatalf("Failed to create empty JSON file: %v", err)
	}

	currentYear := strconv.Itoa(time.Now().Year())
	expected := []string{fmt.Sprintf("Movie1 (%s)", currentYear), fmt.Sprintf("Movie2 (%s)", currentYear)}
	manager := NewLatestMoviesManager(cfg)
	movieStrings := manager.GetLatestMovieInfos(mockGetIMDbInfosByYear, mockBuildMovieInfoStrings)

	if !reflect.DeepEqual(*movieStrings, expected) {
		t.Errorf("Expected %v, but got %v", expected, *movieStrings)
	}
}

func TestIsIndianCountryFilter(t *testing.T) {
	alreadyReturnedMovies = make(map[string]bool)
	movies := []movieinfo.MovieInfo{
		{Title: "MovieIndia1", Country: "India", ReleaseDate: time.Now().AddDate(0, -2, 0).Format("02.01.06")},
		{Title: "MovieIndia2", Country: "India, USA", ReleaseDate: time.Now().AddDate(0, -2, 0).Format("02.01.06")},
		{Title: "MovieIndia3", Country: "UK, India", ReleaseDate: time.Now().AddDate(0, -2, 0).Format("02.01.06")},
		{Title: "MovieOther", Country: "USA", ReleaseDate: time.Now().AddDate(0, -2, 0).Format("02.01.06")},
	}

	filtered := filterAlreadyReturnedMovies(movies)
	if len(filtered) != 1 || filtered[0].Title != "MovieOther" {
		t.Fatalf("Expected only non-Indian movie to remain, got %v", filtered)
	}
}

func TestOlderThanOneMonthFilter(t *testing.T) {
	alreadyReturnedMovies = make(map[string]bool)
	youngMovieDate := time.Now().AddDate(0, 0, -10).Format("02.01.06")
	oldMovieDate := time.Now().AddDate(0, -2, 0).Format("02.01.06")
	movies := []movieinfo.MovieInfo{
		{Title: "RecentMovie", Country: "USA", ReleaseDate: youngMovieDate},
		{Title: "OldMovie", Country: "USA", ReleaseDate: oldMovieDate},
	}

	filtered := filterAlreadyReturnedMovies(movies)
	if len(filtered) != 1 || filtered[0].Title != "OldMovie" {
		t.Fatalf("Expected only old movie to remain, got %v", filtered)
	}
}
