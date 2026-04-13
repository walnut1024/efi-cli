package stock

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

var ipoFields = map[string]string{
	"ISSUER_NAME":   "name",
	"CHECK_STATUS":  "status",
	"REG_ADDRESS":   "region",
	"CSRC_INDUSTRY": "industry",
	"RECOMMEND_ORG": "sponsor",
	"ACCOUNT_FIRM":  "accounting_firm",
	"UPDATE_DATE":   "update_date",
	"ACCEPT_DATE":   "accept_date",
	"TOLIST_MARKET": "target_market",
}

func GetIPOInfo() ([]model.Record, error) {
	records := make([]model.Record, 0)
	page := 1
	for {
		params := url.Values{
			"st":     {"UPDATE_DATE,SECURITY_CODE"},
			"sr":     {"-1,-1"},
			"ps":     {"500"},
			"p":      {fmt.Sprintf("%d", page)},
			"type":   {"RPT_REGISTERED_INFO"},
			"sty":    {"ORG_CODE,ISSUER_NAME,CHECK_STATUS,CHECK_STATUS_CODE,REG_ADDRESS,CSRC_INDUSTRY,RECOMMEND_ORG,LAW_FIRM,ACCOUNT_FIRM,UPDATE_DATE,ACCEPT_DATE,TOLIST_MARKET,SECURITY_CODE"},
			"token":  {"894050c76af8597a853f5b408b759f5d"},
			"client": {"WEB"},
		}

		data, err := client.DefaultClient.Get("https://datacenter-web.eastmoney.com/api/data/get", params)
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
			for fk, ok := range ipoFields {
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

var holderFields = map[string]string{
	"SECURITY_CODE":      "code",
	"SECURITY_NAME_ABBR": "name",
	"HOLDER_NUM":         "holder_num",
	"HOLDER_NUM_RATIO":   "holder_change_pct",
	"HOLDER_NUM_CHANGE":  "holder_change",
	"END_DATE":           "end_date",
	"AVG_MARKET_CAP":     "avg_market_cap",
	"AVG_HOLD_NUM":       "avg_hold_num",
	"TOTAL_MARKET_CAP":   "total_cap",
	"TOTAL_A_SHARES":     "total_shares",
	"HOLD_NOTICE_DATE":   "notice_date",
}

func GetHolderNumber(date string) ([]model.Record, error) {
	records := make([]model.Record, 0)
	page := 1
	for {
		params := url.Values{
			"sortColumns":  {"HOLD_NOTICE_DATE,SECURITY_CODE"},
			"sortTypes":    {"-1,-1"},
			"pageSize":     {"500"},
			"pageNumber":   {fmt.Sprintf("%d", page)},
			"columns":      {"SECURITY_CODE,SECURITY_NAME_ABBR,END_DATE,INTERVAL_CHRATE,AVG_MARKET_CAP,AVG_HOLD_NUM,TOTAL_MARKET_CAP,TOTAL_A_SHARES,HOLD_NOTICE_DATE,HOLDER_NUM,PRE_HOLDER_NUM,HOLDER_NUM_CHANGE,HOLDER_NUM_RATIO,END_DATE,PRE_END_DATE"},
			"quoteColumns": {"f2,f3"},
			"source":       {"WEB"},
			"client":       {"WEB"},
		}
		if date != "" {
			params.Set("filter", fmt.Sprintf("(END_DATE='%s')", date))
			params.Set("reportName", "RPT_HOLDERNUM_DET")
		} else {
			params.Set("reportName", "RPT_HOLDERNUMLATEST")
		}

		data, err := client.DefaultClient.Get("https://datacenter-web.eastmoney.com/api/data/v1/get", params)
		if err != nil {
			return nil, err
		}

		var resp struct {
			Result struct {
				Count int                      `json:"count"`
				Data  []map[string]interface{} `json:"data"`
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
			for fk, ok := range holderFields {
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

var performanceFields = map[string]string{
	"SECURITY_CODE":        "code",
	"SECURITY_NAME_ABBR":   "name",
	"NOTICE_DATE":          "notice_date",
	"TOTAL_OPERATE_INCOME": "revenue",
	"YSTZ":                 "revenue_yoy",
	"YSHZ":                 "revenue_qoq",
	"PARENT_NETPROFIT":     "net_profit",
	"SJLTZ":                "profit_yoy",
	"SJLHZ":                "profit_qoq",
	"BASIC_EPS":            "eps",
	"BPS":                  "bps",
	"WEIGHTAVG_ROE":        "roe",
	"XSMLL":                "gross_margin",
	"MGJYXJJE":             "cashflow_ps",
}

func GetPerformance(date string) ([]model.Record, error) {
	dateFilter := ""
	if date != "" {
		dateFilter = fmt.Sprintf("(REPORTDATE='%s')", date)
	}

	// First page to get total
	params := url.Values{
		"st":     {"NOTICE_DATE,SECURITY_CODE"},
		"sr":     {"-1,-1"},
		"ps":     {"500"},
		"p":      {"1"},
		"type":   {"RPT_LICO_FN_CPD"},
		"sty":    {"ALL"},
		"token":  {"894050c76af8597a853f5b408b759f5d"},
		"client": {"WEB"},
	}
	if dateFilter != "" {
		params.Set("filter", fmt.Sprintf("(SECURITY_TYPE_CODE in (\"058001001\",\"058001008\"))%s", dateFilter))
	}

	records := make([]model.Record, 0)
	page := 1
	for {
		params.Set("p", fmt.Sprintf("%d", page))
		data, err := client.DefaultClient.Get("https://datacenter-web.eastmoney.com/api/data/get", params)
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
			for fk, ok := range performanceFields {
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
