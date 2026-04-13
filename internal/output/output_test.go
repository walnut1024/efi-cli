package output

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/walnut1024/efi-cli/internal/model"
)

func TestOrderedKeysUsesExplicitFieldOrderWithoutGlobalState(t *testing.T) {
	record := model.Record{
		"code": "600000",
		"name": "浦发银行",
		"pct":  1.2,
	}

	keys := orderedKeys(record, []string{"name", "code"})
	want := []string{"name", "code", "pct"}
	if !reflect.DeepEqual(keys, want) {
		t.Fatalf("unexpected keys: got %v want %v", keys, want)
	}
}

func TestOrderedKeysSortsWhenFieldOrderEmpty(t *testing.T) {
	record := model.Record{
		"name": "浦发银行",
		"code": "600000",
	}

	keys := orderedKeys(record, nil)
	want := []string{"code", "name"}
	if !reflect.DeepEqual(keys, want) {
		t.Fatalf("unexpected keys: got %v want %v", keys, want)
	}
}

func TestOrderedKeysUsesFriendlyQuoteOrderByDefault(t *testing.T) {
	record := model.Record{
		"pct":      1.2,
		"price":    10.5,
		"name":     "浦发银行",
		"code":     "600000",
		"turnover": 0.8,
	}

	keys := orderedKeys(record, nil)
	want := []string{"code", "name", "price", "pct", "turnover"}
	if !reflect.DeepEqual(keys, want) {
		t.Fatalf("unexpected keys: got %v want %v", keys, want)
	}
}

func TestOrderedKeysUsesFriendlyKlineOrderByDefault(t *testing.T) {
	record := model.Record{
		"name":   "浦发银行",
		"code":   "600000",
		"high":   10.8,
		"close":  10.5,
		"date":   "2024-01-02",
		"open":   10.1,
		"low":    9.9,
		"amount": 1000,
	}

	keys := orderedKeys(record, nil)
	want := []string{"name", "code", "date", "open", "close", "high", "low", "amount"}
	if !reflect.DeepEqual(keys, want) {
		t.Fatalf("unexpected keys: got %v want %v", keys, want)
	}
}

func TestRenderTableUsesCompactStyle(t *testing.T) {
	records := []model.Record{
		{"code": "600000", "name": "浦发银行", "price": 10.52, "pct": 1.23},
	}

	var buf bytes.Buffer
	if err := renderTable(&buf, records, []string{"code", "name", "price", "pct"}, false); err != nil {
		t.Fatalf("renderTable returned error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "+---") {
		t.Fatalf("expected compact style, got boxed table: %q", out)
	}
	if !strings.Contains(strings.ToLower(out), "code") || !strings.Contains(out, "浦发银行") {
		t.Fatalf("unexpected table output: %q", out)
	}
}

func TestSelectRecordsAppliesFieldFilterAndLimit(t *testing.T) {
	records := []model.Record{
		{"code": "600000", "name": "浦发银行", "pct": 1.2},
		{"code": "000001", "name": "平安银行", "pct": 0.8},
	}

	got, order := SelectRecords(records, "name,code", 1)
	if !reflect.DeepEqual(order, []string{"name", "code"}) {
		t.Fatalf("unexpected order: %v", order)
	}
	if len(got) != 1 {
		t.Fatalf("unexpected record count: %d", len(got))
	}
	if _, ok := got[0]["pct"]; ok {
		t.Fatalf("unexpected filtered field pct present: %v", got[0])
	}
	if got[0]["name"] != "浦发银行" || got[0]["code"] != "600000" {
		t.Fatalf("unexpected record: %v", got[0])
	}
}
