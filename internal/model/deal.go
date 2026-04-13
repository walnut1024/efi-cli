package model

import "strings"

func ParseDealRecords(code string, details []string, preClose interface{}) []Record {
	records := make([]Record, 0, len(details))
	for _, line := range details {
		parts := strings.Split(line, ",")
		r := make(Record)
		r["code"] = code
		if len(parts) > 0 {
			r["time"] = parts[0]
		}
		if len(parts) > 1 {
			r["price"] = CleanValue(parts[1])
		}
		if len(parts) > 2 {
			r["vol"] = CleanValue(parts[2])
		}
		if len(parts) > 3 {
			r["num"] = CleanValue(parts[3])
		}
		r["pre_close"] = preClose
		records = append(records, r)
	}
	return records
}
