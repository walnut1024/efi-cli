package stock

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

func GetHistoryBill(code string) ([]model.Record, error) {
	secid, err := client.ResolveQuoteID(code)
	if err != nil {
		return nil, err
	}

	fields2 := strings.Join(model.BillFieldKeys, ",")
	params := url.Values{
		"lmt":     {"100000"},
		"klt":     {"101"},
		"secid":   {secid},
		"fields1": {"f1,f2,f3,f7"},
		"fields2": {fields2},
	}

	data, err := client.DefaultClient.Get("https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			Klines []string `json:"klines"`
			Name   string   `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	codePart := secid[strings.Index(secid, ".")+1:]
	records := parseBillKlines(resp.Data.Klines, resp.Data.Name, codePart)
	reverseRecords(records)
	return records, nil
}

func GetTodayBill(code string) ([]model.Record, error) {
	secid, err := client.ResolveQuoteID(code)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"lmt":     {"0"},
		"klt":     {"1"},
		"secid":   {secid},
		"fields1": {"f1,f2,f3,f7"},
		"fields2": {"f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63"},
	}

	data, err := client.DefaultClient.Get("https://push2.eastmoney.com/api/qt/stock/fflow/kline/get", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			Klines []string `json:"klines"`
			Name   string   `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	codePart := secid[strings.Index(secid, ".")+1:]
	return parseBillKlines(resp.Data.Klines, resp.Data.Name, codePart), nil
}

func parseBillKlines(klines []string, name, code string) []model.Record {
	records := make([]model.Record, 0, len(klines))
	for _, line := range klines {
		parts := strings.Split(line, ",")
		r := make(model.Record)
		r["name"] = name
		r["code"] = code
		for i, fk := range model.BillFieldKeys {
			if i < len(parts) {
				r[model.BillFields[fk]] = model.CleanValue(parts[i])
			}
		}
		records = append(records, r)
	}
	return records
}

func reverseRecords(records []model.Record) {
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}
}
