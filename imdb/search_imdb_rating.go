package imdb

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type IMDbMovie struct {
	Title     string
	TitleId   string
	IMDb      string
	IMDbVotes string
}

func GetIMDbInfoByTitle(title string) *IMDbMovie {
	var movie IMDbMovie

	// Gets titleId in title.basics.tsv.gz file
	akasFilePath := filepath.Join(downloadDir, basicsFilename)

	file, err := os.Open(akasFilePath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		log.Println(err)
	}
	defer gzipReader.Close()

	scanner := bufio.NewScanner(gzipReader)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) >= 3 && fields[2] == title {
			movie = IMDbMovie{
				Title:   fields[2],
				TitleId: fields[0],
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	// Gets IMDb rating and votes in title.ratings.tsv.gz file
	ratingsFilePath := filepath.Join(downloadDir, ratingsFilename)
	file, err = os.Open(ratingsFilePath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	gzipReader, err = gzip.NewReader(file)
	if err != nil {
		log.Println(err)
	}
	defer gzipReader.Close()

	scanner = bufio.NewScanner(gzipReader)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) >= 2 && fields[0] == movie.TitleId {
			movie.IMDb = fields[1]
			movie.IMDbVotes = fields[2]
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	return &movie
}
