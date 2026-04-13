package stock

import (
	"strings"

	"github.com/walnut1024/efi-cli/internal/analytics"
	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var historySchema = cliutil.CommandSchema{
	Command: "stock history",
	Entity:  "stock_kline",
	Supports: map[string]interface{}{
		"format":      []string{"csv", "json", "table"},
		"raw":         true,
		"schema":      true,
		"list_fields": true,
		"indicators":  []string{"ma", "ema", "macd", "rsi", "boll"},
		"stats":       []string{"total_return", "period_return", "cumulative_pct", "amplitude_avg", "high", "low", "max_drawdown", "start_close", "end_close", "start_date", "end_date", "bars"},
		"summary":     true,
	},
	DefaultFields: model.DefaultKlineOutputFields,
	Fields: []cliutil.FieldSchema{
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "date", Type: "string", Desc: "日期"},
		{Name: "open", Type: "number", Desc: "开盘价"},
		{Name: "close", Type: "number", Desc: "收盘价"},
		{Name: "high", Type: "number", Desc: "最高价"},
		{Name: "low", Type: "number", Desc: "最低价"},
		{Name: "vol", Type: "number", Desc: "成交量"},
		{Name: "amount", Type: "number", Desc: "成交额"},
		{Name: "amplitude", Type: "number", Desc: "振幅"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
		{Name: "chg", Type: "number", Desc: "涨跌额"},
		{Name: "turnover", Type: "number", Desc: "换手率"},
		{Name: "cumulative_pct", Type: "number", Desc: "相对区间起点的累计涨跌幅"},
		{Name: "total_return", Type: "number", Desc: "区间总收益率"},
		{Name: "period_return", Type: "number", Desc: "区间收益率"},
		{Name: "amplitude_avg", Type: "number", Desc: "区间平均振幅"},
		{Name: "max_drawdown", Type: "number", Desc: "区间最大回撤"},
		{Name: "start_date", Type: "string", Desc: "区间起始日期"},
		{Name: "end_date", Type: "string", Desc: "区间结束日期"},
		{Name: "start_close", Type: "number", Desc: "区间起始收盘价"},
		{Name: "end_close", Type: "number", Desc: "区间结束收盘价"},
		{Name: "bars", Type: "number", Desc: "K 线条数"},
		{Name: "ma5", Type: "number", Desc: "5 日均线示例字段"},
		{Name: "ema12", Type: "number", Desc: "12 日指数均线示例字段"},
		{Name: "dif", Type: "number", Desc: "MACD DIF"},
		{Name: "dea", Type: "number", Desc: "MACD DEA"},
		{Name: "macd", Type: "number", Desc: "MACD 柱值"},
		{Name: "rsi14", Type: "number", Desc: "14 周期 RSI 示例字段"},
		{Name: "boll_mid", Type: "number", Desc: "布林中轨"},
		{Name: "boll_up", Type: "number", Desc: "布林上轨"},
		{Name: "boll_low", Type: "number", Desc: "布林下轨"},
	},
}

func GetHistory(code, beg, end string, klt, fqt int, opts analytics.HistoryAnalysisOptions) ([]model.Record, error) {
	records, err := fetchResolvedKlines(code, beg, end, klt, fqt, "stock_history", client.ResolveQuoteID)
	if err != nil {
		return nil, err
	}
	return analytics.AnalyzeHistory(records, opts)
}

func GetHistoryRaw(code, beg, end string, klt, fqt int) ([]byte, error) {
	return fetchResolvedKlinesRaw(code, beg, end, klt, fqt, "stock_history", client.ResolveQuoteID)
}

func parseCSVFlagValues(values []string) []string {
	result := make([]string, 0)
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
	}
	return result
}
