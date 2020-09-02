package mango

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestHTTPClientRetryWrap_do(t *testing.T) {
	testCases := []struct {
		name   string
		client httpClientRetryWrap
		errMsg string
	}{
		{
			name: "Retrying failed",
			client: httpClientRetryWrap{
				config: config{
					count:    2,
					pause:    time.Nanosecond,
					statuses: []int{http.StatusGatewayTimeout},
				},
				client: &httpClientMock{
					returnedResponse: []*http.Response{
						{
							StatusCode: http.StatusGatewayTimeout,
						},
						{
							StatusCode: http.StatusGatewayTimeout,
						},
					},
				},
			},
			errMsg: "request retrying failed",
		},
		{
			name: "Retry is not needed",
			client: httpClientRetryWrap{
				config: config{
					count:    2,
					pause:    time.Nanosecond,
					statuses: []int{http.StatusGatewayTimeout},
				},
				client: &httpClientMock{
					returnedResponse: []*http.Response{
						{
							StatusCode: http.StatusOK,
						},
					},
				},
			},
		},
		{
			name: "Error while request",
			client: httpClientRetryWrap{
				config: config{
					count:    2,
					pause:    time.Nanosecond,
					statuses: []int{http.StatusGatewayTimeout},
				},
				client: &httpClientMock{
					returnedError: errors.New("some error"),
				},
			},
			errMsg: "some error",
		},
		{
			name: "Success after few attempts",
			client: httpClientRetryWrap{
				config: config{
					count:    3,
					pause:    time.Nanosecond,
					statuses: []int{http.StatusGatewayTimeout},
				},
				client: &httpClientMock{
					returnedResponse: []*http.Response{
						{
							StatusCode: http.StatusGatewayTimeout,
						},
						{
							StatusCode: http.StatusGatewayTimeout,
						},
						{
							StatusCode: http.StatusOK,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://test.de", nil)
			if err != nil {
				t.Fatal(err)
			}

			_, err = tc.client.do(req)
			if tc.errMsg != "" {
				if err == nil {
					t.Fatal("should be an error")
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Fatalf("(%s) should contains string (%s)", err.Error(), tc.errMsg)
				}
			} else {
				if err != nil {
					t.Fatal("should not be an error")
				}
			}
		})
	}
}

type httpClientMock struct {
	returnedError    error
	returnedResponse []*http.Response
}

func (s *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	if s.returnedError != nil {
		return nil, s.returnedError
	}

	if len(s.returnedResponse) == 0 {
		return nil, errors.New("httpClientMock: nothing to return")
	}
	res := s.returnedResponse[0]
	s.returnedResponse = s.returnedResponse[1:]

	return res, nil
}
