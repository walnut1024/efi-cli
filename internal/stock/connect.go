package stock

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

// Field maps for datacenter report responses (PascalCase → snake_case).
var connectHistoryFields = map[string]string{
	"TRADE_DATE": "date", "FUND_INFLOW": "fund_inflow", "NET_DEAL_AMT": "net_deal_amt",
	"QUOTA_BALANCE": "quota_balance", "ACCUM_DEAL_AMT": "accum_deal_amt",
	"BUY_AMT": "buy_amt", "SELL_AMT": "sell_amt", "HOLD_MARKET_CAP": "hold_market_cap",
	"LEAD_STOCKS_NAME": "lead_stock", "LEAD_STOCKS_CODE": "lead_code",
	"INDEX_CLOSE_PRICE": "index_close", "INDEX_CHANGE_RATE": "index_pct",
}

var connectHoldRankFields = map[string]string{
	"TRADE_DATE": "date", "SECURITY_CODE": "code", "SECURITY_NAME_ABBR": "name",
	"CLOSE_PRICE": "close", "CHANGE_RATE": "pct",
	"HOLD_SHARES": "hold_shares", "HOLD_MARKET_CAP": "hold_market_cap",
	"HOLD_FLOAT_RATIO": "hold_float_ratio", "HOLD_TOTAL_RATIO": "hold_total_ratio",
	"ADD_SHARES": "add_shares", "ADD_MARKET_CAP": "add_market_cap",
	"ADD_RATIO": "add_pct",
}

var connectHoldDetailFields = map[string]string{
	"HOLD_DATE": "date", "HOLD_SHARES": "hold_shares",
	"HOLD_MARKET_CAP": "hold_market_cap", "HOLD_FLOAT_RATIO": "hold_float_ratio",
	"HOLD_TOTAL_RATIO": "hold_total_ratio", "ADD_SHARES": "add_shares",
	"ADD_MARKET_CAP": "add_market_cap", "ADD_RATIO": "add_pct",
	"CLOSE_PRICE": "close", "CHANGE_RATE": "pct",
}

var connectSummaryFields = map[string]string{
	"TRADE_DATE": "date", "MUTUAL_TYPE": "mutual_type", "MUTUAL_TYPE_NAME": "mutual_type_name",
	"BOARD_TYPE": "board_type", "FUNDS_DIRECTION": "direction",
	"TRADE_QUOTA": "trade_quota", "INDEX_NAME": "index_name",
}

var connectHistorySchema = cliutil.CommandSchema{
	Command: "stock connect history",
	Entity:  "stock_connect_history",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"},
	},
	DefaultFields: []string{"date", "fund_inflow", "net_deal_amt", "buy_amt", "sell_amt", "hold_market_cap"},
	Fields: []cliutil.FieldSchema{
		{Name: "date", Type: "string", Desc: "交易日期"},
		{Name: "fund_inflow", Type: "number", Desc: "资金流入"},
		{Name: "net_deal_amt", Type: "number", Desc: "净成交额"},
		{Name: "buy_amt", Type: "number", Desc: "买入额"},
		{Name: "sell_amt", Type: "number", Desc: "卖出额"},
		{Name: "hold_market_cap", Type: "number", Desc: "持股市值"},
		{Name: "quota_balance", Type: "number", Desc: "余额"},
		{Name: "lead_stock", Type: "string", Desc: "领涨股"},
		{Name: "lead_code", Type: "string", Desc: "领涨股代码"},
	},
}

var connectHoldRankSchema = cliutil.CommandSchema{
	Command: "stock connect hold-rank",
	Entity:  "stock_connect_hold_rank",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"}, "sort": true,
	},
	DefaultFields: []string{"code", "name", "close", "pct", "hold_shares", "hold_market_cap", "add_shares", "add_pct"},
	Fields: []cliutil.FieldSchema{
		{Name: "date", Type: "string", Desc: "日期"},
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "close", Type: "number", Desc: "收盘价"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
		{Name: "hold_shares", Type: "number", Desc: "持股数量"},
		{Name: "hold_market_cap", Type: "number", Desc: "持股市值"},
		{Name: "hold_float_ratio", Type: "number", Desc: "占流通股比"},
		{Name: "hold_total_ratio", Type: "number", Desc: "占总股本比"},
		{Name: "add_shares", Type: "number", Desc: "增持股数"},
		{Name: "add_market_cap", Type: "number", Desc: "增持市值"},
		{Name: "add_pct", Type: "number", Desc: "增持比例"},
	},
}

var connectSummarySchema = cliutil.CommandSchema{
	Command: "stock connect summary",
	Entity:  "stock_connect_summary",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"},
	},
	DefaultFields: []string{"date", "mutual_type_name", "direction", "board_type", "index_name"},
	Fields: []cliutil.FieldSchema{
		{Name: "date", Type: "string", Desc: "交易日期"},
		{Name: "mutual_type_name", Type: "string", Desc: "通道路径"},
		{Name: "direction", Type: "string", Desc: "资金方向"},
		{Name: "board_type", Type: "string", Desc: "板块类型"},
		{Name: "index_name", Type: "string", Desc: "关联指数"},
	},
}

var connectAHSchema = cliutil.CommandSchema{
	Command: "stock connect ah",
	Entity:  "stock_connect_ah",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"}, "sort": true,
	},
	DefaultFields: []string{"code", "name", "price", "pct"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "price", Type: "number", Desc: "最新价"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
	},
}

var connectRealtimeSchema = cliutil.CommandSchema{
	Command: "stock connect realtime",
	Entity:  "stock_connect_realtime",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"},
	},
	DefaultFields: []string{"time", "sh_connect", "sz_connect", "north_total"},
	Fields: []cliutil.FieldSchema{
		{Name: "time", Type: "string", Desc: "时间"},
		{Name: "sh_connect", Type: "number", Desc: "沪股通(万元)"},
		{Name: "sz_connect", Type: "number", Desc: "深股通(万元)"},
		{Name: "north_total", Type: "number", Desc: "北向合计(万元)"},
	},
}

func GetConnectSummary() ([]model.Record, error) {
	data, err := client.FetchDatacenterReport("RPT_MUTUAL_QUOTA", nil, 1, 100)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", "", "connect_summary", err, nil)
	}
	return mapDatacenterRecords(data.Data, connectSummaryFields), nil
}

func GetConnectHistory(connType string, limit int) ([]model.Record, error) {
	mutualType := connectMutualType(connType)
	filters := url.Values{}
	if mutualType != "" {
		filters.Set("filter", fmt.Sprintf("(MUTUAL_TYPE=\"%s\")", mutualType))
	}
	filters.Set("sortColumns", "TRADE_DATE")
	filters.Set("sortTypes", "-1")

	data, err := client.FetchDatacenterReportAll("RPT_MUTUAL_DEAL_HISTORY", filters, limit)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", connType, "connect_history", err, nil)
	}
	return mapDatacenterRecords(data, connectHistoryFields), nil
}

func GetConnectRealtime() ([]model.Record, error) {
	params := url.Values{
		"fields1": {"f1,f2,f3,f4"},
		"fields2": {"f51,f52,f53,f54,f55,f56"},
	}
	data, err := client.DefaultClient.Get("https://push2.eastmoney.com/api/qt/kamtbs.rtmin/get", params)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", "", "connect_realtime", err, nil)
	}

	var resp struct {
		Data struct {
			S2N []string `json:"s2n"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", "", "connect_realtime", err, nil)
	}

	records := make([]model.Record, 0, len(resp.Data.S2N))
	for _, line := range resp.Data.S2N {
		parts := strings.Split(line, ",")
		if len(parts) < 4 {
			continue
		}
		r := make(model.Record)
		r["time"] = parts[0]
		r["sh_connect"] = model.CleanValue(parts[1])
		r["sz_connect"] = model.CleanValue(parts[2])
		r["north_total"] = model.CleanValue(parts[3])
		records = append(records, r)
	}
	return records, nil
}

func GetConnectHoldRank(rankType, indicator, date string) ([]model.Record, error) {
	mutualType := connectRankMutualType(rankType)
	filter := fmt.Sprintf("(MUTUAL_TYPE=\"%s\")", mutualType)
	if indicator != "" {
		filter += fmt.Sprintf("(INDICATOR_TYPE=\"%s\")", indicator)
	}
	if date != "" {
		filter += fmt.Sprintf("(TRADE_DATE=\"%s\")", date)
	}
	filters := url.Values{
		"filter":      {filter},
		"sortColumns": {"HOLD_MARKET_CAP"},
		"sortTypes":   {"-1"},
	}

	data, err := client.FetchDatacenterReportAll("RPT_MUTUAL_STOCK_NORTHSTA", filters, 50)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", rankType, "connect_hold_rank", err, nil)
	}
	return mapDatacenterRecords(data, connectHoldRankFields), nil
}

func GetConnectHold(code, start, end string) ([]model.Record, error) {
	filters := url.Values{
		"filter":      {fmt.Sprintf("(SECURITY_CODE=\"%s\")", code)},
		"sortColumns": {"HOLD_DATE"},
		"sortTypes":   {"-1"},
	}
	if start != "" || end != "" {
		dateFilter := fmt.Sprintf("(SECURITY_CODE=\"%s\")", code)
		if start != "" {
			dateFilter += fmt.Sprintf("(HOLD_DATE>='%s')", start)
		}
		if end != "" {
			dateFilter += fmt.Sprintf("(HOLD_DATE<='%s')", end)
		}
		filters.Set("filter", dateFilter)
	}

	data, err := client.FetchDatacenterReportAll("RPT_MUTUAL_HOLD_DET", filters, 100)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", code, "connect_hold", err, nil)
	}

	records := mapDatacenterRecords(data, connectHoldDetailFields)
	for i := range records {
		records[i]["code"] = code
	}
	return records, nil
}

func GetConnectAH() ([]model.Record, error) {
	fs := "b:BK0500"
	fields := "f12,f14,f2,f3"
	params := url.Values{
		"pn": {"1"}, "pz": {"200"}, "po": {"1"}, "np": {"1"},
		"fltt": {"2"}, "invt": {"2"}, "fid": {"f3"},
		"fs": {fs}, "fields": {fields},
	}
	data, err := client.DefaultClient.Get("https://push2.eastmoney.com/api/qt/clist/get", params)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", "", "connect_ah", err, nil)
	}

	var resp struct {
		Data struct {
			Diff []map[string]interface{} `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", "", "connect_ah", err, nil)
	}

	ahMap := map[string]string{"f12": "code", "f14": "name", "f2": "price", "f3": "pct"}
	return mapDatacenterRecords(resp.Data.Diff, ahMap), nil
}

func connectMutualType(connType string) string {
	switch connType {
	case "north":
		return "001"
	case "sh":
		return "001"
	case "sz":
		return "003"
	case "south":
		return "002"
	case "hshk":
		return "004"
	case "szhk":
		return "004"
	default:
		return "001"
	}
}

func connectRankMutualType(rankType string) string {
	switch rankType {
	case "north":
		return "1"
	case "sh":
		return "1"
	case "sz":
		return "3"
	default:
		return "1"
	}
}

func mapDatacenterRecords(items []map[string]interface{}, fieldMap map[string]string) []model.Record {
	records := make([]model.Record, 0, len(items))
	for _, item := range items {
		r := make(model.Record)
		for dk, rk := range fieldMap {
			if v, exists := item[dk]; exists {
				if rk == "date" || strings.HasSuffix(rk, "_date") {
					if s, ok := v.(string); ok && len(s) > 10 {
						v = s[:10]
					}
				}
				r[rk] = v
			}
		}
		records = append(records, r)
	}
	return records
}
