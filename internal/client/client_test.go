package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetReturnsErrorOnNon2xxStatus(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{
			Timeout: time.Second,
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadGateway,
					Body:       io.NopCloser(strings.NewReader("bad gateway")),
					Header:     make(http.Header),
					Request:    req,
				}, nil
			}),
		},
	}

	_, err := c.Get("https://example.com", url.Values{"q": {"test"}})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
	if !strings.Contains(err.Error(), "unexpected status 502") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearchQuotesFreshBypassesCache(t *testing.T) {
	originalClient := DefaultClient
	originalCache := searchCache
	defer func() {
		DefaultClient = originalClient
		searchCache = originalCache
	}()

	searchCache = map[string][]*Quote{
		"pingan": {
			{Code: "cached", Name: "cached"},
		},
	}

	DefaultClient = &Client{
		httpClient: &http.Client{
			Timeout: time.Second,
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				payload := map[string]any{
					"QuotationCodeTable": map[string]any{
						"Data": []map[string]any{
							{"Code": "000001", "Name": "平安银行", "QuoteID": "0.000001", "MktNum": "0"},
						},
					},
				}
				body, _ := json.Marshal(payload)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(string(body))),
					Header:     make(http.Header),
					Request:    req,
				}, nil
			}),
		},
	}

	resp, err := SearchQuotesFresh("pingan", 1)
	if err != nil {
		t.Fatalf("SearchQuotesFresh returned error: %v", err)
	}
	if resp == nil || len(resp.Candidates) != 1 {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if resp.Candidates[0].Code != "000001" {
		t.Fatalf("expected fresh response, got %#v", resp.Candidates[0])
	}
	if len(resp.Raw) == 0 {
		t.Fatal("expected raw payload from fresh search")
	}
}

func TestMain(m *testing.M) {
	searchCache = make(map[string][]*Quote)
	os.Exit(m.Run())
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
