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

var hotRankSchema = cliutil.CommandSchema{
	Command: "stock hot-rank",
	Entity:  "stock_hot_rank",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"}, "sort": true,
	},
	DefaultFields: []string{"rank", "code", "name", "price", "pct"},
	Fields: []cliutil.FieldSchema{
		{Name: "rank", Type: "number", Desc: "排名"},
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "price", Type: "number", Desc: "最新价"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
	},
}

var hotRankHistorySchema = cliutil.CommandSchema{
	Command: "stock hot-rank history",
	Entity:  "stock_hot_rank_history",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"},
	},
	DefaultFields: []string{"date", "rank", "code", "new_fans", "loyal_fans"},
	Fields: []cliutil.FieldSchema{
		{Name: "date", Type: "string", Desc: "日期"},
		{Name: "rank", Type: "number", Desc: "排名"},
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "new_fans", Type: "number", Desc: "新增粉丝"},
		{Name: "loyal_fans", Type: "number", Desc: "忠诚粉丝"},
	},
}

const hotRankBaseURL = "https://emappdata.eastmoney.com/stockrank/"

func GetHotRank(limit int) ([]model.Record, error) {
	payload := map[string]any{
		"appId":      "appId01",
		"globalId":   "786e4c21-70dc-435a-93bb-38",
		"marketType": "",
		"pageNo":     1,
		"pageSize":   limit,
	}
	data, err := client.DefaultClient.PostJSON(hotRankBaseURL+"getAllCurrentList", payload)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", "", "hot_rank", err, nil)
	}

	var resp struct {
		Data []struct {
			SC string `json:"sc"`
			RK int    `json:"rk"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", "", "hot_rank", err, nil)
	}

	// Resolve stock codes to get name/price
	records := make([]model.Record, 0, len(resp.Data))
	codes := make([]string, 0, len(resp.Data))
	for _, item := range resp.Data {
		code := normalizeHotRankCode(item.SC)
		codes = append(codes, code)
	}

	// Batch resolve names via ulist API (one request instead of N)
	secids := make([]string, 0, len(codes))
	for _, code := range codes {
		secids = append(secids, client.RuleBasedSecID(code))
	}
	nameMap := make(map[string]string)
	if len(secids) > 0 {
		secidStr := strings.Join(secids, ",")
		fields := "f12,f14"
		items, err := client.FetchQuotesBySecIDs(secidStr, fields)
		if err == nil {
			for _, item := range items {
				if code, ok := item["f12"]; ok {
					if name, ok2 := item["f14"]; ok2 {
						nameMap[fmt.Sprintf("%v", code)] = fmt.Sprintf("%v", name)
					}
				}
			}
		}
	}

	for i, item := range resp.Data {
		r := make(model.Record)
		r["rank"] = item.RK
		code := codes[i]
		r["code"] = code
		if name, ok := nameMap[code]; ok {
			r["name"] = name
		}
		records = append(records, r)
	}
	return records, nil
}

func GetHotRankHistory(code string) ([]model.Record, error) {
	payload := map[string]any{
		"appId":     "appId01",
		"globalId":  "786e4c21-70dc-435a-93bb-38",
		"stockCode": code,
		"pageNo":    1,
		"pageSize":  30,
	}
	data, err := client.DefaultClient.PostJSON(hotRankBaseURL+"getHisList", payload)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", code, "hot_rank_history", err, nil)
	}

	var resp struct {
		Data []struct {
			Rank int    `json:"rank"`
			Date string `json:"date"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", code, "hot_rank_history", err, nil)
	}

	records := make([]model.Record, 0, len(resp.Data))
	for _, item := range resp.Data {
		r := make(model.Record)
		r["date"] = item.Date
		r["rank"] = item.Rank
		r["code"] = code
		records = append(records, r)
	}
	return records, nil
}

func GetHotRankKeyword(code string) ([]model.Record, error) {
	payload := map[string]any{
		"appId":    "appId01",
		"globalId": "786e4c21-70dc-435a-93bb-38",
		"code":     code,
		"pageNo":   1,
		"pageSize": 20,
	}
	data, err := client.DefaultClient.PostJSON(hotRankBaseURL+"getHotStockRankList", payload)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", code, "hot_rank_keyword", err, nil)
	}

	var resp struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", code, "hot_rank_keyword", err, nil)
	}

	records := make([]model.Record, 0, len(resp.Data))
	for _, item := range resp.Data {
		r := make(model.Record)
		if v, ok := item["concept_name"]; ok {
			r["concept_name"] = v
		}
		if v, ok := item["concept_code"]; ok {
			r["concept_code"] = v
		}
		if v, ok := item["hot_value"]; ok {
			r["hot_value"] = v
		}
		r["code"] = code
		records = append(records, r)
	}
	return records, nil
}

func normalizeHotRankCode(sc string) string {
	if strings.HasPrefix(sc, "SH") {
		return sc[2:]
	}
	if strings.HasPrefix(sc, "SZ") {
		return sc[2:]
	}
	return sc
}

func extractDatacenterFieldNames(reportName string) ([]string, error) {
	resp, err := client.FetchDatacenterReport(reportName, nil, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, nil
	}
	keys := make([]string, 0, len(resp.Data[0]))
	for k := range resp.Data[0] {
		keys = append(keys, k)
	}
	return keys, nil
}

func buildFilter(parts ...string) string {
	return "(" + strings.Join(parts, "") + ")"
}

func addDateFilter(base, start, end string) string {
	return addDateFilterCol("TRADE_DATE", base, start, end)
}

func addDateFilterCol(col, base, start, end string) string {
	var parts []string
	if base != "" {
		parts = append(parts, base)
	}
	if start != "" {
		parts = append(parts, fmt.Sprintf("(%s>='%s')", col, start))
	}
	if end != "" {
		parts = append(parts, fmt.Sprintf("(%s<='%s')", col, end))
	}
	return "(" + strings.Join(parts, "") + ")"
}

func datacenterFilter(code, field string) string {
	return fmt.Sprintf("(%s=\"%s\")", field, code)
}

func fetchDatacenter(reportName string, filters url.Values, limit int) ([]map[string]interface{}, error) {
	return client.FetchDatacenterReportAll(reportName, filters, limit)
}

func parseDatacenterToRecords(data []map[string]interface{}, fieldMap map[string]string) []model.Record {
	return mapDatacenterRecords(data, fieldMap)
}
