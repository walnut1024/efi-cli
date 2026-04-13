package stock

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

func GetDividend(code string, limit int) ([]model.Record, error) {
	// Resolve market prefix for F10 API: SH for 6xx, SZ for 0xx/3xx, BJ for 4xx/8xx
	prefix := "SZ"
	if strings.HasPrefix(code, "6") || strings.HasPrefix(code, "5") {
		prefix = "SH"
	} else if strings.HasPrefix(code, "4") || strings.HasPrefix(code, "8") {
		prefix = "BJ"
	}
	f10Code := prefix + code

	data, err := client.DefaultClient.Get(
		"https://emweb.securities.eastmoney.com/PC_HSF10/BonusFinancing/PageAjax",
		url.Values{"code": {f10Code}},
	)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Fhyx []map[string]interface{} `json:"fhyx"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	records := make([]model.Record, 0, len(resp.Fhyx))
	for _, item := range resp.Fhyx {
		r := make(model.Record)
		r["code"] = code
		r["name"] = cleanNil(strVal(item, "SECURITY_NAME_ABBR"))
		r["notice_date"] = cleanDate(cleanNil(strVal(item, "NOTICE_DATE")))
		r["plan"] = cleanNil(strVal(item, "IMPL_PLAN_PROFILE"))
		r["progress"] = cleanNil(strVal(item, "ASSIGN_PROGRESS"))
		r["record_date"] = cleanDate(cleanNil(strVal(item, "EQUITY_RECORD_DATE")))
		r["ex_date"] = cleanDate(cleanNil(strVal(item, "EX_DIVIDEND_DATE")))
		r["pay_date"] = cleanDate(cleanNil(strVal(item, "PAY_CASH_DATE")))
		records = append(records, r)
	}

	if limit > 0 && len(records) > limit {
		records = records[:limit]
	}
	return records, nil
}

func cleanDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

func cleanNil(s string) string {
	if s == "<nil>" {
		return ""
	}
	return s
}

func strVal(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}
