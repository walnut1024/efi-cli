package model

import "fmt"

func ParseQuoteRecords(items []map[string]interface{}) []Record {
	records := make([]Record, 0, len(items))
	for _, item := range items {
		r := make(Record)
		for fk, ok := range QuoteFields {
			if v, exists := item[fk]; exists {
				r[ok] = v
			}
		}
		if mktNum, ok := item["f13"]; ok {
			if code, ok2 := item["f12"]; ok2 {
				r["secid"] = fmt.Sprintf("%v.%v", mktNum, code)
			}
		}
		records = append(records, r)
	}
	return records
}
