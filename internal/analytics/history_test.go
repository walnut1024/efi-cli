package analytics

import (
	"math"
	"testing"

	"github.com/walnut1024/efi-cli/internal/model"
)

func TestAnalyzeHistoryAddsMAAndCumulativePct(t *testing.T) {
	records := []model.Record{
		{"code": "000001", "name": "测试", "date": "2024-01-01", "close": 10.0, "high": 10.0, "low": 9.0, "amplitude": 1.0},
		{"code": "000001", "name": "测试", "date": "2024-01-02", "close": 11.0, "high": 11.0, "low": 10.0, "amplitude": 2.0},
		{"code": "000001", "name": "测试", "date": "2024-01-03", "close": 12.0, "high": 12.0, "low": 11.0, "amplitude": 3.0},
	}

	got, err := AnalyzeHistory(records, HistoryAnalysisOptions{Indicators: []string{"ma:2"}})
	if err != nil {
		t.Fatalf("AnalyzeHistory returned error: %v", err)
	}

	if got[1]["ma2"] != 10.5 {
		t.Fatalf("unexpected ma2: %v", got[1]["ma2"])
	}
	if !almostEqual(got[2]["cumulative_pct"], 20.0) {
		t.Fatalf("unexpected cumulative_pct: %v", got[2]["cumulative_pct"])
	}
}

func TestAnalyzeHistorySummaryStats(t *testing.T) {
	records := []model.Record{
		{"code": "000001", "name": "测试", "date": "2024-01-01", "close": 10.0, "high": 10.0, "low": 9.0, "amplitude": 1.0},
		{"code": "000001", "name": "测试", "date": "2024-01-02", "close": 8.0, "high": 8.5, "low": 7.5, "amplitude": 2.0},
		{"code": "000001", "name": "测试", "date": "2024-01-03", "close": 12.0, "high": 12.5, "low": 11.0, "amplitude": 3.0},
	}

	got, err := AnalyzeHistory(records, HistoryAnalysisOptions{
		Stats:   []string{"total_return", "high", "low", "amplitude_avg", "max_drawdown"},
		Summary: true,
	})
	if err != nil {
		t.Fatalf("AnalyzeHistory returned error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("unexpected summary length: %d", len(got))
	}
	if !almostEqual(got[0]["total_return"], 20.0) {
		t.Fatalf("unexpected total_return: %v", got[0]["total_return"])
	}
	if got[0]["high"] != 12.5 {
		t.Fatalf("unexpected high: %v", got[0]["high"])
	}
	if got[0]["low"] != 7.5 {
		t.Fatalf("unexpected low: %v", got[0]["low"])
	}
	if !almostEqual(got[0]["amplitude_avg"], 2.0) {
		t.Fatalf("unexpected amplitude_avg: %v", got[0]["amplitude_avg"])
	}
	if !almostEqual(got[0]["max_drawdown"], -20.0) {
		t.Fatalf("unexpected max_drawdown: %v", got[0]["max_drawdown"])
	}
}

func TestParseIndicatorRequestsRejectsUnsupportedIndicator(t *testing.T) {
	_, err := ParseIndicatorRequests([]string{"foo"})
	if err == nil {
		t.Fatal("expected error for unsupported indicator")
	}
}

func almostEqual(v interface{}, target float64) bool {
	n, ok := v.(float64)
	if !ok {
		return false
	}
	return math.Abs(n-target) < 1e-9
}
