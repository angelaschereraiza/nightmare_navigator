package tmdb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nightmare_navigator/internal/config"
)

func TestGetTMDbInfoByTitle(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/search/movie") {
			query := r.URL.Query().Get("query")
			if query == "Title 1" {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"results":[{"id":1,"overview":"Overview 1","release_date":"2024-05-01","title":"Title 1"}]}`)
			} else {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"results":[]}`)
			}
		}
		if strings.Contains(r.URL.Path, "/movie/1") {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"runtime":120}`)
		}
	}))
	defer mockServer.Close()

	cfg := config.Config{
		TMDb: config.TMDb{
			ApiKey: "dummy_api_key",
			ApiURL: mockServer.URL,
		},
	}

	manager := NewTMDbManager(cfg)

	tests := []struct {
		title       string
		year        string
		expected    *TMDbInfo
		expectError bool
	}{
		{"Title 1", "2024", &TMDbInfo{ID: 1, Overview: "Overview 1", ReleaseDate: "01.05.24", Runtime: 120, Title: "Title 1"}, false},
		{"Non-existent Title", "2024", nil, true},
		{"Title 1", "2025", nil, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Title: %s, Year: %s", test.title, test.year), func(t *testing.T) {
			movieInfo := manager.GetTMDbInfoByTitle(test.title, test.year)

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
