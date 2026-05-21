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
	mockGetIMDbInfosByYear := func(cfg config.Config, year string, getOMDbInfoByTitle func(string) *movieinfo.MovieInfo) []movieinfo.MovieInfo {
		return []movieinfo.MovieInfo{
			{Title: "Movie1", Year: year},
			{Title: "Movie2", Year: year},
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
	movies := []movieinfo.MovieInfo{
		{Title: "MovieIndia1", Country: "India"},
		{Title: "MovieIndia2", Country: "India, USA"},
		{Title: "MovieIndia3", Country: "UK, India"},
		{Title: "MovieOther", Country: "USA"},
	}

	filtered := filterAlreadyReturnedMovies(movies)
	if len(filtered) != 1 || filtered[0].Title != "MovieOther" {
		t.Fatalf("Expected only non-Indian movie to remain, got %v", filtered)
	}
}
