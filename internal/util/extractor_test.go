package util

import (
	"testing"
	"time"
)

func TestExtractCount(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"There are 15 items", 15},
		{"Count: 25", 25},
		{"No numbers here", defaultCount},
		{"123abc456", 123},
		{"", defaultCount},
	}

	for _, test := range tests {
		result := ExtractCount(test.input)
		if result != test.expected {
			t.Errorf("ExtractCount(%q) = %d; want %d", test.input, result, test.expected)
		}
	}
}

func TestExtractGenres(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"Sci-Fi and Fantasy movies", []string{"Horror", "Sci-Fi", "Fantasy"}},
		{"Mystery and Thriller are great", []string{"Horror", "Thriller", "Mystery"}},
		{"Animation movies are sugoi", []string{"Horror", "Animation"}},
		{"No genres mentioned", []string{"Horror"}},
		{"Sci-Fi is blyat", []string{"Horror", "Sci-Fi"}},
	}

	for _, test := range tests {
		result := ExtractGenres(test.input)
		if !equalStringSlices(result, test.expected) {
			t.Errorf("ExtractGenres(%q) = %v; want %v", test.input, result, test.expected)
		}
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

func TestExtractDate(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
	}{
		{"The event is on 25.12.21", parseDate("25.12.21")},
		{"Date: 01.01.20", parseDate("01.01.20")},
		{"No date here", time.Now()},
		{"", time.Now()},
	}

	for _, test := range tests {
		result := ExtractDate(test.input)
		if !datesAreClose(result, test.expected) {
			t.Errorf("ExtractDate(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func parseDate(dateStr string) time.Time {
	const dateFormat = "02.01.06"
	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		panic(err)
	}
	return date
}

func datesAreClose(a, b time.Time) bool {
	diff := a.Sub(b)
	return diff < time.Second && diff > -time.Second
}
