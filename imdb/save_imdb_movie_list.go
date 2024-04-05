package imdb

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	baseURL         = "https://datasets.imdbws.com/"
	basicsFilename  = "title.basics.tsv.gz"
	ratingsFilename = "title.ratings.tsv.gz"
	downloadDir     = "imdb"
)

type TitleBasics struct {
	Tconst        string `json:"tconst"`
	PrimaryTitle  string `json:"primaryTitle"`
	OriginalTitle string `json:"originalTitle"`
	Genres        string `json:"genres"`
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
	Genres        string `json:"genres"`
	AverageRating string `json:"averageRating"`
	NumVotes      string `json:"numVotes"`
}

func SaveLatestIMDbRatings() {
	// Delete existing files if they exist
	basicsFilePath := filepath.Join(downloadDir, basicsFilename)
	// if _, err := os.Stat(basicsFilePath); err == nil {
	// 	os.Remove(basicsFilePath)
	// }

	ratingsFilePath := filepath.Join(downloadDir, ratingsFilename)
	// if _, err := os.Stat(ratingsFilePath); err == nil {
	// 	os.Remove(ratingsFilePath)
	// }

	// // Download title.basics.tsv.gz
	// basicsURL := baseURL + basicsFilename
	// if err := downloadFile(basicsURL, basicsFilePath); err != nil {
	// 	log.Fatal("Error downloading title.basics.tsv.gz:", err)
	// }
	// log.Println("Downloaded", basicsFilename)

	// // Download title.ratings.tsv.gz
	// ratingsURL := baseURL + ratingsFilename
	// if err := downloadFile(ratingsURL, ratingsFilePath); err != nil {
	// 	log.Fatal("Error downloading title.ratings.tsv.gz:", err)
	// }
	// log.Println("Downloaded", ratingsFilename)

	// Open title.basics.tsv.gz
	basicsFile, err := os.Open(basicsFilePath)
	if err != nil {
		log.Fatal("Error opening title.basics.tsv.gz:", err)
	}
	defer basicsFile.Close()

	// Open title.ratings.tsv.gz
	ratingsFile, err := os.Open(ratingsFilePath)
	if err != nil {
		log.Fatal("Error opening title.ratings.tsv.gz:", err)
	}
	defer ratingsFile.Close()

	outputFilePath := "horror_movies.json"
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatal("Error creating output file:", err)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)

	ratings := make(map[string]*TitleRatings)

	// Read title.ratings.tsv.gz
	ratingsScanner := bufio.NewScanner(ratingsFile)
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

	// Read title.basics.tsv.gz and filter horror movies with average rating >= 5
	basicsScanner := bufio.NewScanner(basicsFile)
	for basicsScanner.Scan() {
		line := basicsScanner.Text()
		fields := strings.Split(line, "\t")
		fmt.Println("aha", len(fields))
		fmt.Println("aha2", fields)
		if len(fields) >= 9 && fields[1] == "movie" {
			for _, genre := range strings.Split(fields[8], ",") {
				if strings.Contains(genre, "Horror") {
					if rating, ok := ratings[fields[0]]; ok && rating.AverageRating >= "5" {
						movieInfo := MovieInfo{
							Tconst:        fields[0],
							PrimaryTitle:  fields[2],
							OriginalTitle: fields[3],
							Genres:        fields[8],
							AverageRating: rating.AverageRating,
							NumVotes:      rating.NumVotes,
						}
						err := encoder.Encode(movieInfo)
						if err != nil {
							log.Println("Error encoding movie info:", err)
						}
						break
					}
				}
			}
		}
	}
	if err := basicsScanner.Err(); err != nil {
		log.Fatal("Error scanning title.basics.tsv.gz:", err)
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
