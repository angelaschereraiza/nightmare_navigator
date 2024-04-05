package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ExtractCount(text string) int {
	count := 20
	numStr := ""
	for _, char := range text {
		if unicode.IsDigit(char) {
			numStr += string(char)
		} else if numStr != "" {
			break
		}
	}

	if numStr != "" {
		count, _ = strconv.Atoi(numStr)
	}

	return count
}

func ExtractGenres(text string) []int {
	genres := []int{27}

	genreRegex := map[int]*regexp.Regexp{
		878: regexp.MustCompile(`\b(sci[\s-]?fi)\b`),
		14:  regexp.MustCompile(`\b(fantasy)\b`),
		53:  regexp.MustCompile(`\b(thriller)\b`),
		16:  regexp.MustCompile(`\b(animation)\b`),
	}

	for genre, regex := range genreRegex {
		if regex.MatchString(strings.ToLower(text)) {
			genres = append(genres, genre)
		}
	}

	return genres
}

func ExtractDate(text string) time.Time {
	now := time.Now()

	if text == "" {
		return now
	}

	re := regexp.MustCompile(`\d{2}\.\d{2}\.\d{2}`)
	dateStr := re.FindString(text)
	if dateStr == "" {
		return now
	}

	date, err := time.Parse("02.01.06", dateStr)
	if err != nil {
		log.Println(err)
		return now
	}

	return date
}
