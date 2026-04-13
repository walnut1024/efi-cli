package model

import "testing"

func TestParseKlineRecords(t *testing.T) {
	records := ParseKlineRecords("测试指数", "000001", []string{
		"2024-01-02,10.1,10.5,10.8,9.9,100,1000,2.1,3.5,0.4,1.2",
	})

	if len(records) != 1 {
		t.Fatalf("unexpected record count: %d", len(records))
	}

	record := records[0]
	if got := record["name"]; got != "测试指数" {
		t.Fatalf("unexpected name: %v", got)
	}
	if got := record["code"]; got != "000001" {
		t.Fatalf("unexpected code: %v", got)
	}
	if got := record["close"]; got != 10.5 {
		t.Fatalf("unexpected close: %v", got)
	}
	if got := record["vol"]; got != int64(100) {
		t.Fatalf("unexpected vol: %v", got)
	}
}
