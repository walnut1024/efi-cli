package analytics

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

type IndicatorRequest struct {
	Name   string
	Params []int
}

type HistoryAnalysisOptions struct {
	Indicators []string
	Stats      []string
	Summary    bool
}

type historyBar struct {
	Name      string
	Code      string
	Date      string
	Close     float64
	High      float64
	Low       float64
	Amplitude float64
	Valid     bool
}

func AnalyzeHistory(records []model.Record, opts HistoryAnalysisOptions) ([]model.Record, error) {
	indicatorReqs, err := ParseIndicatorRequests(opts.Indicators)
	if err != nil {
		return nil, err
	}

	out := cloneRecords(records)
	bars := toHistoryBars(out)

	if len(indicatorReqs) > 0 || hasStat(opts.Stats, "cumulative_pct") {
		applyCumulativePct(out, bars)
	}

	if len(indicatorReqs) > 0 {
		for _, req := range indicatorReqs {
			switch req.Name {
			case "ma":
				applyMA(out, bars, req.Params)
			case "ema":
				applyEMA(out, bars, req.Params)
			case "macd":
				applyMACD(out, bars)
			case "rsi":
				applyRSI(out, bars, req.Params)
			case "boll":
				applyBOLL(out, bars, req.Params)
			}
		}
	}

	if len(opts.Stats) == 0 && !opts.Summary {
		return out, nil
	}

	summary, err := ComputeSummary(out, bars, opts.Stats)
	if err != nil {
		return nil, err
	}
	if opts.Summary {
		return []model.Record{summary}, nil
	}
	for i := range out {
		for key, value := range summary {
			if key == "name" || key == "code" {
				continue
			}
			if _, exists := out[i][key]; exists {
				continue
			}
			out[i][key] = value
		}
	}
	return out, nil
}

func ParseIndicatorRequests(specs []string) ([]IndicatorRequest, error) {
	reqs := make([]IndicatorRequest, 0)
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		name, params, err := parseIndicatorSpec(spec)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, IndicatorRequest{Name: name, Params: params})
	}
	return reqs, nil
}

func ComputeSummary(records []model.Record, bars []historyBar, stats []string) (model.Record, error) {
	if len(records) == 0 || len(bars) == 0 {
		return model.Record{}, nil
	}

	first := bars[0]
	last := bars[len(bars)-1]

	if !first.Valid || !last.Valid || first.Close == 0 {
		return nil, cliutil.NewError(cliutil.ErrInvalidArg, "参数错误", "history", "summary", fmt.Errorf("insufficient close data"), nil)
	}

	highest := first.High
	lowest := first.Low
	amplitudeSum := 0.0
	amplitudeCount := 0.0
	peak := first.Close
	maxDrawdown := 0.0

	for _, bar := range bars {
		if !bar.Valid {
			continue
		}
		if bar.High > highest {
			highest = bar.High
		}
		if lowest == 0 || (bar.Low > 0 && bar.Low < lowest) {
			lowest = bar.Low
		}
		if bar.Amplitude != 0 {
			amplitudeSum += bar.Amplitude
			amplitudeCount++
		}
		if bar.Close > peak {
			peak = bar.Close
		}
		if peak != 0 {
			drawdown := (bar.Close/peak - 1) * 100
			if drawdown < maxDrawdown {
				maxDrawdown = drawdown
			}
		}
	}

	totalReturn := (last.Close/first.Close - 1) * 100
	summary := model.Record{
		"name":           first.Name,
		"code":           first.Code,
		"start_date":     first.Date,
		"end_date":       last.Date,
		"start_close":    first.Close,
		"end_close":      last.Close,
		"total_return":   totalReturn,
		"cumulative_pct": totalReturn,
		"high":           highest,
		"low":            lowest,
		"bars":           len(records),
		"max_drawdown":   maxDrawdown,
	}
	if amplitudeCount > 0 {
		summary["amplitude_avg"] = amplitudeSum / amplitudeCount
	} else {
		summary["amplitude_avg"] = nil
	}

	if len(stats) == 0 {
		return summary, nil
	}

	filtered := model.Record{
		"name": first.Name,
		"code": first.Code,
	}
	for _, key := range stats {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if value, ok := summary[key]; ok {
			filtered[key] = value
			continue
		}
		return nil, cliutil.NewError(cliutil.ErrInvalidArg, "参数错误", key, "history_stats", fmt.Errorf("unsupported stat"), nil)
	}
	return filtered, nil
}

func applyCumulativePct(records []model.Record, bars []historyBar) {
	if len(records) == 0 || len(bars) == 0 || !bars[0].Valid || bars[0].Close == 0 {
		return
	}
	base := bars[0].Close
	for i, bar := range bars {
		if !bar.Valid {
			continue
		}
		records[i]["cumulative_pct"] = (bar.Close/base - 1) * 100
	}
}

func applyMA(records []model.Record, bars []historyBar, windows []int) {
	for _, window := range windows {
		if window <= 0 {
			continue
		}
		key := fmt.Sprintf("ma%d", window)
		sum := 0.0
		for i := range bars {
			if !bars[i].Valid {
				continue
			}
			sum += bars[i].Close
			if i >= window && bars[i-window].Valid {
				sum -= bars[i-window].Close
			}
			if i+1 >= window {
				records[i][key] = sum / float64(window)
			}
		}
	}
}

func applyEMA(records []model.Record, bars []historyBar, windows []int) {
	for _, window := range windows {
		if window <= 0 {
			continue
		}
		key := fmt.Sprintf("ema%d", window)
		multiplier := 2.0 / float64(window+1)
		var ema float64
		var initialized bool
		for i, bar := range bars {
			if !bar.Valid {
				continue
			}
			if !initialized {
				ema = bar.Close
				initialized = true
			} else {
				ema = (bar.Close-ema)*multiplier + ema
			}
			records[i][key] = ema
		}
	}
}

func applyMACD(records []model.Record, bars []historyBar) {
	var ema12 float64
	var ema26 float64
	var dea float64
	var init12 bool
	var init26 bool
	for i, bar := range bars {
		if !bar.Valid {
			continue
		}
		if !init12 {
			ema12 = bar.Close
			init12 = true
		} else {
			ema12 = (bar.Close-ema12)*(2.0/13.0) + ema12
		}
		if !init26 {
			ema26 = bar.Close
			init26 = true
		} else {
			ema26 = (bar.Close-ema26)*(2.0/27.0) + ema26
		}
		dif := ema12 - ema26
		if i == 0 {
			dea = dif
		} else {
			dea = (dif-dea)*(2.0/10.0) + dea
		}
		records[i]["dif"] = dif
		records[i]["dea"] = dea
		records[i]["macd"] = (dif - dea) * 2
	}
}

func applyRSI(records []model.Record, bars []historyBar, windows []int) {
	for _, window := range windows {
		if window <= 0 {
			continue
		}
		key := fmt.Sprintf("rsi%d", window)
		if len(bars) <= 1 {
			continue
		}
		gains := make([]float64, len(bars))
		losses := make([]float64, len(bars))
		for i := 1; i < len(bars); i++ {
			if !bars[i].Valid || !bars[i-1].Valid {
				continue
			}
			change := bars[i].Close - bars[i-1].Close
			if change > 0 {
				gains[i] = change
			} else {
				losses[i] = -change
			}
		}
		initialGain := 0.0
		initialLoss := 0.0
		var avgGain, avgLoss float64
		for i := 1; i < len(bars); i++ {
			if i <= window {
				initialGain += gains[i]
				initialLoss += losses[i]
			}
			if i < window {
				continue
			}
			if i == window {
				avgGain = initialGain / float64(window)
				avgLoss = initialLoss / float64(window)
			} else {
				avgGain = (avgGain*(float64(window)-1) + gains[i]) / float64(window)
				avgLoss = (avgLoss*(float64(window)-1) + losses[i]) / float64(window)
			}
			if avgLoss == 0 {
				records[i][key] = 100.0
				continue
			}
			rs := avgGain / avgLoss
			records[i][key] = 100 - (100 / (1 + rs))
		}
	}
}

func applyBOLL(records []model.Record, bars []historyBar, params []int) {
	window := 20
	multiplier := 2.0
	if len(params) > 0 && params[0] > 0 {
		window = params[0]
	}
	if len(params) > 1 && params[1] > 0 {
		multiplier = float64(params[1])
	}
	values := make([]float64, 0, window)
	for i, bar := range bars {
		if !bar.Valid {
			continue
		}
		values = append(values, bar.Close)
		if len(values) > window {
			values = values[1:]
		}
		if len(values) < window {
			continue
		}
		mid := average(values)
		std := stddev(values, mid)
		records[i]["boll_mid"] = mid
		records[i]["boll_up"] = mid + multiplier*std
		records[i]["boll_low"] = mid - multiplier*std
	}
}

func parseIndicatorSpec(spec string) (string, []int, error) {
	parts := strings.SplitN(strings.ToLower(spec), ":", 2)
	name := strings.TrimSpace(parts[0])
	switch name {
	case "ma":
		params, err := parseIntList(defaultWhenEmpty(getTail(parts), "5,10,20"))
		return name, params, err
	case "ema":
		params, err := parseIntList(defaultWhenEmpty(getTail(parts), "12,26"))
		return name, params, err
	case "macd":
		return name, nil, nil
	case "rsi":
		params, err := parseIntList(defaultWhenEmpty(getTail(parts), "14"))
		return name, params, err
	case "boll":
		params, err := parseIntList(defaultWhenEmpty(getTail(parts), "20,2"))
		return name, params, err
	default:
		return "", nil, cliutil.NewError(cliutil.ErrInvalidArg, "参数错误", spec, "history_indicators", fmt.Errorf("unsupported indicator"), nil)
	}
}

func parseIntList(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	values := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		v, err := strconv.Atoi(part)
		if err != nil {
			return nil, cliutil.NewError(cliutil.ErrInvalidArg, "参数错误", s, "history_indicators", err, nil)
		}
		values = append(values, v)
	}
	return values, nil
}

func defaultWhenEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func getTail(parts []string) string {
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func toHistoryBars(records []model.Record) []historyBar {
	bars := make([]historyBar, 0, len(records))
	for _, record := range records {
		closeValue, closeOK := toFloat(record["close"])
		highValue, _ := toFloat(record["high"])
		lowValue, _ := toFloat(record["low"])
		amplitudeValue, _ := toFloat(record["amplitude"])
		bars = append(bars, historyBar{
			Name:      asString(record["name"]),
			Code:      asString(record["code"]),
			Date:      asString(record["date"]),
			Close:     closeValue,
			High:      highValue,
			Low:       lowValue,
			Amplitude: amplitudeValue,
			Valid:     closeOK,
		})
	}
	return bars
}

func cloneRecords(records []model.Record) []model.Record {
	cloned := make([]model.Record, 0, len(records))
	for _, record := range records {
		item := make(model.Record, len(record))
		for k, v := range record {
			item[k] = v
		}
		cloned = append(cloned, item)
	}
	return cloned
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func stddev(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range values {
		diff := value - mean
		sum += diff * diff
	}
	return math.Sqrt(sum / float64(len(values)))
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(n), 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}

func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func hasStat(stats []string, target string) bool {
	for _, stat := range stats {
		if strings.TrimSpace(stat) == target {
			return true
		}
	}
	return false
}
