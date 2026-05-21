package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"nightmare_navigator/internal/config"
	movieinfo "nightmare_navigator/pkg/movie_info"
	"nightmare_navigator/pkg/omdb"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var alreadyReturnedMovies = make(map[string]bool)
var alreadyReturnedMoviesFile = ""

type GetIMDbInfosByYearFunc func(config.Config, string, func(string) *movieinfo.MovieInfo) []movieinfo.MovieInfo

type LatestMoviesManager struct {
	cfg config.Config
}

func NewLatestMoviesManager(cfg config.Config) *LatestMoviesManager {
	return &LatestMoviesManager{cfg: cfg}
}

func (mgr *LatestMoviesManager) GetLatestMovieInfos(getIMDbInfosByYear GetIMDbInfosByYearFunc, buildMovieInfoStrings BuildMovieInfoStringsFunc) *[]string {
	alreadyReturnedMoviesFile = filepath.Join(mgr.cfg.General.DataDir, mgr.cfg.General.AlreadyReturnedMoviesJSON)

	if err := os.MkdirAll(mgr.cfg.General.DataDir, 0755); err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	if err := loadAlreadyReturnedMovies(); err != nil {
		log.Println("Error loading already returned movies:", err)
	}

	getOMDbInfoByTitle := func(title string) *movieinfo.MovieInfo {
		manager := omdb.NewOMDbManager(mgr.cfg)
		omdbInfo := manager.GetOMDbInfoByTitle(title)
		if omdbInfo == nil {
			return nil
		}
		return &movieinfo.MovieInfo{
			Rated:   omdbInfo.Rated,
			Country: omdbInfo.Country,
		}
	}

	currentYear := strconv.Itoa(time.Now().Year())
	movieInfos := getIMDbInfosByYear(mgr.cfg, currentYear, getOMDbInfoByTitle)
	newMovies := filterAlreadyReturnedMovies(movieInfos)

	for _, movie := range newMovies {
		alreadyReturnedMovies[movie.Title] = true
	}

	if err := saveReturnedMovies(); err != nil {
		log.Println("Error saving already returned movies:", err)
	}

	return buildMovieInfoStrings(newMovies)
}

func filterAlreadyReturnedMovies(imdbRatingsMovies []movieinfo.MovieInfo) []movieinfo.MovieInfo {
	filteredMovies := make([]movieinfo.MovieInfo, 0, len(imdbRatingsMovies))
	for _, movie := range imdbRatingsMovies {
		if alreadyReturnedMovies[movie.Title] {
			continue
		}
		if isIndianCountry(movie.Country) {
			continue
		}
		if !movieIsOlderThanOneMonth(movie) {
			continue
		}
		filteredMovies = append(filteredMovies, movie)
	}
	return filteredMovies
}

func isIndianCountry(country string) bool {
	country = strings.TrimSpace(strings.ToLower(country))
	if country == "" {
		return false
	}

	separators := []string{",", ";", "/", "|"}
	for _, sep := range separators {
		country = strings.ReplaceAll(country, sep, ",")
	}

	parts := strings.Split(country, ",")
	for _, part := range parts {
		if strings.TrimSpace(part) == "india" {
			return true
		}
	}
	return false
}

func movieIsOlderThanOneMonth(movie movieinfo.MovieInfo) bool {
	if movie.ReleaseDate == "" {
		return false
	}

	releaseDate, err := parseReleaseDate(movie.ReleaseDate)
	if err != nil {
		log.Println("Error parsing release date for movie", movie.Title, err)
		return false
	}

	now := time.Now()
	return !releaseDate.AddDate(0, 1, 0).After(now)
}

func parseReleaseDate(dateStr string) (time.Time, error) {
	formats := []string{"02.01.06", "02.01.2006"}
	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}

func saveReturnedMovies() error {
	movieTitles := make([]string, 0, len(alreadyReturnedMovies))
	for title := range alreadyReturnedMovies {
		movieTitles = append(movieTitles, title)
	}

	jsonData, err := json.MarshalIndent(movieTitles, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(alreadyReturnedMoviesFile, jsonData, 0644)
}

func loadAlreadyReturnedMovies() error {
	jsonData, err := os.ReadFile(alreadyReturnedMoviesFile)
	if err != nil {
		if os.IsNotExist(err) {
			if err := createEmptyJsonFile(alreadyReturnedMoviesFile); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	var movieTitles []string
	if err := json.Unmarshal(jsonData, &movieTitles); err != nil {
		return err
	}

	for _, title := range movieTitles {
		alreadyReturnedMovies[title] = true
	}
	return nil
}

func createEmptyJsonFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	emptyJson := "[]"
	if _, err := file.WriteString(emptyJson); err != nil {
		return err
	}

	return nil
}
