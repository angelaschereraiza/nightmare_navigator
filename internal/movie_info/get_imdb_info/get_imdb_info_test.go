package get_imdb_info

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"io/fs"
	movieinfo "nightmare_navigator/internal/movie_info"
)

func createTempIMDbJSON(t *testing.T, content string) {
	t.Helper()

	err := os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}

	err = os.WriteFile(filepath.Join(downloadDir, jsonFilename), []byte(content), fs.FileMode(0644))
	if err != nil {
		t.Fatalf("could not write temp file: %v", err)
	}
}

func TestLoadIMDbData(t *testing.T) {
	jsonContent := `
	[
		{
			"data": [
				{
					"averageRating": "6.7",
					"description": "A rescue crew investigates a spaceship that disappeared into a black hole and has now returned...with someone or something new on-board.",
					"genres": "Horror, Sci-Fi",
					"numVotes": "155651",
					"originalTitle": "Event Horizon",
					"primaryTitle": "Event Horizon",
					"releaseDate": "15.08.97",
					"runtime": 96,
					"tconst": "tt0119081"
				}
			],
			"startYear": "1997"
		}
	]`

	createTempIMDbJSON(t, jsonContent)
	defer os.RemoveAll(downloadDir)

	expected := []movieinfo.MovieInfo{
		{
			Description:   "A rescue crew investigates a spaceship that disappeared into a black hole and has now returned...with someone or something new on-board.",
			IMDb:          "6.7",
			IMDbVotes:     "155651",
			Genres:        "Horror, Sci-Fi",
			OriginalTitle: "Event Horizon",
			ReleaseDate:   "15.08.97",
			Runtime:       96,
			Title:         "Event Horizon",
			TitleId:       "tt0119081",
			Year:          "1997",
		},
	}

	movieInfos, err := loadIMDbData()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if !equal(movieInfos, expected) {
		t.Fatalf("Expected %v, but got %v", expected, movieInfos)
	}
}

func TestGetIMDbInfosByYear(t *testing.T) {
	jsonContent := `
	[
		{
			"data": [
				{
					"averageRating": "6.7",
					"description": "A rescue crew investigates a spaceship that disappeared into a black hole and has now returned...with someone or something new on-board.",
					"genres": "Horror, Sci-Fi",
					"numVotes": "155651",
					"originalTitle": "Event Horizon",
					"primaryTitle": "Event Horizon",
					"releaseDate": "15.08.97",
					"runtime": 96,
					"tconst": "tt0119081"
				}
			],
			"startYear": "1997"
		}
	]`

	createTempIMDbJSON(t, jsonContent)
	defer os.RemoveAll(downloadDir)

	expected := []movieinfo.MovieInfo{
		{
			Description:   "A rescue crew investigates a spaceship that disappeared into a black hole and has now returned...with someone or something new on-board.",
			IMDb:          "6.7",
			IMDbVotes:     "155651",
			Genres:        "Horror, Sci-Fi",
			OriginalTitle: "Event Horizon",
			ReleaseDate:   "15.08.97",
			Runtime:       96,
			Title:         "Event Horizon",
			TitleId:       "tt0119081",
			Year:          "1997",
			Rated:         "R",
			Country:       "UK, USA",
		},
	}

	mockGetOMDbInfoByTitle := func(title string) *movieinfo.MovieInfo {
		return &movieinfo.MovieInfo{
			Rated:   "R",
			Country: "UK, USA",
		}
	}

	moviesByYear := GetIMDbInfosByYear("1997", mockGetOMDbInfoByTitle)
	if !equal(moviesByYear, expected) {
		t.Fatalf("Expected %v, but got %v", expected, moviesByYear)
	}
}

func TestGetIMDbInfosByDateAndGenre(t *testing.T) {
	jsonContent := `
	[
		{
			"data": [
				{
					"averageRating": "6.7",
					"description": "A rescue crew investigates a spaceship that disappeared into a black hole and has now returned...with someone or something new on-board.",
					"genres": "Horror, Sci-Fi",
					"numVotes": "155651",
					"originalTitle": "Event Horizon",
					"primaryTitle": "Event Horizon",
					"releaseDate": "15.08.97",
					"runtime": 96,
					"tconst": "tt0119081"
				}
			],
			"startYear": "1997"
		}
	]`

	createTempIMDbJSON(t, jsonContent)
	defer os.RemoveAll(downloadDir)

	expected := []movieinfo.MovieInfo{
		{
			Description:   "A rescue crew investigates a spaceship that disappeared into a black hole and has now returned...with someone or something new on-board.",
			IMDb:          "6.7",
			IMDbVotes:     "155651",
			Genres:        "Horror, Sci-Fi",
			OriginalTitle: "Event Horizon",
			ReleaseDate:   "15.08.97",
			Runtime:       96,
			Title:         "Event Horizon",
			TitleId:       "tt0119081",
			Year:          "1997",
			Rated:         "R",
			Country:       "UK, USA",
		},
	}

	mockGetOMDbInfoByTitle := func(title string) *movieinfo.MovieInfo {
		return &movieinfo.MovieInfo{
			Rated:   "R",
			Country: "UK, USA",
		}
	}

	date, _ := time.Parse("02.01.06", "01.01.21")
	result := GetIMDbInfosByDateAndGenre(1, []string{"Horror", "Sci-Fi"}, date, mockGetOMDbInfoByTitle)
	if !equal(*result, expected) {
		t.Fatalf("Expected %v, but got %v", expected, *result)
	}
}

func equal(a, b []movieinfo.MovieInfo) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
