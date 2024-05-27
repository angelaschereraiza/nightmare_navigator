package movie_info

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"nightmare_navigator/internal/config"
	"testing"
)

func TestGetOMDbInfoByTitle(t *testing.T) {
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

	manager := NewOMDbManager(cfg)

	tests := []struct {
		title       string
		expected    *OMDbMovieInfo
		expectError bool
	}{
		{"Event Horizon", &OMDbMovieInfo{Title: "Event Horizon", Rated: "R", Country: "UK, USA"}, false},
		{"NonExistentMovie", nil, true},
		{"", nil, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Title: %s", test.title), func(t *testing.T) {
			movieInfo := manager.GetOMDbInfoByTitle(test.title)

			if test.expectError {
				if movieInfo != nil {
					t.Errorf("Expected error, but got movie info: %v", movieInfo)
				}
			} else {
				if movieInfo == nil {
					t.Error("Expected movie info, but got nil")
				} else {
					if movieInfo.Title != test.expected.Title ||
						movieInfo.Rated != test.expected.Rated ||
						movieInfo.Country != test.expected.Country {
						t.Errorf("Unexpected movie info. Expected: %v, Actual: %v", test.expected, movieInfo)
					}
				}
			}
		})
	}
}
