package movie_info

import (
	"reflect"
	"strings"
	"testing"
)

func TestBuildMovieInfoString(t *testing.T) {
	tests := []struct {
		movieInfo MovieInfo
		expected  string
	}{
		{
			movieInfo: MovieInfo{
				Title:         "Event Horizon",
				OriginalTitle: "",
				IMDb:          "6.7",
				IMDbVotes:     "250,000",
				TitleId:       "tt0119081",
				Country:       "USA",
				Rated:         "R",
				Genres:        "Horror, Sci-Fi",
				ReleaseDate:   "15 Aug 1997",
				Runtime:       96,
				Description:   "A rescue crew investigates a spaceship that disappeared into a black hole and has now returned...with someone or something new on-board.",
			},
			expected: `Title: Event Horizon
				IMDb Rating: 6.7
				IMDb Votes: 250,000
				IMDb Link: https://www.imdb.com/title/tt0119081
				Country: USA
				Rated: R
				Genres: Horror, Sci-Fi
				Released: 15 Aug 1997
				Runtime: 96 minutes
				Description: A rescue crew investigates a spaceship that disappeared into a black hole and has now returned...with someone or something new on-board.
			`,
		},
		{
			movieInfo: MovieInfo{
				Title:         "The Grudge",
				OriginalTitle: "呪怨",
				IMDb:          "4.3",
				IMDbVotes:     "10,000",
				TitleId:       "tt3612126",
				Country:       "Japan",
				Rated:         "R",
				Genres:        "Horror, Mystery",
				ReleaseDate:   "03 Jan 2020",
				Runtime:       94,
				Description:   "A house is cursed by a vengeful ghost that dooms those who enter it with a violent death.",
			},
			expected: `Title: The Grudge
				Original Title: 呪怨
				IMDb Rating: 4.3
				IMDb Votes: 10,000
				IMDb Link: https://www.imdb.com/title/tt3612126
				Country: Japan
				Rated: R
				Genres: Horror, Mystery
				Released: 03 Jan 2020
				Runtime: 94 minutes
				Description: A house is cursed by a vengeful ghost that dooms those who enter it with a violent death.
			`,
		},
	}

	for _, test := range tests {
		result := buildMovieInfoString(test.movieInfo)
		result = strings.TrimSpace(result)
		expected := strings.TrimSpace(test.expected)

		if result != expected {
			t.Errorf("BuildMovieInfoString() = %q; want %q", result, expected)
		}
	}
}

func TestBuildMovieInfoStrings(t *testing.T) {
	movieInfos := []MovieInfo{
		{
			Title:         "Event Horizon",
			OriginalTitle: "",
			IMDb:          "6.7",
			IMDbVotes:     "250,000",
			TitleId:       "tt0119081",
			Country:       "USA",
			Rated:         "R",
			Genres:        "Horror, Sci-Fi",
			ReleaseDate:   "15 Aug 1997",
			Runtime:       96,
			Description:   "A rescue crew investigates a spaceship that disappeared into a black hole and has now returned...with someone or something new on-board.",
		},
		{
			Title:         "The Grudge",
			OriginalTitle: "呪怨",
			IMDb:          "4.3",
			IMDbVotes:     "10,000",
			TitleId:       "tt3612126",
			Country:       "Japan",
			Rated:         "R",
			Genres:        "Horror, Mystery",
			ReleaseDate:   "03 Jan 2020",
			Runtime:       94,
			Description:   "A house is cursed by a vengeful ghost that dooms those who enter it with a violent death.",
		},
	}

	expected := []string{
		`Title: Event Horizon
		IMDb Rating: 6.7
		IMDb Votes: 250,000
		IMDb Link: https://www.imdb.com/title/tt0119081
		Country: USA
		Rated: R
		Genres: Horror, Sci-Fi
		Released: 15 Aug 1997
		Runtime: 96 minutes
		Description: A rescue crew investigates a spaceship that disappeared into a black hole and has now returned...with someone or something new on-board.
		`,
				`Title: The Grudge
		Original Title: 呪怨
		IMDb Rating: 4.3
		IMDb Votes: 10,000
		IMDb Link: https://www.imdb.com/title/tt3612126
		Country: Japan
		Rated: R
		Genres: Horror, Mystery
		Released: 03 Jan 2020
		Runtime: 94 minutes
		Description: A house is cursed by a vengeful ghost that dooms those who enter it with a violent death.
	`,
	}

	result := BuildMovieInfoStrings(movieInfos)
	for i, r := range *result {
		(*result)[i] = strings.TrimSpace(r)
	}
	for i, e := range expected {
		expected[i] = strings.TrimSpace(e)
	}

	if !reflect.DeepEqual(*result, expected) {
		t.Errorf("BuildMovieInfoStrings() = %v; want %v", *result, expected)
	}
}
