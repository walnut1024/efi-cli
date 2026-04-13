package stock

import (
	"fmt"
	"net/url"

	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var restrictedSummaryFields = map[string]string{
	"FREE_DATE": "date", "LIFT_NUM": "stock_count", "PLAN_LIFT_NUM": "plan_stock_count",
	"ABLE_FREE_SHARES": "free_shares", "MARKET_CAP": "market_cap",
	"INDEX_PRICE": "index_close", "CHANGE_RATE": "index_pct",
}

var restrictedDetailFields = map[string]string{
	"SECURITY_CODE": "code", "SECURITY_NAME_ABBR": "name", "FREE_DATE": "free_date",
	"FREE_SHARES": "free_shares", "ABLE_FREE_SHARES": "actual_free_shares",
	"LIFT_MARKET_CAP": "market_cap", "FREE_RATIO": "free_ratio",
	"TOTAL_RATIO": "total_ratio", "FREE_SHARES_TYPE": "free_shares_type",
	"TOTALSHARES_RATIO": "total_shares_ratio", "BATCH_HOLDER_NUM": "batch_holder_num",
	"B20_ADJCHRATE": "pct_b20", "A20_ADJCHRATE": "pct_a20",
}

var restrictedSchema = cliutil.CommandSchema{
	Command: "stock restricted summary",
	Entity:  "stock_restricted_summary",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"},
	},
	DefaultFields: []string{"date", "stock_count", "free_shares", "market_cap", "index_close", "index_pct"},
	Fields: []cliutil.FieldSchema{
		{Name: "date", Type: "string", Desc: "解禁日期"},
		{Name: "stock_count", Type: "number", Desc: "解禁公司数"},
		{Name: "free_shares", Type: "number", Desc: "解禁股数(万股)"},
		{Name: "market_cap", Type: "number", Desc: "解禁市值(万元)"},
		{Name: "index_close", Type: "number", Desc: "指数收盘"},
		{Name: "index_pct", Type: "number", Desc: "指数涨跌幅"},
	},
}

var restrictedDetailSchema = cliutil.CommandSchema{
	Command: "stock restricted detail",
	Entity:  "stock_restricted_detail",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"},
	},
	DefaultFields: []string{"code", "name", "free_date", "free_shares", "market_cap", "free_ratio", "total_ratio"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "free_date", Type: "string", Desc: "解禁日期"},
		{Name: "free_shares", Type: "number", Desc: "解禁股数"},
		{Name: "actual_free_shares", Type: "number", Desc: "实际解禁数"},
		{Name: "market_cap", Type: "number", Desc: "解禁市值"},
		{Name: "free_ratio", Type: "number", Desc: "解禁比例"},
		{Name: "total_ratio", Type: "number", Desc: "占总股本比"},
	},
}

func GetRestrictedSummary(start, end, market string) ([]model.Record, error) {
	filter := ""
	if start != "" || end != "" {
		filter = addDateFilterCol("FREE_DATE", "", start, end)
	}
	filters := url.Values{
		"sortColumns": {"FREE_DATE"},
		"sortTypes":   {"-1"},
	}
	if filter != "" {
		filters.Set("filter", filter)
	}
	data, err := fetchDatacenter("RPT_LIFTDAY_STA", filters, 100)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", market, "restricted_summary", err, nil)
	}
	return parseDatacenterToRecords(data, restrictedSummaryFields), nil
}

func GetRestrictedDetail(start, end string) ([]model.Record, error) {
	filter := ""
	if start != "" || end != "" {
		filter = addDateFilterCol("FREE_DATE", "", start, end)
	}
	filters := url.Values{
		"sortColumns": {"FREE_DATE"},
		"sortTypes":   {"-1"},
	}
	if filter != "" {
		filters.Set("filter", filter)
	}
	data, err := fetchDatacenter("RPT_LIFT_STAGE", filters, 100)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", "", "restricted_detail", err, nil)
	}
	return parseDatacenterToRecords(data, restrictedDetailFields), nil
}

func GetRestrictedQueue(code string) ([]model.Record, error) {
	filters := url.Values{
		"filter":      {fmt.Sprintf("(SECURITY_CODE=\"%s\")", code)},
		"sortColumns": {"FREE_DATE"},
		"sortTypes":   {"1"},
	}
	data, err := fetchDatacenter("RPT_LIFT_STAGE", filters, 50)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", code, "restricted_queue", err, nil)
	}
	return parseDatacenterToRecords(data, restrictedDetailFields), nil
}

func GetRestrictedHolders(code, date string) ([]model.Record, error) {
	holderFields := map[string]string{
		"HOLDER_NAME": "holder_name", "ADD_SHARES": "add_shares",
		"ACTUAL_SHARES": "actual_shares", "ADD_MARKET_CAP": "add_market_cap",
		"LOCK_MONTHS": "lock_months", "RESIDUAL_SHARES": "residual_shares",
		"FREE_TYPE": "free_type", "PROGRESS": "progress",
	}
	filter := fmt.Sprintf("(SECURITY_CODE=\"%s\")", code)
	if date != "" {
		filter = fmt.Sprintf("(SECURITY_CODE=\"%s\")(FREE_DATE=\"%s\")", code, date)
	}
	filters := url.Values{
		"filter":      {filter},
		"sortColumns": {"FREE_DATE"},
		"sortTypes":   {"-1"},
	}
	data, err := fetchDatacenter("RPT_LIFT_GD", filters, 100)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", code, "restricted_holders", err, nil)
	}
	return parseDatacenterToRecords(data, holderFields), nil
}
