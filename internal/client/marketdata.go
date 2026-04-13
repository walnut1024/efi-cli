package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	clistAPIURL   = "https://push2.eastmoney.com/api/qt/clist/get"
	ulistAPIURL   = "https://push2.eastmoney.com/api/qt/ulist.np/get"
	detailAPIURL  = "https://push2.eastmoney.com/api/qt/stock/details/get"
	defaultPageSz = 200
)

type QuoteListResponse struct {
	Total int
	Diff  []map[string]interface{}
}

type DealDetailResponse struct {
	Details  []string
	PrePrice interface{}
}

func FetchQuoteList(fs, fields string, pn, pz int) (*QuoteListResponse, error) {
	resp, _, err := FetchQuoteListRaw(fs, fields, pn, pz)
	return resp, err
}

func FetchQuoteListRaw(fs, fields string, pn, pz int) (*QuoteListResponse, []byte, error) {
	if pz <= 0 {
		pz = defaultPageSz
	}
	params := url.Values{
		"pn":     {fmt.Sprintf("%d", pn)},
		"pz":     {fmt.Sprintf("%d", pz)},
		"po":     {"1"},
		"np":     {"1"},
		"fltt":   {"2"},
		"invt":   {"2"},
		"fid":    {"f12"},
		"fs":     {fs},
		"fields": {fields},
	}
	data, err := DefaultClient.Get(clistAPIURL, params)
	if err != nil {
		return nil, nil, err
	}

	var resp struct {
		Data struct {
			Total int                      `json:"total"`
			Diff  []map[string]interface{} `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, data, err
	}
	return &QuoteListResponse{Total: resp.Data.Total, Diff: resp.Data.Diff}, data, nil
}

func FetchQuotesBySecIDs(secids, fields string) ([]map[string]interface{}, error) {
	items, _, err := FetchQuotesBySecIDsRaw(secids, fields)
	return items, err
}

func FetchQuotesBySecIDsRaw(secids, fields string) ([]map[string]interface{}, []byte, error) {
	params := url.Values{
		"fields": {fields},
		"fltt":   {"2"},
		"secids": {secids},
	}
	data, err := DefaultClient.Get(ulistAPIURL, params)
	if err != nil {
		return nil, nil, err
	}

	var resp struct {
		Data struct {
			Diff []map[string]interface{} `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, data, err
	}
	return resp.Data.Diff, data, nil
}

func FetchDealDetails(secid string, maxCount int) (*DealDetailResponse, error) {
	if maxCount <= 0 {
		maxCount = 10000
	}

	params := url.Values{
		"secid":   {secid},
		"fields1": {"f1,f2,f3,f4,f5"},
		"fields2": {"f51,f52,f53,f54,f55"},
		"pos":     {fmt.Sprintf("-%d", maxCount)},
	}
	data, err := DefaultClient.Get(detailAPIURL, params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			Details  []string    `json:"details"`
			PrePrice interface{} `json:"prePrice"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &DealDetailResponse{Details: resp.Data.Details, PrePrice: resp.Data.PrePrice}, nil
}
