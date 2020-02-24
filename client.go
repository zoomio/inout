package inout

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const (
	maxAttempts    = 3
	attemptBackOff = time.Millisecond * 30
)

var (
	webClient = &http.Client{
		Timeout: time.Second * 3,
	}

	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type fetchResult struct {
	resp     *http.Response
	err      error
	attempts int
}

func fetch(src string) *fetchResult {
	var err error
	var resp *http.Response
	var doRetry bool
	var attempts int
	backOff := attemptBackOff
	for {
		attempts++
		resp, doRetry, err = fetchSrc(src)
		if !doRetry {
			return &fetchResult{
				resp:     resp,
				err:      err,
				attempts: attempts,
			}
		}
		if attempts >= maxAttempts {
			return &fetchResult{
				resp:     resp,
				err:      fmt.Errorf("exceeded max amount of attempts [%d]", maxAttempts),
				attempts: attempts,
			}
		}
		time.Sleep(backOff)
		backOff *= time.Duration(math.Max(7, 15*rnd.Float64()))
	}

}

func fetchSrc(src string) (resp *http.Response, doRetry bool, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, src, nil)
	if err != nil {
		err = fmt.Errorf("error in creating request for source=%s: %w", src, err)
		return
	}
	resp, err = webClient.Do(req)
	if err != nil {
		err = fmt.Errorf("error in calling provided source=%s: %w", src, err)
		return
	}
	if resp.StatusCode == 429 || resp.StatusCode >= 500 {
		doRetry = true
	}
	return
}
