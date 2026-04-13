package model

import "testing"

func TestParseQuoteRecordsIncludesSecID(t *testing.T) {
	records := ParseQuoteRecords([]map[string]interface{}{
		{
			"f12": "600000",
			"f13": 1,
			"f14": "浦发银行",
			"f2":  10.5,
			"f3":  1.2,
		},
	})

	if len(records) != 1 {
		t.Fatalf("unexpected record count: %d", len(records))
	}
	record := records[0]
	if got := record["code"]; got != "600000" {
		t.Fatalf("unexpected code: %v", got)
	}
	if got := record["secid"]; got != "1.600000" {
		t.Fatalf("unexpected secid: %v", got)
	}
	if got := record["price"]; got != 10.5 {
		t.Fatalf("unexpected price: %v", got)
	}
}
