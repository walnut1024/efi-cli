package model

import "strings"

func ParseKlineRecords(name, code string, klines []string) []Record {
	records := make([]Record, 0, len(klines))
	for _, line := range klines {
		parts := strings.Split(line, ",")
		r := make(Record)
		if name != "" {
			r["name"] = name
		}
		if code != "" {
			r["code"] = code
		}
		for i, fk := range KlineFieldKeys {
			if i < len(parts) {
				r[KlineFields[fk]] = CleanValue(parts[i])
			}
		}
		records = append(records, r)
	}
	return records
}
