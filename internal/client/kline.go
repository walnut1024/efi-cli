package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

const (
	klineAPIURL     = "https://push2his.eastmoney.com/api/qt/stock/kline/get"
	defaultBegDate  = "19000101"
	defaultEndDate  = "20500101"
	defaultKlineKlt = 101
	defaultKlineFqt = 1
)

type KlineResponse struct {
	Name   string
	Klines []string
}

func FetchKlines(secid, beg, end string, klt, fqt int) (*KlineResponse, error) {
	resp, _, err := FetchKlinesRaw(secid, beg, end, klt, fqt)
	return resp, err
}

func FetchKlinesRaw(secid, beg, end string, klt, fqt int) (*KlineResponse, []byte, error) {
	if beg == "" {
		beg = defaultBegDate
	}
	if end == "" {
		end = defaultEndDate
	}
	if klt == 0 {
		klt = defaultKlineKlt
	}
	if fqt == 0 {
		fqt = defaultKlineFqt
	}

	params := url.Values{
		"fields1": {"f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13"},
		"fields2": {strings.Join([]string{"f51", "f52", "f53", "f54", "f55", "f56", "f57", "f58", "f59", "f60", "f61"}, ",")},
		"beg":     {beg},
		"end":     {end},
		"rtntype": {"6"},
		"secid":   {secid},
		"klt":     {fmt.Sprintf("%d", klt)},
		"fqt":     {fmt.Sprintf("%d", fqt)},
	}

	data, err := DefaultClient.Get(klineAPIURL, params)
	if err != nil {
		return nil, nil, err
	}

	var resp struct {
		Data struct {
			Klines []string `json:"klines"`
			Name   string   `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, data, err
	}

	return &KlineResponse{
		Name:   resp.Data.Name,
		Klines: resp.Data.Klines,
	}, data, nil
}

func SecIDCode(secid string) string {
	if idx := strings.Index(secid, "."); idx >= 0 && idx+1 < len(secid) {
		return secid[idx+1:]
	}
	return secid
}
