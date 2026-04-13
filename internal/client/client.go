package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	UserAgent       = "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; Touch; rv:11.0) like Gecko"
	AcceptHeader    = "*/*"
	AcceptLanguage  = "zh-CN,zh;q=0.8,zh-TW;q=0.7,en-US;q=0.3,en;q=0.2"
	MaxRetries      = 3
	MaxResponseBody = 32 * 1024 * 1024 // 32 MB
	RetryDelay      = 1 * time.Second
	MaxConcurrency  = 10
	RequestTimeout  = 10 * time.Second
)

type Client struct {
	httpClient *http.Client
}

var DefaultClient = &Client{
	httpClient: &http.Client{
		Timeout: RequestTimeout,
	},
}

func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

func (c *Client) Get(rawURL string, params url.Values) ([]byte, error) {
	return c.GetWithHeaders(rawURL, params, nil)
}

func (c *Client) GetWithHeaders(rawURL string, params url.Values, headers map[string]string) ([]byte, error) {
	reqURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if params != nil {
		q := reqURL.Query()
		for k, vs := range params {
			q[k] = vs
		}
		reqURL.RawQuery = q.Encode()
	}

	var lastErr error
	for i := 0; i < MaxRetries; i++ {
		if i > 0 {
			time.Sleep(RetryDelay)
		}
		req, err := http.NewRequest("GET", reqURL.String(), nil)
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("User-Agent", UserAgent)
		req.Header.Set("Accept", AcceptHeader)
		req.Header.Set("Accept-Language", AcceptLanguage)
		if headers != nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
		req.Header.Set("Referer", "https://quote.eastmoney.com/")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, MaxResponseBody))
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			lastErr = fmt.Errorf("unexpected status %d: %s", resp.StatusCode, truncateBody(body))
			continue
		}
		return body, nil
	}
	return nil, fmt.Errorf("after %d retries: %w", MaxRetries, lastErr)
}

func (c *Client) PostJSON(rawURL string, jsonData interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}

	var lastErr error
	for i := 0; i < MaxRetries; i++ {
		if i > 0 {
			time.Sleep(RetryDelay)
		}
		req, err := http.NewRequest("POST", rawURL, bytes.NewReader(jsonBytes))
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", UserAgent)
		req.Header.Set("Referer", "https://quote.eastmoney.com/")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, MaxResponseBody))
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			lastErr = fmt.Errorf("unexpected status %d: %s", resp.StatusCode, truncateBody(body))
			continue
		}
		return body, nil
	}
	return nil, fmt.Errorf("after %d retries: %w", MaxRetries, lastErr)
}

func (c *Client) FetchPages(total int, pageSize int, fetchFn func(page int) ([]byte, error)) ([][]byte, error) {
	pages := total / pageSize
	if total%pageSize != 0 {
		pages++
	}
	if pages < 1 {
		pages = 1
	}

	results := make([][]byte, pages)
	type result struct {
		idx  int
		data []byte
		err  error
	}

	sem := make(chan struct{}, MaxConcurrency)
	ch := make(chan result, pages)

	for i := 0; i < pages; i++ {
		go func(idx int) {
			sem <- struct{}{}
			data, err := fetchFn(idx + 1)
			<-sem
			ch <- result{idx: idx, data: data, err: err}
		}(i)
	}

	for i := 0; i < pages; i++ {
		r := <-ch
		if r.err != nil {
			return nil, fmt.Errorf("page %d: %w", r.idx+1, r.err)
		}
		results[r.idx] = r.data
	}
	return results, nil
}

func truncateBody(body []byte) string {
	const max = 256
	if len(body) <= max {
		return string(body)
	}
	return string(body[:max]) + "..."
}
