package movie_info

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nightmare_navigator/internal/config"
)

func mockGetIMDbInfosByDateAndGenre(cfg config.Config, count int, genres []string, date time.Time, getOMDbInfoByTitle func(string) *MovieInfo) *[]MovieInfo {
	movie := MovieInfo{
		Title:       "Event Horizon",
		Genres:      "Horror, Sci-Fi",
		ReleaseDate: "15.08.97",
		Year:        "1997",
	}
	return &[]MovieInfo{movie}
}

func mockBuildMovieInfoStrings(movies []MovieInfo) *[]string {
	result := []string{"Event Horizon (1997) - Horror, Sci-Fi"}
	return &result
}

func TestGetFilteredMovieInfos(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		switch query.Get("t") {
		case "Event Horizon":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"Title":"Event Horizon","Rated":"R","Country":"UK, USA"}`)
		case "NonExistentMovie":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"Response":"False","Error":"Movie not found!"}`)
		default:
			http.Error(w, "Movie not found", http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	cfg := config.Config{
		OMDb: config.OMDb{
			ApiKey: "dummy_api_key",
			ApiURL: mockServer.URL,
		},
	}

	date, _ := time.Parse("02.01.06", "01.01.21")
	result := GetFilteredMovieInfos(1, []string{"Horror", "Sci-Fi"}, date, mockGetIMDbInfosByDateAndGenre, mockBuildMovieInfoStrings, cfg)
	expected := &[]string{"Event Horizon (1997) - Horror, Sci-Fi"}

	if !equalStringSlices(*result, *expected) {
		t.Fatalf("Expected %v, but got %v", *expected, *result)
	}
}

func equalStringSlices(a, b []string) bool {
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
