package stock

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

var billboardFields = map[string]string{
	"SECURITY_CODE":      "code",
	"SECURITY_NAME_ABBR": "name",
	"TRADE_DATE":         "date",
	"EXPLAIN":            "explain",
	"CLOSE_PRICE":        "close",
	"CHANGE_RATE":        "pct",
	"TURNOVER_RATE":      "turnover",
	"NET_BUY_AMT":        "net_buy",
	"BUY_AMT":            "buy",
	"SELL_AMT":           "sell",
	"BILLBOARD_AMT":      "billboard_amt",
	"ACCUMULATE_AMT":     "total_amt",
	"NET_BUY_RATE":       "net_buy_pct",
	"BILLBOARD_RATE":     "billboard_pct",
	"FREE_MARKET_CAP":    "float_cap",
	"REASON_TYPE":        "reason",
}

func GetBillboard(start, end string) ([]model.Record, error) {
	if start == "" || end == "" {
		return nil, fmt.Errorf("start and end dates required (YYYY-MM-DD)")
	}

	records := make([]model.Record, 0)
	page := 1
	maxPages := 40 // cap at ~20000 records to avoid unbounded queries
	for page <= maxPages {
		params := url.Values{
			"sortColumns": {"TRADE_DATE,SECURITY_CODE"},
			"sortTypes":   {"-1,1"},
			"pageSize":    {"500"},
			"pageNumber":  {fmt.Sprintf("%d", page)},
			"reportName":  {"RPT_DAILYBILLBOARD_DETAILS"},
			"columns":     {"ALL"},
			"source":      {"WEB"},
			"client":      {"WEB"},
			"filter":      {fmt.Sprintf("(TRADE_DATE<='%s')(TRADE_DATE>='%s')", end, start)},
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
			for fk, ok := range billboardFields {
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
