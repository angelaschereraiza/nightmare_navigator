package movie_infos_test

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	movieinfos "nightmare_navigator/internal/movie_infos"
)

func TestSaveLatestIMDbRatings(t *testing.T) {
	// Setup mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data string
		if strings.Contains(r.URL.String(), movieinfos.BasicsFilename) {
			data = "tconst\ttitleType\tprimaryTitle\toriginalTitle\tisAdult\tstartYear\tendYear\truntimeMinutes\tgenres\n" +
				"tt0119081\tmovie\tEvent Horizon\tEvent Horizon\t0\t1997\t\\N\t96\tHorror,Sci-Fi\n" +
				"tt0391198\tmovie\tThe Grudge\tThe Grudge\t0\t2004\t\\N\t92\tHorror,Mystery,Thriller\n"
		} else if strings.Contains(r.URL.String(), movieinfos.RatingsFilename) {
			data = "tconst\taverageRating\tnumVotes\n" +
				"tt0119081\t6.7\t120000\n" +
				"tt0391198\t5.9\t90000\n"
		}
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		gz.Write([]byte(data))
		gz.Close()
		w.Write(buf.Bytes())
	}))
	defer server.Close()

	oldIMDbBaseURL := movieinfos.IMDbBaseURL
	movieinfos.IMDbBaseURL = server.URL + "/"
	defer func() { movieinfos.IMDbBaseURL = oldIMDbBaseURL }()

	tmpDir, err := os.MkdirTemp("", "test_movieinfos")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDownloadDir := movieinfos.DownloadDir
	movieinfos.DownloadDir = tmpDir
	defer func() { movieinfos.DownloadDir = oldDownloadDir }()

	movieinfos.SaveLatestIMDbRatings()

	jsonFile := filepath.Join(tmpDir, movieinfos.JSONFilename)
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		t.Fatalf("Expected JSON file to be created: %v", jsonFile)
	}

	fileContent, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	expectedSubstrings := []string{
		`"primaryTitle": "Event Horizon"`,
		`"primaryTitle": "The Grudge"`,
	}

	for _, substring := range expectedSubstrings {
		if !strings.Contains(string(fileContent), substring) {
			t.Errorf("Expected JSON to contain %q, got %s", substring, string(fileContent))
		}
	}
}
