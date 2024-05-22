package telegram_bot

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	movieinfos "nightmare_navigator/internal/movie_infos"
)

type MockHTTPClient struct{}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var data string
	if strings.Contains(req.URL.String(), movieinfos.BasicsFilename) {
		data = "tconst\ttitleType\tprimaryTitle\toriginalTitle\tisAdult\tstartYear\tendYear\truntimeMinutes\tgenres\n" +
			"tt0119081\tmovie\tEvent Horizon\tEvent Horizon\t0\t1997\t\\N\t96\tHorror,Sci-Fi\n" +
			"tt0391198\tmovie\tThe Grudge\tThe Grudge\t0\t2004\t\\N\t92\tHorror,Mystery,Thriller\n"
	} else if strings.Contains(req.URL.String(), movieinfos.RatingsFilename) {
		data = "tconst\taverageRating\tnumVotes\n" +
			"tt0119081\t6.7\t120000\n" +
			"tt0391198\t5.9\t90000\n"
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(io.Reader(bytes.NewBufferString(data))),
	}, nil
}

func TestRunTelegramBotWithInvalidHTTPClient(t *testing.T) {
	// Invalid Mock HTTP client
	http.DefaultClient = nil // Set to nil to force error

	// Run Telegram Bot
	RunTelegramBot()
}

func TestRunTelegramBotWithInvalidTimer(t *testing.T) {
	// Invalid Mock Timer
	timer := time.NewTimer(0 * time.Second)

	// Run Telegram Bot
	go RunTelegramBot()

	select {
	case <-timer.C:
	}
}
