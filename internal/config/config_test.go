package config

import (
	"os"
	"testing"
)

// Helper function to create a temporary file with the provided content
func createTempFile(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpfile.Name()
}

func TestLoadConfig(t *testing.T) {
	// Success case
	validConfig := `
telegram_bot:
  token: "12345"
  channel_name: "test_channel"
general:
  data_dir: "data"
  already_returned_movies_json: "movies.json"
imdb:
  basics_filename: "title.basics.tsv.gz"
  json_filename: "imdb_movie_infos.json"
  min_rating: 5.0
  min_votes: 1000
  ratings_filename: "title.ratings.tsv.gz"
  imdb_base_url: "https://datasets.imdbws.com/"
tmdb:
  api_key: "tmdb_api_key"
  api_url: "https://api.themoviedb.org/3"
omdb:
  api_key: "omdb_api_key"
  api_url: "http://www.omdbapi.com/"
`
	tmpfile := createTempFile(t, validConfig)
	defer os.Remove(tmpfile)

	config, err := LoadConfig(tmpfile)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.TelegramBot.Token != "12345" {
		t.Errorf("expected token to be '12345', got %s", config.TelegramBot.Token)
	}

	if config.TelegramBot.ChannelName != "test_channel" {
		t.Errorf("expected channel name to be 'test_channel', got %s", config.TelegramBot.ChannelName)
	}

	// Failure case: File not found
	_, err = LoadConfig("non_existent_file.yaml")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Failure case: Invalid YAML content
	invalidConfig := `
telegram_bot:
  token: "12345
  channel_name: "test_channel"
`
	tmpfile = createTempFile(t, invalidConfig)
	defer os.Remove(tmpfile)

	_, err = LoadConfig(tmpfile)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}