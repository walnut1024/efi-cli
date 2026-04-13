package analytics

import (
	"testing"

	"github.com/walnut1024/efi-cli/internal/model"
)

func TestBuildAlignedComparison(t *testing.T) {
	left := []model.Record{
		{"code": "AAA", "name": "左侧", "date": "2024-01-01", "close": 10.0},
		{"code": "AAA", "name": "左侧", "date": "2024-01-02", "close": 11.0},
		{"code": "AAA", "name": "左侧", "date": "2024-01-03", "close": 12.0},
	}
	right := []model.Record{
		{"code": "BBB", "name": "右侧", "date": "2024-01-02", "close": 20.0},
		{"code": "BBB", "name": "右侧", "date": "2024-01-03", "close": 22.0},
		{"code": "BBB", "name": "右侧", "date": "2024-01-04", "close": 24.0},
	}

	got, err := BuildAlignedComparison(left, right)
	if err != nil {
		t.Fatalf("BuildAlignedComparison returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("unexpected aligned rows: %d", len(got))
	}
	if got[0]["date"] != "2024-01-02" {
		t.Fatalf("unexpected first date: %v", got[0]["date"])
	}
	if !almostEqual(got[1]["left_cumulative_pct"], 20.0) {
		t.Fatalf("unexpected left cumulative pct: %v", got[1]["left_cumulative_pct"])
	}
	if !almostEqual(got[1]["right_cumulative_pct"], 10.0) {
		t.Fatalf("unexpected right cumulative pct: %v", got[1]["right_cumulative_pct"])
	}
}
