package stock

import (
	"fmt"
	"strings"

	"github.com/walnut1024/efi-cli/internal/analytics"
	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/index"
	"github.com/walnut1024/efi-cli/internal/model"
)

var compareSchema = cliutil.CommandSchema{
	Command: "stock compare",
	Entity:  "stock_compare",
	Supports: map[string]interface{}{
		"format":      []string{"csv", "json", "table"},
		"schema":      true,
		"list_fields": true,
		"align_date":  true,
		"metric":      []string{"total_return", "period_return", "cumulative_pct", "amplitude_avg", "high", "low", "max_drawdown", "start_close", "end_close", "start_date", "end_date", "bars"},
	},
	DefaultFields: []string{"code", "name", "start_date", "end_date", "start_close", "end_close", "total_return", "max_drawdown", "high", "low", "bars"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "start_date", Type: "string", Desc: "区间起始日期"},
		{Name: "end_date", Type: "string", Desc: "区间结束日期"},
		{Name: "start_close", Type: "number", Desc: "区间起始收盘价"},
		{Name: "end_close", Type: "number", Desc: "区间结束收盘价"},
		{Name: "total_return", Type: "number", Desc: "区间总收益率"},
		{Name: "period_return", Type: "number", Desc: "区间收益率"},
		{Name: "cumulative_pct", Type: "number", Desc: "累计涨跌幅"},
		{Name: "amplitude_avg", Type: "number", Desc: "区间平均振幅"},
		{Name: "high", Type: "number", Desc: "区间最高价"},
		{Name: "low", Type: "number", Desc: "区间最低价"},
		{Name: "max_drawdown", Type: "number", Desc: "区间最大回撤"},
		{Name: "bars", Type: "number", Desc: "K 线数量"},
		{Name: "date", Type: "string", Desc: "对齐日期"},
		{Name: "left_code", Type: "string", Desc: "左侧标的代码"},
		{Name: "left_name", Type: "string", Desc: "左侧标的名称"},
		{Name: "left_close", Type: "number", Desc: "左侧收盘价"},
		{Name: "left_cumulative_pct", Type: "number", Desc: "左侧累计涨跌幅"},
		{Name: "right_code", Type: "string", Desc: "右侧标的代码"},
		{Name: "right_name", Type: "string", Desc: "右侧标的名称"},
		{Name: "right_close", Type: "number", Desc: "右侧收盘价"},
		{Name: "right_cumulative_pct", Type: "number", Desc: "右侧累计涨跌幅"},
		{Name: "spread_pct", Type: "number", Desc: "左右累计涨跌幅差值"},
	},
}

func CompareTargets(left, right, beg, end string, klt int, alignDate bool, metric string) ([]model.Record, error) {
	leftHistory, err := getComparableHistory(left, beg, end, klt)
	if err != nil {
		return nil, err
	}
	rightHistory, err := getComparableHistory(right, beg, end, klt)
	if err != nil {
		return nil, err
	}

	if alignDate {
		return analytics.BuildAlignedComparison(leftHistory, rightHistory)
	}

	stats := compareStats(metric)
	leftSummary, err := analytics.AnalyzeHistory(leftHistory, analytics.HistoryAnalysisOptions{
		Stats:   stats,
		Summary: true,
	})
	if err != nil {
		return nil, err
	}
	rightSummary, err := analytics.AnalyzeHistory(rightHistory, analytics.HistoryAnalysisOptions{
		Stats:   stats,
		Summary: true,
	})
	if err != nil {
		return nil, err
	}

	result := append([]model.Record{}, leftSummary...)
	result = append(result, rightSummary...)
	if metric != "" {
		cfg := cliutil.OutputConfig{Sort: metric}
		return cfg.ApplySort(result)
	}
	return result, nil
}

func compareStats(metric string) []string {
	base := []string{"start_date", "end_date", "start_close", "end_close", "total_return", "max_drawdown", "high", "low", "bars"}
	if strings.TrimSpace(metric) == "" {
		return base
	}
	seen := map[string]struct{}{}
	result := make([]string, 0, len(base)+1)
	for _, item := range append(base, metric) {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func getComparableHistory(target, beg, end string, klt int) ([]model.Record, error) {
	return fetchResolvedKlines(target, beg, end, klt, 1, "stock_compare", resolveCompareTarget)
}

func resolveCompareTarget(target string) (string, error) {
	if strings.Contains(target, ".") {
		return target, nil
	}
	if secid, ok := index.IndexSecID[target]; ok {
		return secid, nil
	}
	secid, err := client.ResolveQuoteID(target)
	if err == nil {
		return secid, nil
	}
	if cliErr, ok := cliutil.AsCLIError(err); ok && cliErr.Kind == cliutil.ErrCodeNotFound {
		secid = index.ResolveSecID(target)
		if secid != "" && secid != target {
			return secid, nil
		}
	}
	return "", err
}

func validateCompareMetric(metric string) error {
	if strings.TrimSpace(metric) == "" {
		return nil
	}
	for _, item := range compareSchema.Supports["metric"].([]string) {
		if item == metric {
			return nil
		}
	}
	return cliutil.NewError(cliutil.ErrInvalidArg, "参数错误", metric, "compare_metric", fmt.Errorf("unsupported metric"), nil)
}
