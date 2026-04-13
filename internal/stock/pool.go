package stock

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

// Pool field maps for each topic pool type.
var poolFieldMaps = map[string]map[string]string{
	"zt": {
		"f12": "code", "f14": "name", "f2": "price", "f3": "pct", "f6": "amount",
		"f8": "turnover", "f62": "seal_amount", "f17": "first_seal_time",
		"f18": "last_seal_time", "f20": "total_cap", "f21": "float_cap",
		"f33": "break_count", "f84": "consecutive", "f128": "industry",
	},
	"dt": {
		"f12": "code", "f14": "name", "f2": "price", "f3": "pct", "f6": "amount",
		"f9": "pe", "f8": "turnover", "f62": "seal_amount", "f18": "last_seal_time",
		"f20": "total_cap", "f21": "float_cap", "f33": "break_count",
		"f84": "consecutive", "f128": "industry",
	},
	"zb": {
		"f12": "code", "f14": "name", "f2": "price", "f3": "pct", "f6": "amount",
		"f8": "turnover", "f20": "total_cap", "f21": "float_cap", "f128": "industry",
	},
	"strong": {
		"f12": "code", "f14": "name", "f2": "price", "f3": "pct", "f6": "amount",
		"f8": "turnover", "f20": "total_cap", "f21": "float_cap", "f128": "industry",
	},
	"sub-new": {
		"f12": "code", "f14": "name", "f2": "price", "f3": "pct", "f6": "amount",
		"f8": "turnover", "f20": "total_cap", "f21": "float_cap", "f128": "industry",
	},
	"zt-prev": {
		"f12": "code", "f14": "name", "f2": "price", "f3": "pct", "f6": "amount",
		"f8": "turnover", "f20": "total_cap", "f21": "float_cap", "f128": "industry",
	},
}

// poolEndpoints maps pool type to push2ex URL path.
var poolEndpoints = map[string]string{
	"zt":      "getTopicZTPool",
	"zt-prev": "getYesterdayZTPool",
	"dt":      "getTopicDTPool",
	"zb":      "getTopicZBPool",
	"strong":  "getTopicQSPool",
	"sub-new": "getTopicCXPooll",
}

var poolZTSchema = cliutil.CommandSchema{
	Command: "stock pool zt",
	Entity:  "stock_pool_zt",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"}, "sort": true,
	},
	DefaultFields: []string{"code", "name", "pct", "price", "amount", "float_cap", "turnover", "seal_amount", "consecutive", "industry"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
		{Name: "price", Type: "number", Desc: "最新价"},
		{Name: "amount", Type: "number", Desc: "成交额"},
		{Name: "float_cap", Type: "number", Desc: "流通市值"},
		{Name: "total_cap", Type: "number", Desc: "总市值"},
		{Name: "turnover", Type: "number", Desc: "换手率"},
		{Name: "seal_amount", Type: "number", Desc: "封单金额"},
		{Name: "first_seal_time", Type: "string", Desc: "首次封板时间"},
		{Name: "last_seal_time", Type: "string", Desc: "最后封板时间"},
		{Name: "break_count", Type: "number", Desc: "炸板次数"},
		{Name: "consecutive", Type: "number", Desc: "连板数"},
		{Name: "industry", Type: "string", Desc: "所属行业"},
	},
}

var poolDTSchema = cliutil.CommandSchema{
	Command: "stock pool dt",
	Entity:  "stock_pool_dt",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"}, "sort": true,
	},
	DefaultFields: []string{"code", "name", "pct", "price", "amount", "float_cap", "turnover", "seal_amount", "consecutive", "industry"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
		{Name: "price", Type: "number", Desc: "最新价"},
		{Name: "amount", Type: "number", Desc: "成交额"},
		{Name: "float_cap", Type: "number", Desc: "流通市值"},
		{Name: "total_cap", Type: "number", Desc: "总市值"},
		{Name: "pe", Type: "number", Desc: "市盈率"},
		{Name: "turnover", Type: "number", Desc: "换手率"},
		{Name: "seal_amount", Type: "number", Desc: "封单金额"},
		{Name: "last_seal_time", Type: "string", Desc: "最后封板时间"},
		{Name: "break_count", Type: "number", Desc: "炸板次数"},
		{Name: "consecutive", Type: "number", Desc: "连板数"},
		{Name: "industry", Type: "string", Desc: "所属行业"},
	},
}

var poolCommonSchema = cliutil.CommandSchema{
	Command: "stock pool",
	Entity:  "stock_pool",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"}, "sort": true,
	},
	DefaultFields: []string{"code", "name", "pct", "price", "amount", "float_cap", "turnover", "industry"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
		{Name: "price", Type: "number", Desc: "最新价"},
		{Name: "amount", Type: "number", Desc: "成交额"},
		{Name: "float_cap", Type: "number", Desc: "流通市值"},
		{Name: "total_cap", Type: "number", Desc: "总市值"},
		{Name: "turnover", Type: "number", Desc: "换手率"},
		{Name: "industry", Type: "string", Desc: "所属行业"},
	},
}

func GetPool(poolType, date string) ([]model.Record, error) {
	endpoint, ok := poolEndpoints[poolType]
	if !ok {
		return nil, fmt.Errorf("unknown pool type: %s", poolType)
	}
	if date == "" {
		date = "today"
	}

	fieldMap := poolFieldMaps[poolType]
	params := url.Values{
		"ut":   {"7eea3edcaed734bea9cb3f6a79746447"},
		"dpt":  {"wz.ztzt"},
		"date": {date},
	}

	apiURL := "https://push2ex.eastmoney.com/" + endpoint
	data, err := client.DefaultClient.Get(apiURL, params)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", poolType, "stock_pool", err, nil)
	}

	var resp struct {
		RC   int `json:"rc"`
		Data struct {
			Pool []map[string]interface{} `json:"pool"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", poolType, "stock_pool", err, nil)
	}

	if resp.RC != 0 || resp.Data.Pool == nil {
		return nil, nil
	}

	records := make([]model.Record, 0, len(resp.Data.Pool))
	for _, item := range resp.Data.Pool {
		r := make(model.Record)
		for fk, rk := range fieldMap {
			if v, exists := item[fk]; exists {
				r[rk] = v
			}
		}
		records = append(records, r)
	}
	return records, nil
}

func GetPoolRaw(poolType, date string) ([]byte, error) {
	endpoint, ok := poolEndpoints[poolType]
	if !ok {
		return nil, fmt.Errorf("unknown pool type: %s", poolType)
	}
	if date == "" {
		date = "today"
	}
	params := url.Values{
		"ut":   {"7eea3edcaed734bea9cb3f6a79746447"},
		"dpt":  {"wz.ztzt"},
		"date": {date},
	}
	apiURL := "https://push2ex.eastmoney.com/" + endpoint
	return client.DefaultClient.Get(apiURL, params)
}

func poolSchemaForType(poolType string) *cliutil.CommandSchema {
	switch poolType {
	case "zt":
		return &poolZTSchema
	case "dt":
		return &poolDTSchema
	default:
		return &poolCommonSchema
	}
}
