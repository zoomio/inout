package inout

import (
	"context"
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
	req      *http.Request
	resp     *http.Response
	err      error
	attempts int
}

func fetch(ctx context.Context, src string) *fetchResult {
	var req *http.Request
	var err error
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
	if err != nil {
		return &fetchResult{err: fmt.Errorf("error in creating request for source=%s: %w", src, err)}
	}

	var resp *http.Response
	var doRetry bool
	var attempts int
	backOff := attemptBackOff
	for {
		attempts++
		resp, doRetry, err = fetchReq(req)
		if !doRetry {
			return &fetchResult{
				req:      req,
				resp:     resp,
				err:      err,
				attempts: attempts,
			}
		}
		if attempts >= maxAttempts {
			return &fetchResult{
				req:      req,
				resp:     resp,
				err:      fmt.Errorf("exceeded max amount of attempts [%d]", maxAttempts),
				attempts: attempts,
			}
		}
		time.Sleep(backOff)
		backOff *= time.Duration(math.Max(7, 15*rnd.Float64()))
	}

}

func fetchReq(req *http.Request) (resp *http.Response, doRetry bool, err error) {
	resp, err = webClient.Do(req)
	if err != nil {
		err = fmt.Errorf("error in calling GET on provided source=%s: %w", req.URL.String(), err)
		return
	}
	if resp.StatusCode == 429 || resp.StatusCode >= 500 {
		doRetry = true
	}
	return
}
