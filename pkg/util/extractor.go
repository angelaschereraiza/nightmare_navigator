package util

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const defaultCount = 10

var defaultGenres = []string{"Horror"}

var genreRegex = map[string]*regexp.Regexp{
	"Sci-Fi":    regexp.MustCompile(`\b(sci[\s-]?fi)\b`),
	"Fantasy":   regexp.MustCompile(`\b(fantasy)\b`),
	"Thriller":  regexp.MustCompile(`\b(thriller)\b`),
	"Animation": regexp.MustCompile(`\b(animation)\b`),
	"Mystery":   regexp.MustCompile(`\b(mystery)\b`),
}

// ExtractCount extracts the number of movies to be returned
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

// ExtractGenres extracts predefined genres from the text. Defaults to ["Horror"]
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

// ExtractDate extracts the first date from the text and accepts formats like
// "14.06.19", "14.06.2019", "2012", "01. Januar 2012" or "Januar 2012".
// Defaults to the current date.
func ExtractDate(text string) time.Time {
	now := time.Now()
	if text == "" {
		return now
	}

	lowerText := strings.ToLower(text)

	if date := parseNumericDate(lowerText); !date.IsZero() {
		return date
	}

	if date := parseDayMonthYear(lowerText); !date.IsZero() {
		return date
	}

	if date := parseMonthYear(lowerText); !date.IsZero() {
		return date
	}

	if year := parseYearOnly(lowerText); year != 0 {
		return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	}

	return now
}

func parseNumericDate(text string) time.Time {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\b\d{1,2}[.\-/]\d{1,2}[.\-/](?:\d{2}|\d{4})\b`),
		regexp.MustCompile(`\b(?:\d{4})[.\-/]\d{1,2}[.\-/]\d{1,2}\b`),
	}

	formats := []string{
		"2.1.2006",
		"02.01.2006",
		"2-1-2006",
		"02-01-2006",
		"2/1/2006",
		"02/01/2006",
		"2.1.06",
		"02.01.06",
		"2-1-06",
		"02-01-06",
		"2/1/06",
		"02/01/06",
		"2006-1-2",
		"2006-01-02",
		"2006/1/2",
		"2006/01/02",
		"2006.1.2",
		"2006.01.02",
	}

	for _, pattern := range patterns {
		if dateStr := pattern.FindString(text); dateStr != "" {
			for _, format := range formats {
				if date, err := time.Parse(format, dateStr); err == nil {
					return date
				}
			}
		}
	}

	return time.Time{}
}

func parseDayMonthYear(text string) time.Time {
	re := regexp.MustCompile(`(?i)\b(\d{1,2})\s*[.\-/]?\s*([a-zäöüß]+)\s*[.\-/]?\s*(\d{4})\b`)
	m := re.FindStringSubmatch(text)
	if len(m) != 4 {
		return time.Time{}
	}

	monthName, ok := normalizeMonthName(m[2])
	if !ok {
		return time.Time{}
	}

	dateStr := fmt.Sprintf("%s %s %s", m[1], monthName, m[3])
	date, err := time.Parse("2 January 2006", dateStr)
	if err != nil {
		return time.Time{}
	}

	return date
}

func parseMonthYear(text string) time.Time {
	re := regexp.MustCompile(`(?i)\b([a-zäöüß]+)\s+(\d{4})\b`)
	m := re.FindStringSubmatch(text)
	if len(m) != 3 {
		return time.Time{}
	}

	monthName, ok := normalizeMonthName(m[1])
	if !ok {
		return time.Time{}
	}

	dateStr := fmt.Sprintf("%s %s", monthName, m[2])
	date, err := time.Parse("January 2006", dateStr)
	if err != nil {
		return time.Time{}
	}

	return date
}

func parseYearOnly(text string) int {
	re := regexp.MustCompile(`\b(19\d{2}|20\d{2})\b`)
	m := re.FindStringSubmatch(text)
	if len(m) != 2 {
		return 0
	}

	year, err := strconv.Atoi(m[1])
	if err != nil {
		return 0
	}

	return year
}

func normalizeMonthName(raw string) (string, bool) {
	key := strings.ToLower(strings.TrimSpace(raw))
	key = strings.TrimSuffix(key, ".")
	key = strings.ReplaceAll(key, "ä", "a")
	key = strings.ReplaceAll(key, "ö", "o")
	key = strings.ReplaceAll(key, "ü", "u")
	key = strings.ReplaceAll(key, "ß", "ss")

	months := map[string]string{
		"jan":       "January",
		"january":   "January",
		"januar":    "January",
		"feb":       "February",
		"februar":   "February",
		"mar":       "March",
		"march":     "March",
		"marz":      "March",
		"april":     "April",
		"apr":       "April",
		"may":       "May",
		"mai":       "May",
		"jun":       "June",
		"juni":      "June",
		"jul":       "July",
		"juli":      "July",
		"aug":       "August",
		"august":    "August",
		"sep":       "September",
		"sept":      "September",
		"september": "September",
		"okt":       "October",
		"oct":       "October",
		"oktober":   "October",
		"nov":       "November",
		"november":  "November",
		"dez":       "December",
		"dec":       "December",
		"dezember":  "December",
	}

	monthName, ok := months[key]
	return monthName, ok
}
