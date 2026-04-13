package index

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

var IndexSecID = map[string]string{
	"上证指数":   "1.000001",
	"深证成指":   "0.399001",
	"创业板指":   "0.399006",
	"沪深300":  "1.000300",
	"中证500":  "1.000905",
	"中证1000": "0.000852",
	"上证50":   "1.000016",
	"科创50":   "1.000688",
	"恒生指数":   "100.HSI",
}

func resolveSecID(code string) string {
	if strings.Contains(code, ".") {
		return code
	}
	if secID, ok := IndexSecID[code]; ok {
		return secID
	}
	return code
}

func ResolveSecID(code string) string {
	return resolveSecID(code)
}

func GetQuote(codes []string) ([]model.Record, error) {
	secIDs := make([]string, len(codes))
	for i, c := range codes {
		secIDs[i] = resolveSecID(c)
	}

	params := url.Values{
		"fields": {"f1,f2,f3,f4,f5,f6,f7,f8,f12,f13,f14,f15,f16,f17,f18"},
		"secids": {strings.Join(secIDs, ",")},
	}
	data, err := client.DefaultClient.Get("https://push2.eastmoney.com/api/qt/ulist.np/get", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			Diff []map[string]interface{} `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	fieldMap := map[string]string{
		"f12": "code",
		"f14": "name",
		"f2":  "price",
		"f3":  "pct",
		"f4":  "chg",
		"f15": "high",
		"f16": "low",
		"f17": "open",
		"f5":  "vol",
		"f6":  "amount",
		"f18": "pre_close",
	}

	records := make([]model.Record, 0, len(resp.Data.Diff))
	for _, item := range resp.Data.Diff {
		r := make(model.Record)
		for fk, ok := range fieldMap {
			if v, exists := item[fk]; exists {
				r[ok] = v
			}
		}
		records = append(records, r)
	}
	return records, nil
}

func GetHistory(code, beg, end string, klt int) ([]model.Record, error) {
	secID := resolveSecID(code)
	if beg == "" {
		beg = "19000101"
	}
	if end == "" {
		end = "20500101"
	}
	if klt == 0 {
		klt = 101
	}

	resp, err := client.FetchKlines(secID, beg, end, klt, 1)
	if err != nil {
		return nil, err
	}
	return model.ParseKlineRecords(resp.Name, "", resp.Klines), nil
}
