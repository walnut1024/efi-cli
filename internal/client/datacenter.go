package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
)

const datacenterBaseURL = "https://datacenter-web.eastmoney.com/api/data/v1/get"

// DatacenterResponse holds the parsed response from datacenter-web API.
type DatacenterResponse struct {
	Data  []map[string]interface{}
	Count int
	Pages int
}

// FetchDatacenterReport fetches a single page from datacenter-web API.
func FetchDatacenterReport(reportName string, filters url.Values, page, pageSize int) (*DatacenterResponse, error) {
	params := url.Values{
		"reportName": {reportName},
		"columns":    {"ALL"},
		"pageSize":   {fmt.Sprintf("%d", pageSize)},
		"pageNo":     {fmt.Sprintf("%d", page)},
	}
	if filters != nil {
		for k, vs := range filters {
			params[k] = vs
		}
	}

	data, err := DefaultClient.Get(datacenterBaseURL, params)
	if err != nil {
		return nil, err
	}

	return parseDatacenterResponse(data)
}

// FetchDatacenterReportAll fetches all pages until limit is reached or data exhausted.
// Pages 2+ are fetched concurrently using the shared MaxConcurrency semaphore.
func FetchDatacenterReportAll(reportName string, filters url.Values, limit int) ([]map[string]interface{}, error) {
	const pageSize = 500

	first, err := FetchDatacenterReport(reportName, filters, 1, pageSize)
	if err != nil {
		return nil, err
	}
	if len(first.Data) == 0 {
		return nil, nil
	}

	all := first.Data
	if limit > 0 && len(all) >= limit {
		return all[:limit], nil
	}
	if first.Pages <= 1 {
		return all, nil
	}

	remaining := first.Pages - 1
	type pageResult struct {
		data []map[string]interface{}
		err  error
	}
	results := make([]pageResult, remaining)

	var wg sync.WaitGroup
	sem := make(chan struct{}, MaxConcurrency)
	for p := 2; p <= first.Pages; p++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			resp, fetchErr := FetchDatacenterReport(reportName, filters, page, pageSize)
			if fetchErr != nil {
				results[page-2] = pageResult{err: fetchErr}
				return
			}
			results[page-2] = pageResult{data: resp.Data}
		}(p)
	}
	wg.Wait()

	for _, res := range results {
		if res.err != nil {
			return all, nil // return partial results on error, same as before
		}
		all = append(all, res.data...)
		if limit > 0 && len(all) >= limit {
			return all[:limit], nil
		}
	}
	return all, nil
}

func parseDatacenterResponse(data []byte) (*DatacenterResponse, error) {
	var resp struct {
		Result struct {
			Data  []map[string]interface{} `json:"data"`
			Count int                      `json:"count"`
			Pages int                      `json:"pages"`
		} `json:"result"`
		Success bool   `json:"success"`
		Message string `json:"message"`
		Code    int    `json:"code"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("datacenter error (code %d): %s", resp.Code, resp.Message)
	}
	if resp.Result.Data == nil {
		return &DatacenterResponse{Data: []map[string]interface{}{}, Count: 0, Pages: 0}, nil
	}
	return &DatacenterResponse{
		Data:  resp.Result.Data,
		Count: resp.Result.Count,
		Pages: resp.Result.Pages,
	}, nil
}
