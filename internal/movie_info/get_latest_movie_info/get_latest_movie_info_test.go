package get_latest_movie_info

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"nightmare_navigator/internal/config"
	movieinfo "nightmare_navigator/internal/movie_info"
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

	expected := []string{"Movie1 (2024)", "Movie2 (2024)"}
	manager := NewLatestMoviesManager(cfg)
	movieStrings := manager.GetLatestMovieInfos(mockGetIMDbInfosByYear, mockBuildMovieInfoStrings)

	if !reflect.DeepEqual(*movieStrings, expected) {
		t.Errorf("Expected %v, but got %v", expected, *movieStrings)
	}
}
