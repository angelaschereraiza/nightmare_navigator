package imdb

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	baseURL         = "https://datasets.imdbws.com/"
	basicsFilename  = "title.basics.tsv.gz"
	ratingsFilename = "title.ratings.tsv.gz"
	jsonFilename    = "imdb_ratings.json"
	downloadDir     = "data"
)

type TitleBasics struct {
	Tconst        string `json:"tconst"`
	PrimaryTitle  string `json:"primaryTitle"`
	OriginalTitle string `json:"originalTitle"`
	Genres        string `json:"genres"`
	StartYear     string `json:"startYear"`
}

type TitleRatings struct {
	Tconst        string `json:"tconst"`
	AverageRating string `json:"averageRating"`
	NumVotes      string `json:"numVotes"`
}

type MovieInfo struct {
	Tconst        string `json:"tconst"`
	PrimaryTitle  string `json:"primaryTitle"`
	OriginalTitle string `json:"originalTitle"`
	Year          string `json:"startYear"`
	Genres        string `json:"genres"`
	AverageRating string `json:"averageRating"`
	NumVotes      string `json:"numVotes"`
}

func SaveLatestIMDbRatings() {
	// Create the directory if it does not exist
	err := os.MkdirAll(downloadDir, 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	// Delete existing files if they exist
	basicsFilePath := filepath.Join(downloadDir, basicsFilename)
	if _, err := os.Stat(basicsFilePath); err == nil {
		os.Remove(basicsFilePath)
	}

	ratingsFilePath := filepath.Join(downloadDir, ratingsFilename)
	if _, err := os.Stat(ratingsFilePath); err == nil {
		os.Remove(ratingsFilePath)
	}

	jsonFilePath := filepath.Join(downloadDir, jsonFilename)
	if _, err := os.Stat(jsonFilePath); err == nil {
		os.Remove(jsonFilePath)
	}

	// // Download files
	if err := downloadFile(baseURL+basicsFilename, basicsFilePath); err != nil {
		log.Fatal("Error downloading title.basics.tsv.gz:", err)
	}

	if err := downloadFile(baseURL+ratingsFilename, ratingsFilePath); err != nil {
		log.Fatal("Error downloading title.ratings.tsv.gz:", err)
	}

	// Open files
	basicsFile, err := os.Open(basicsFilePath)
	if err != nil {
		log.Fatal("Error opening title.basics.tsv.gz:", err)
	}
	defer basicsFile.Close()

	ratingsFile, err := os.Open(ratingsFilePath)
	if err != nil {
		log.Fatal("Error opening title.ratings.tsv.gz:", err)
	}
	defer ratingsFile.Close()

	// Creates a json file to save the data
	outputFile, err := os.Create(downloadDir + "/" + jsonFilename)
	if err != nil {
		log.Fatal("Error creating output file:", err)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "\t")

	var movies []MovieInfo

	ratings := make(map[string]*TitleRatings)

	// Create a gzip.Reader to read the files
	gzipBasicsFileReader, err := gzip.NewReader(basicsFile)
	if err != nil {
		log.Fatal("Error creating gzip reader:", err)
	}
	defer gzipBasicsFileReader.Close()

	gzipRatingsFileReader, err := gzip.NewReader(ratingsFile)
	if err != nil {
		log.Fatal("Error creating gzip reader:", err)
	}
	defer gzipRatingsFileReader.Close()

	// Read title.ratings.tsv.gz
	ratingsScanner := bufio.NewScanner(gzipRatingsFileReader)
	for ratingsScanner.Scan() {
		line := ratingsScanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) >= 3 {
			ratings[fields[0]] = &TitleRatings{
				Tconst:        fields[0],
				AverageRating: fields[1],
				NumVotes:      fields[2],
			}
		}
	}

	// Read title.basics.tsv.gz and filter horror movies with average rating >= 5 and the count of votes >= 1000
	basicsScanner := bufio.NewScanner(gzipBasicsFileReader)
	for basicsScanner.Scan() {
		line := basicsScanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) >= 9 && fields[1] == "movie" {
			for _, genre := range strings.Split(fields[8], ",") {
				if strings.Contains(genre, "Horror") && !strings.Contains(genre, "Romance") && !strings.Contains(genre, "Family") {
					if rating, ok := ratings[fields[0]]; ok && rating.AverageRating >= "5" {
						numVotes, err := strconv.Atoi(rating.NumVotes)
						if err != nil {
							log.Printf("Error converting NumVotes to int: %v", err)
							continue
						}
						if numVotes >= 1000 {
							movieInfo := MovieInfo{
								Tconst:        fields[0],
								PrimaryTitle:  fields[2],
								OriginalTitle: fields[3],
								Year:          fields[5],
								Genres:        fields[8],
								AverageRating: rating.AverageRating,
								NumVotes:      rating.NumVotes,
							}
							movies = append(movies, movieInfo)
							break
						}
					}
				}
			}
		}
	}
	if err := basicsScanner.Err(); err != nil {
		log.Fatal("Error scanning title.basics.tsv.gz:", err)
	}

	err = encoder.Encode(movies)
	if err != nil {
		log.Fatal("Error encoding movie info:", err)
	}

	log.Println("Updated imdb rating list")
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
	if err != nil {
		return err
	}

	return nil
}
