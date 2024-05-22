package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Default values
const defaultCount = 10
var defaultGenres = []string{"Horror"}

// Precompiled regular expressions for genre extraction
var genreRegex = map[string]*regexp.Regexp{
	"Sci-Fi":    regexp.MustCompile(`\b(sci[\s-]?fi)\b`),
	"Fantasy":   regexp.MustCompile(`\b(fantasy)\b`),
	"Thriller":  regexp.MustCompile(`\b(thriller)\b`),
	"Animation": regexp.MustCompile(`\b(animation)\b`),
	"Mystery":   regexp.MustCompile(`\b(mystery)\b`),
}

func ExtractCount(text string) int {
	numStr := ""
	for _, char := range text {
		if unicode.IsDigit(char) {
			numStr += string(char)
		} else if len(numStr) > 0 {
			break
		}
	}

	if numStr == "" {
		return defaultCount
	}

	count, err := strconv.Atoi(numStr)
	if err != nil {
		log.Printf("Error converting %s to int: %v", numStr, err)
		return defaultCount
	}
	return count
}

// ExtractGenres extracts predefined genres from the text. Defaults to ["Horror"].
func ExtractGenres(text string) []string {
	genres := make([]string, len(defaultGenres))
	copy(genres, defaultGenres)

	lowerText := strings.ToLower(text)
	for genre, regex := range genreRegex {
		if regex.MatchString(lowerText) {
			genres = append(genres, genre)
		}
	}
	return genres
}

// ExtractDate extracts the first date in the format "DD.MM.YY" from the text. Defaults to the current date.
func ExtractDate(text string) time.Time {
	const dateFormat = "02.01.06"
	now := time.Now()

	if text == "" {
		return now
	}

	re := regexp.MustCompile(`\d{2}\.\d{2}\.\d{2}`)
	dateStr := re.FindString(text)
	if dateStr == "" {
		return now
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		log.Printf("Error parsing date %s: %v", dateStr, err)
		return now
	}

	return date
}