package stock

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

func GetBelongBoard(code string) ([]model.Record, error) {
	quote, err := client.SearchQuote(code)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, fmt.Errorf("code not found: %s", code)
	}

	params := url.Values{
		"forcect": {"1"},
		"spt":     {"3"},
		"fields":  {"f1,f12,f152,f3,f14,f128,f136"},
		"pi":      {"0"},
		"pz":      {"1000"},
		"po":      {"1"},
		"fid":     {"f3"},
		"fid0":    {"f4003"},
		"invt":    {"2"},
		"secid":   {quote.QuoteID},
	}

	data, err := client.DefaultClient.Get("https://push2.eastmoney.com/api/qt/slist/get", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			Diff map[string]map[string]interface{} `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	boardFieldMap := map[string]string{
		"f12": "board_code",
		"f14": "board_name",
		"f3":  "board_pct",
	}

	records := make([]model.Record, 0, len(resp.Data.Diff))
	for _, item := range resp.Data.Diff {
		r := make(model.Record)
		r["code"] = quote.Code
		r["name"] = quote.Name
		for fk, ok := range boardFieldMap {
			if v, exists := item[fk]; exists {
				r[ok] = v
			}
		}
		records = append(records, r)
	}
	return records, nil
}

func GetMembers(indexCode string) ([]model.Record, error) {
	// Resolve index code to find the actual index
	quote, err := client.SearchQuote(indexCode)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, fmt.Errorf("index not found: %s", indexCode)
	}

	params := url.Values{
		"IndexCode":     {quote.Code},
		"pageIndex":     {"1"},
		"pageSize":      {"10000"},
		"deviceid":      {"1234567890"},
		"version":       {"6.9.9"},
		"product":       {"EFund"},
		"plat":          {"Iphone"},
		"ServerVersion": {"6.9.9"},
	}

	data, err := client.DefaultClient.Get("https://fundztapi.eastmoney.com/FundSpecialApiNew/FundSpecialZSB30ZSCFG", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Datas []map[string]interface{} `json:"Datas"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	memberFields := map[string]string{
		"IndexCode":    "idx_code",
		"IndexName":    "idx_name",
		"StockCode":    "code",
		"StockName":    "name",
		"MARKETCAPPCT": "weight",
	}

	records := make([]model.Record, 0, len(resp.Datas))
	for _, item := range resp.Datas {
		r := make(model.Record)
		for fk, ok := range memberFields {
			if v, exists := item[fk]; exists {
				if strings.HasPrefix(ok, "weight") {
					r[ok] = model.CleanFloat(fmt.Sprintf("%v", v))
				} else {
					r[ok] = v
				}
			}
		}
		records = append(records, r)
	}
	return records, nil
}
