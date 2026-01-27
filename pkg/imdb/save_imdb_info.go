package imdb

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"nightmare_navigator/internal/config"
	"nightmare_navigator/pkg/tmdb"
)

type TitleBasics struct {
	Tconst        string
	PrimaryTitle  string
	OriginalTitle string
	Genres        string
	StartYear     string
}

type TitleRatings struct {
	Tconst        string
	AverageRating float64
	NumVotes      int
}

type IMDbMovieInfo struct {
	AverageRating string `json:"averageRating"`
	Genres        string `json:"genres"`
	NumVotes      string `json:"numVotes"`
	OriginalTitle string `json:"originalTitle"`
	PrimaryTitle  string `json:"primaryTitle"`
	ReleaseDate   string `json:"releaseDate"`
	Runtime       int    `json:"runtime"`
	Description   string `json:"description"`
	Tconst        string `json:"tconst"`
}

type IMDbManager struct {
	cfg config.Config
}

func NewIMDbManager(cfg config.Config) *IMDbManager {
	return &IMDbManager{cfg: cfg}
}

func (mgr *IMDbManager) SaveLatestIMDbRatings() {
	if err := os.MkdirAll(mgr.cfg.General.DataDir, 0755); err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	mgr.cleanupFiles([]string{mgr.cfg.IMDb.BasicsFilename, mgr.cfg.IMDb.RatingsFilename, mgr.cfg.IMDb.JSONFilename})

	if err := mgr.downloadAndExtractFiles(); err != nil {
		log.Fatalf("Error downloading files: %v", err)
	}

	basicsFile, err := os.Open(filepath.Join(mgr.cfg.General.DataDir, mgr.cfg.IMDb.BasicsFilename))
	if err != nil {
		log.Fatalf("Error opening title.basics.tsv.gz: %v", err)
	}
	defer basicsFile.Close()

	ratingsFile, err := os.Open(filepath.Join(mgr.cfg.General.DataDir, mgr.cfg.IMDb.RatingsFilename))
	if err != nil {
		log.Fatalf("Error opening title.ratings.tsv.gz: %v", err)
	}
	defer ratingsFile.Close()

	outputFile, err := os.Create(filepath.Join(mgr.cfg.General.DataDir, mgr.cfg.IMDb.JSONFilename))
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	imdbMovies := mgr.loadIMDbMovies(basicsFile, ratingsFile)

	mgr.writeJSON(outputFile, imdbMovies)

	log.Println("Updated IMDb movie list")
}

func (mgr *IMDbManager) cleanupFiles(files []string) {
	for _, file := range files {
		path := filepath.Join(mgr.cfg.General.DataDir, file)
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
		}
	}
}

func (mgr *IMDbManager) downloadAndExtractFiles() error {
	files := []string{mgr.cfg.IMDb.BasicsFilename, mgr.cfg.IMDb.RatingsFilename}
	for _, file := range files {
		if err := downloadFile(mgr.cfg.IMDb.IMDbBaseUrl+file, filepath.Join(mgr.cfg.General.DataDir, file)); err != nil {
			return fmt.Errorf("error downloading %s: %v", file, err)
		}
	}
	return nil
}

func (mgr *IMDbManager) loadIMDbMovies(basicsFile, ratingsFile *os.File) map[string][]IMDbMovieInfo {
	ratings := loadRatings(ratingsFile)
	return mgr.loadBasicsAndFilter(basicsFile, ratings)
}

func loadRatings(file *os.File) map[string]*TitleRatings {
	ratings := make(map[string]*TitleRatings)
	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatalf("Error creating gzip reader for ratings file: %v", err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) >= 3 {
			averageRating, _ := strconv.ParseFloat(fields[1], 64)
			numVotes, _ := strconv.Atoi(fields[2])
			ratings[fields[0]] = &TitleRatings{
				Tconst:        fields[0],
				AverageRating: averageRating,
				NumVotes:      numVotes,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading ratings file: %v", err)
	}

	return ratings
}

func (mgr *IMDbManager) loadBasicsAndFilter(file *os.File, ratings map[string]*TitleRatings) map[string][]IMDbMovieInfo {
	moviesByYear := make(map[string][]IMDbMovieInfo)
	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatalf("Error creating gzip reader for basics file: %v", err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) >= 9 && fields[1] == "movie" {
			mgr.filterMovieAndGetAdditionalInfo(fields, ratings, moviesByYear)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading basics file: %v", err)
	}

	return moviesByYear
}

func (mgr *IMDbManager) filterMovieAndGetAdditionalInfo(fields []string, ratings map[string]*TitleRatings, moviesByYear map[string][]IMDbMovieInfo) {
	if rating, ok := ratings[fields[0]]; ok && rating.AverageRating >= mgr.cfg.IMDb.MinRating && rating.NumVotes >= mgr.cfg.IMDb.MinVotes {
		genres := strings.Split(fields[8], ",")
		if containsGenre(genres, "Horror") && !containsGenre(genres, "Romance") && !containsGenre(genres, "Family") {
			movieInfo := mgr.createMovieInfo(fields, rating)
			startYear := fields[5]
			mgr.addAdditionalInfo(startYear, movieInfo)

			if movieInfo == nil || movieInfo.ReleaseDate == "" {
				return
			}

			moviesByYear[startYear] = append(moviesByYear[startYear], *movieInfo)
		}
	}
}

func (mgr *IMDbManager) createMovieInfo(fields []string, rating *TitleRatings) *IMDbMovieInfo {
	genres := strings.Split(fields[8], ",")
	var genresFormatted strings.Builder

	for i, genre := range genres {
		if i > 0 && i < len(genres)-1 {
			genresFormatted.WriteString(", ")
		}
		if i == len(genres)-1 && i != 0 {
			genresFormatted.WriteString(" and ")
		}
		genresFormatted.WriteString(genre)
	}

	return &IMDbMovieInfo{
		Tconst:        fields[0],
		PrimaryTitle:  fields[2],
		OriginalTitle: fields[3],
		Genres:        genresFormatted.String(),
		AverageRating: fmt.Sprintf("%.1f", rating.AverageRating),
		NumVotes:      strconv.Itoa(rating.NumVotes),
	}
}

func (mgr *IMDbManager) addAdditionalInfo(year string, imdbInfo *IMDbMovieInfo) {
	tmdbManager := tmdb.NewTMDbManager(mgr.cfg)
	tmdbInfo := tmdbManager.GetTMDbInfoByTitle(imdbInfo.PrimaryTitle, year)

	if tmdbInfo != nil {
		imdbInfo.Description = tmdbInfo.Overview
		imdbInfo.Runtime = tmdbInfo.Runtime
		imdbInfo.ReleaseDate = tmdbInfo.ReleaseDate
	}
}

func containsGenre(genres []string, genre string) bool {
	for _, g := range genres {
		if g == genre {
			return true
		}
	}
	return false
}

func (mgr *IMDbManager) writeJSON(outputFile *os.File, moviesByYear map[string][]IMDbMovieInfo) {
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "\t")

	years := sortedYears(moviesByYear)

	var result []map[string]interface{}
	for _, year := range years {
		movies := moviesByYear[year]
		sortIMDbMoviesByReleaseDate(movies)
		yearData := map[string]interface{}{
			"startYear": year,
			"data":      movies,
		}
		result = append(result, yearData)
	}

	if err := encoder.Encode(result); err != nil {
		log.Fatalf("Error encoding movie info: %v", err)
	}
}

func sortedYears(moviesByYear map[string][]IMDbMovieInfo) []string {
	years := make([]string, 0, len(moviesByYear))
	for year := range moviesByYear {
		years = append(years, year)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(years)))
	return years
}

func sortIMDbMoviesByReleaseDate(movies []IMDbMovieInfo) {
	sort.Slice(movies, func(i, j int) bool {
		releaseDateI, errI := time.Parse("02.01.06", movies[i].ReleaseDate)
		releaseDateJ, errJ := time.Parse("02.01.06", movies[j].ReleaseDate)
		if errI != nil || errJ != nil {
			return false
		}
		return releaseDateI.Before(releaseDateJ)
	})
}

func downloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
