package stock

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var quoteSchema = cliutil.CommandSchema{
	Command: "stock quote",
	Entity:  "stock_quote",
	Supports: map[string]interface{}{
		"format":      []string{"csv", "json", "table"},
		"raw":         true,
		"schema":      true,
		"list_fields": true,
		"sort":        true,
	},
	DefaultFields: []string{"code", "name", "price", "pct", "chg", "open", "high", "low", "pre_close", "vol", "amount", "turnover", "vol_ratio", "pe", "mkt_cap", "float_cap", "latest_trade_date", "update_ts", "secid", "mkt_num"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "price", Type: "number", Desc: "最新价"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
		{Name: "chg", Type: "number", Desc: "涨跌额"},
		{Name: "open", Type: "number", Desc: "开盘价"},
		{Name: "high", Type: "number", Desc: "最高价"},
		{Name: "low", Type: "number", Desc: "最低价"},
		{Name: "pre_close", Type: "number", Desc: "昨收价"},
		{Name: "vol", Type: "number", Desc: "成交量"},
		{Name: "amount", Type: "number", Desc: "成交额"},
		{Name: "turnover", Type: "number", Desc: "换手率"},
		{Name: "vol_ratio", Type: "number", Desc: "量比"},
		{Name: "pe", Type: "number", Desc: "市盈率"},
		{Name: "mkt_cap", Type: "number", Desc: "总市值"},
		{Name: "float_cap", Type: "number", Desc: "流通市值"},
		{Name: "mkt_num", Type: "string", Desc: "市场编号"},
		{Name: "latest_trade_date", Type: "string", Desc: "最新交易日"},
		{Name: "update_ts", Type: "number", Desc: "更新时间戳"},
		{Name: "secid", Type: "string", Desc: "市场证券标识"},
	},
}

func GetQuote(codes []string) ([]model.Record, error) {
	return fetchQuoteRecords(codes)
}

func GetQuoteRaw(codes []string) ([]byte, error) {
	return fetchQuoteRaw(codes)
}

func GetBaseInfo(code string) (model.Record, error) {
	secid, err := client.ResolveQuoteID(code)
	if err != nil {
		return nil, err
	}
	return getBaseInfoBySecID(secid)
}

func getBaseInfoBySecID(secid string) (model.Record, error) {
	fields := make([]string, 0, len(model.BaseInfoFields))
	for k := range model.BaseInfoFields {
		fields = append(fields, k)
	}
	params := url.Values{
		"ut":     {"fa5fd1943c7b386f172d6893dbfba10b"},
		"invt":   {"2"},
		"fltt":   {"2"},
		"fields": {strings.Join(fields, ",")},
		"secid":  {secid},
	}

	data, err := client.DefaultClient.Get("https://push2.eastmoney.com/api/qt/stock/get", params)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", secid, "stock_base_info", err, nil)
	}

	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", secid, "stock_base_info", err, nil)
	}

	r := make(model.Record)
	for fk, ok := range model.BaseInfoFields {
		if v, exists := resp.Data[fk]; exists {
			r[ok] = v
		}
	}
	return r, nil
}
