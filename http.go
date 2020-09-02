package mango

import (
	"fmt"
	"net/http"
	"time"
)

type config struct {
	count    int
	pause    time.Duration
	statuses httpStatusList
}

type httpStatusList []int

func (s httpStatusList) Consist(status int) bool {
	for _, item := range s {
		if item == status {
			return true
		}
	}

	return false
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpClientRetryWrap struct {
	config config
	client httpClient
}

func NewDefaultHTTPClientRetryWrap(client httpClient) *httpClientRetryWrap {
	var (
		httpRetryDefaultCount    = 3
		httpRetryDefaultPause    = time.Second
		httpRetryDefaultStatuses = []int{524} //could consist custom statuses e.g. 524: A timeout occurred from CloudFlare
	)

	return &httpClientRetryWrap{
		config: config{
			count:    httpRetryDefaultCount,
			pause:    httpRetryDefaultPause,
			statuses: httpRetryDefaultStatuses,
		},
		client: client,
	}
}

func (w *httpClientRetryWrap) do(req *http.Request) (*http.Response, error) {
	for i := w.config.count; i > 0; i-- {
		resp, err := w.client.Do(req)
		if err != nil {
			return nil, err
		}

		if !w.config.statuses.Consist(resp.StatusCode) {
			return resp, nil
		}

		time.Sleep(w.config.pause)
	}

	return nil, fmt.Errorf("request retrying failed for URL: %s", req.URL)
}
