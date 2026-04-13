package bond

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

var bondBaseInfoFields = map[string]string{
	"SECURITY_CODE":         "code",
	"SECURITY_NAME_ABBR":    "name",
	"CONVERT_STOCK_CODE":    "stock_code",
	"SECURITY_SHORT_NAME":   "stock_name",
	"RATING":                "rating",
	"PUBLIC_START_DATE":     "sub_date",
	"ACTUAL_ISSUE_SCALE":    "issue_size",
	"LISTING_DATE":          "list_date",
	"EXPIRE_DATE":           "expire_date",
	"BOND_EXPIRE":           "term",
	"INTEREST_RATE_EXPLAIN": "rate_desc",
}

// GetBaseInfo returns bond basic info
func GetBaseInfo(code string) (model.Record, error) {
	params := url.Values{
		"reportName": {"RPT_BOND_CB_LIST"},
		"columns":    {"ALL"},
		"source":     {"WEB"},
		"client":     {"WEB"},
		"filter":     {fmt.Sprintf(`(SECURITY_CODE="%s")`, code)},
	}

	data, err := client.DefaultClient.Get("https://datacenter-web.eastmoney.com/api/data/v1/get", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result struct {
			Data []map[string]interface{} `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if len(resp.Result.Data) == 0 {
		return nil, fmt.Errorf("bond code not found: %s", code)
	}

	item := resp.Result.Data[0]
	r := make(model.Record)
	for fk, ok := range bondBaseInfoFields {
		if v, exists := item[fk]; exists {
			r[ok] = v
		}
	}
	return r, nil
}

// GetAllBaseInfo returns all bond base info (paginated)
func GetAllBaseInfo() ([]model.Record, error) {
	records := make([]model.Record, 0)
	page := 1
	for {
		params := url.Values{
			"sortColumns": {"PUBLIC_START_DATE"},
			"sortTypes":   {"-1"},
			"pageSize":    {"500"},
			"pageNumber":  {fmt.Sprintf("%d", page)},
			"reportName":  {"RPT_BOND_CB_LIST"},
			"columns":     {"ALL"},
			"source":      {"WEB"},
			"client":      {"WEB"},
		}

		data, err := client.DefaultClient.Get("https://datacenter-web.eastmoney.com/api/data/v1/get", params)
		if err != nil {
			return nil, err
		}

		var resp struct {
			Result struct {
				Data []map[string]interface{} `json:"data"`
			} `json:"result"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, err
		}
		if len(resp.Result.Data) == 0 {
			break
		}

		for _, item := range resp.Result.Data {
			r := make(model.Record)
			for fk, ok := range bondBaseInfoFields {
				if v, exists := item[fk]; exists {
					r[ok] = v
				}
			}
			records = append(records, r)
		}
		page++
	}
	return records, nil
}

// GetRealtime returns realtime bond quotes
func GetRealtime() ([]model.Record, error) {
	resp, err := client.FetchQuoteList("b:MK0354", "f12,f14,f3,f2,f15,f16,f17,f4,f8,f10,f9,f5,f6,f18,f20,f21,f13", 1, 200)
	if err != nil {
		return nil, err
	}
	return model.ParseQuoteRecords(resp.Diff), nil
}

// GetHistory returns bond K-line history
func GetHistory(secid, beg, end string, klt, fqt int) ([]model.Record, error) {
	resp, err := client.FetchKlines(secid, beg, end, klt, fqt)
	if err != nil {
		return nil, err
	}
	return model.ParseKlineRecords(resp.Name, client.SecIDCode(secid), resp.Klines), nil
}
