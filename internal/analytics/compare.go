package analytics

import (
	"fmt"
	"sort"

	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

func BuildAlignedComparison(left, right []model.Record) ([]model.Record, error) {
	leftMap, leftBase, err := buildSeriesMap(left)
	if err != nil {
		return nil, err
	}
	rightMap, rightBase, err := buildSeriesMap(right)
	if err != nil {
		return nil, err
	}

	commonDates := make([]string, 0)
	for date := range leftMap {
		if _, ok := rightMap[date]; ok {
			commonDates = append(commonDates, date)
		}
	}
	sort.Strings(commonDates)

	results := make([]model.Record, 0, len(commonDates))
	for _, date := range commonDates {
		l := leftMap[date]
		r := rightMap[date]
		leftClose, lok := toFloat(l["close"])
		rightClose, rok := toFloat(r["close"])
		if !lok || !rok {
			continue
		}

		leftPct := 0.0
		rightPct := 0.0
		if leftBase != 0 {
			leftPct = (leftClose/leftBase - 1) * 100
		}
		if rightBase != 0 {
			rightPct = (rightClose/rightBase - 1) * 100
		}

		results = append(results, model.Record{
			"date":                 date,
			"left_code":            l["code"],
			"left_name":            l["name"],
			"left_close":           leftClose,
			"left_cumulative_pct":  leftPct,
			"right_code":           r["code"],
			"right_name":           r["name"],
			"right_close":          rightClose,
			"right_cumulative_pct": rightPct,
			"spread_pct":           leftPct - rightPct,
		})
	}
	return results, nil
}

func buildSeriesMap(records []model.Record) (map[string]model.Record, float64, error) {
	if len(records) == 0 {
		return map[string]model.Record{}, 0, nil
	}
	baseClose, ok := toFloat(records[0]["close"])
	if !ok || baseClose == 0 {
		return nil, 0, cliutil.NewError(cliutil.ErrInvalidArg, "参数错误", "compare", "align_date", fmt.Errorf("missing close data"), nil)
	}
	result := make(map[string]model.Record, len(records))
	for _, record := range records {
		date := asString(record["date"])
		if date == "" {
			continue
		}
		result[date] = record
	}
	return result, baseClose, nil
}
