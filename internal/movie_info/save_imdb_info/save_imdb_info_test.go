package save_imdb_info

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nightmare_navigator/internal/config"
)

func TestSaveLatestIMDbRatings(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data string
		if strings.Contains(r.URL.String(), "title.basics.tsv.gz") {
			data = "tconst\ttitleType\tprimaryTitle\toriginalTitle\tisAdult\tstartYear\tendYear\truntimeMinutes\tgenres\n" +
				"tt0119081\tmovie\tEvent Horizon\tEvent Horizon\t0\t1997\t\\N\t96\tHorror,Sci-Fi\n" +
				"tt0391198\tmovie\tThe Grudge\tThe Grudge\t0\t2004\t\\N\t92\tHorror,Mystery,Thriller\n"
		} else if strings.Contains(r.URL.String(), "title.ratings.tsv.gz") {
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
	defer mockServer.Close()

	cfg := config.Config{
		IMDb: config.IMDb{
			BasicsFilename:  "title.basics.tsv.gz",
			RatingsFilename: "title.ratings.tsv.gz",
			JSONFilename:    "imdb_movie_infos_test.json",
			DownloadDir:     "data",
			IMDbBaseUrl:     mockServer.URL + "/",
			MinRating:       5.0,
			MinVotes:        1000,
		},
		TMDb: config.TMDb{
			ApiKey: "6882b8441ce200fda300c1e46eeb3e64",
			ApiURL: "https://api.themoviedb.org/3",
		},
	}

	tmpDir, err := os.MkdirTemp("", "test_movieinfos")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg.IMDb.DownloadDir = tmpDir

	manager := NewSaveIMDbInfoManager(cfg)
	manager.SaveLatestIMDbRatings()

	jsonFile := filepath.Join(tmpDir, cfg.IMDb.JSONFilename)
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
