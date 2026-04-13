package model

import "testing"

func TestParseDealRecords(t *testing.T) {
	records := ParseDealRecords("600000", []string{
		"09:30,10.50,100,3",
	}, 10.2)

	if len(records) != 1 {
		t.Fatalf("unexpected record count: %d", len(records))
	}
	record := records[0]
	if got := record["time"]; got != "09:30" {
		t.Fatalf("unexpected time: %v", got)
	}
	if got := record["price"]; got != 10.5 {
		t.Fatalf("unexpected price: %v", got)
	}
	if got := record["vol"]; got != int64(100) {
		t.Fatalf("unexpected vol: %v", got)
	}
	if got := record["pre_close"]; got != 10.2 {
		t.Fatalf("unexpected pre_close: %v", got)
	}
}
