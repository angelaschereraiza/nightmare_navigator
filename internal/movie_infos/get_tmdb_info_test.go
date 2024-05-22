package movie_infos_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	movieinfos "nightmare_navigator/internal/movie_infos"
	"testing"
)

func TestGetTMDbInfoByTitle(t *testing.T) {
	// Mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve search results
		if r.URL.Path == "/search/movie" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"results":[{"id":1,"overview":"Overview 1","release_date":"2024-05-01","title":"Title 1"}]}`)
		}
		// Serve movie details
		if r.URL.Path == "/movie/1" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"runtime":120}`)
		}
	}))
	defer mockServer.Close()

	// Override base URL
	movieinfos.BaseURL = mockServer.URL

	// Test cases
	tests := []struct {
		title       string
		year        string
		expected    *movieinfos.TMDbInfo
		expectError bool
	}{
		{"Title 1", "2024", &movieinfos.TMDbInfo{ID: 1, Overview: "Overview 1", ReleaseDate: "01.05.24", Runtime: 120, Title: "Title 1"}, false},
		{"Non-existent Title", "2024", nil, true},
		{"Title 1", "2025", nil, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Title: %s, Year: %s", test.title, test.year), func(t *testing.T) {
			movieInfo := movieinfos.GetTMDbInfoByTitle(test.title, test.year)

			if test.expectError {
				if movieInfo != nil {
					t.Errorf("Expected error, but got movie info: %v", movieInfo)
				}
			} else {
				if movieInfo == nil {
					t.Error("Expected movie info, but got nil")
				} else {
					if movieInfo.ID != test.expected.ID ||
						movieInfo.Overview != test.expected.Overview ||
						movieInfo.ReleaseDate != test.expected.ReleaseDate ||
						movieInfo.Runtime != test.expected.Runtime ||
						movieInfo.Title != test.expected.Title {
						t.Errorf("Unexpected movie info. Expected: %v, Actual: %v", test.expected, movieInfo)
					}
				}
			}
		})
	}
}
