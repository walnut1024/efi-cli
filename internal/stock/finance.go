package stock

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

const dataCenter = "https://datacenter-web.eastmoney.com/api/data/v1/get"

var financeReportMap = map[string]string{
	"income":   "RPT_DMSK_FN_INCOME",
	"balance":  "RPT_DMSK_FN_BALANCE",
	"cashflow": "RPT_DMSK_FN_CASHFLOW",
}

var financeFieldMaps = map[string]map[string]string{
	"income": {
		"REPORT_DATE":             "report_date",
		"TOTAL_OPERATE_INCOME":    "revenue",
		"OPERATE_PROFIT":          "operate_profit",
		"PARENT_NETPROFIT":        "net_profit",
		"DEDUCT_PARENT_NETPROFIT": "deduct_net_profit",
		"OPERATE_EXPENSE":         "operate_expense",
		"SALE_EXPENSE":            "sale_expense",
		"MANAGE_EXPENSE":          "manage_expense",
		"FINANCE_EXPENSE":         "finance_expense",
		"INCOME_TAX":              "income_tax",
		"TOTAL_OPERATE_COST":      "total_cost",
	},
	"balance": {
		"REPORT_DATE":       "report_date",
		"TOTAL_ASSETS":      "total_assets",
		"TOTAL_LIABILITIES": "total_liabilities",
		"TOTAL_EQUITY":      "total_equity",
		"MONETARYFUNDS":     "cash",
		"ACCOUNTS_RECE":     "receivables",
		"INVENTORY":         "inventory",
		"FIXED_ASSET":       "fixed_assets",
	},
	"cashflow": {
		"REPORT_DATE":     "report_date",
		"NETCASH_OPERATE": "operate_cashflow",
		"NETCASH_INVEST":  "invest_cashflow",
		"NETCASH_FINANCE": "finance_cashflow",
		"CCE_ADD":         "cash_change",
		"SALES_SERVICES":  "sales_receipt",
		"PAY_STAFF_CASH":  "staff_payment",
	},
}

func GetFinance(code, fType string, limit int) ([]model.Record, error) {
	reportName, ok := financeReportMap[fType]
	if !ok {
		return nil, fmt.Errorf("unknown finance type: %s (use income/balance/cashflow)", fType)
	}

	fieldMap := financeFieldMaps[fType]

	params := url.Values{
		"reportName":  {reportName},
		"columns":     {"ALL"},
		"filter":      {fmt.Sprintf(`(SECURITY_CODE="%s")`, code)},
		"pageSize":    {fmt.Sprintf("%d", limit)},
		"sortColumns": {"REPORT_DATE"},
		"sortTypes":   {"-1"},
		"source":      {"WEB"},
		"client":      {"WEB"},
	}

	data, err := client.DefaultClient.Get(dataCenter, params)
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

	records := make([]model.Record, 0, len(resp.Result.Data))
	for _, item := range resp.Result.Data {
		r := make(model.Record)
		r["code"] = code
		for apiField, outField := range fieldMap {
			if v, ok := item[apiField]; ok {
				if s, ok := v.(string); ok {
					r[outField] = model.CleanValue(s)
				} else {
					r[outField] = v
				}
			}
		}
		records = append(records, r)
	}
	return records, nil
}
