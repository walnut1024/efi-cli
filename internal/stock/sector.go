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

var sectorListFields = "f12,f14,f2,f3,f4,f20,f8,f104,f105,f128,f140"
var sectorMemberFields = "f12,f14,f2,f3,f4,f5,f6,f7,f15,f16,f17,f18,f8,f9,f23"

var sectorListFieldMap = map[string]string{
	"f12": "code", "f14": "name", "f2": "price", "f3": "pct", "f4": "chg",
	"f20": "total_cap", "f8": "turnover", "f104": "up_count",
	"f105": "down_count", "f128": "lead_stock", "f140": "lead_code",
}

var sectorMemberFieldMap = map[string]string{
	"f12": "code", "f14": "name", "f2": "price", "f3": "pct", "f4": "chg",
	"f5": "vol", "f6": "amt", "f7": "amplitude", "f15": "high", "f16": "low",
	"f17": "open", "f18": "pre_close", "f8": "turnover", "f9": "pe", "f23": "pb",
}

var sectorListSchema = cliutil.CommandSchema{
	Command: "stock sector list",
	Entity:  "stock_sector_list",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"}, "sort": true,
	},
	DefaultFields: []string{"code", "name", "pct", "chg", "total_cap", "turnover", "up_count", "down_count", "lead_stock"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "板块代码"},
		{Name: "name", Type: "string", Desc: "板块名称"},
		{Name: "price", Type: "number", Desc: "板块指数"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
		{Name: "chg", Type: "number", Desc: "涨跌额"},
		{Name: "total_cap", Type: "number", Desc: "总市值"},
		{Name: "turnover", Type: "number", Desc: "换手率"},
		{Name: "up_count", Type: "number", Desc: "上涨家数"},
		{Name: "down_count", Type: "number", Desc: "下跌家数"},
		{Name: "lead_stock", Type: "string", Desc: "领涨股"},
		{Name: "lead_code", Type: "string", Desc: "领涨股代码"},
	},
}

var sectorMembersSchema = cliutil.CommandSchema{
	Command: "stock sector members",
	Entity:  "stock_sector_members",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"}, "sort": true,
	},
	DefaultFields: []string{"code", "name", "pct", "price", "chg", "vol", "amt", "turnover", "pe", "pb"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "price", Type: "number", Desc: "最新价"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
		{Name: "chg", Type: "number", Desc: "涨跌额"},
		{Name: "vol", Type: "number", Desc: "成交量"},
		{Name: "amt", Type: "number", Desc: "成交额"},
		{Name: "amplitude", Type: "number", Desc: "振幅"},
		{Name: "high", Type: "number", Desc: "最高"},
		{Name: "low", Type: "number", Desc: "最低"},
		{Name: "open", Type: "number", Desc: "开盘"},
		{Name: "pre_close", Type: "number", Desc: "昨收"},
		{Name: "turnover", Type: "number", Desc: "换手率"},
		{Name: "pe", Type: "number", Desc: "市盈率"},
		{Name: "pb", Type: "number", Desc: "市净率"},
	},
}

func sectorTypeFS(sectorType string) string {
	switch sectorType {
	case "industry":
		return "m:90+t:2"
	case "concept":
		return "m:90+t:3"
	default:
		return "m:90+t:3"
	}
}

func GetSectorList(sectorType string) ([]model.Record, error) {
	fs := sectorTypeFS(sectorType)
	params := url.Values{
		"pn": {"1"}, "pz": {"1000"}, "po": {"1"}, "np": {"1"},
		"fltt": {"2"}, "invt": {"2"}, "fid": {"f3"},
		"fs": {fs}, "fields": {sectorListFields},
	}
	data, err := client.DefaultClient.Get("https://push2.eastmoney.com/api/qt/clist/get", params)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", sectorType, "sector_list", err, nil)
	}

	var resp struct {
		Data struct {
			Total int                      `json:"total"`
			Diff  []map[string]interface{} `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", sectorType, "sector_list", err, nil)
	}

	return parseSectorRecords(resp.Data.Diff, sectorListFieldMap), nil
}

func GetSectorMembers(code string, limit int) ([]model.Record, error) {
	fs := "b:" + code + "+f:!50"
	pz := 500
	if limit > 0 && limit < pz {
		pz = limit
	}
	params := url.Values{
		"pn": {"1"}, "pz": {fmt.Sprintf("%d", pz)}, "po": {"1"}, "np": {"1"},
		"fltt": {"2"}, "invt": {"2"}, "fid": {"f12"},
		"fs": {fs}, "fields": {sectorMemberFields},
	}
	data, err := client.DefaultClient.Get("https://push2.eastmoney.com/api/qt/clist/get", params)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", code, "sector_members", err, nil)
	}

	var resp struct {
		Data struct {
			Total int                      `json:"total"`
			Diff  []map[string]interface{} `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", code, "sector_members", err, nil)
	}

	return parseSectorRecords(resp.Data.Diff, sectorMemberFieldMap), nil
}

func GetSectorHistory(code, beg, end string, klt, fqt int) ([]model.Record, error) {
	secid := "90." + code
	resp, err := client.FetchKlines(secid, beg, end, klt, fqt)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", code, "sector_history", err, nil)
	}

	records := make([]model.Record, 0, len(resp.Klines))
	for _, line := range resp.Klines {
		parts := strings.Split(line, ",")
		r := make(model.Record)
		for i, fk := range model.KlineFieldKeys {
			if i < len(parts) {
				r[model.KlineFields[fk]] = model.CleanValue(parts[i])
			}
		}
		records = append(records, r)
	}
	return records, nil
}

func GetSectorQuote(code string) ([]model.Record, error) {
	secid := "90." + code
	fields := make([]string, 0, len(sectorListFieldMap))
	for k := range sectorListFieldMap {
		fields = append(fields, k)
	}
	params := url.Values{
		"ut":   {"fa5fd1943c7b386f172d6893dbfba10b"},
		"fltt": {"2"}, "invt": {"2"},
		"fields": {strings.Join(fields, ",")},
		"secid":  {secid},
	}
	data, err := client.DefaultClient.Get("https://push2.eastmoney.com/api/qt/stock/get", params)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", code, "sector_quote", err, nil)
	}

	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", code, "sector_quote", err, nil)
	}

	r := make(model.Record)
	for fk, rk := range sectorListFieldMap {
		if v, exists := resp.Data[fk]; exists {
			r[rk] = v
		}
	}
	return []model.Record{r}, nil
}

func parseSectorRecords(items []map[string]interface{}, fieldMap map[string]string) []model.Record {
	records := make([]model.Record, 0, len(items))
	for _, item := range items {
		r := make(model.Record)
		for fk, rk := range fieldMap {
			if v, exists := item[fk]; exists {
				r[rk] = v
			}
		}
		records = append(records, r)
	}
	return records
}
