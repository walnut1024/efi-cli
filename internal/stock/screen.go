package stock

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

type ScreenFilter struct {
	PEMax  float64
	PEMin  float64
	PBMax  float64
	PBMin  float64
	ROEMin float64
	CapMin float64 // 亿元
	CapMax float64
	PctMin float64
	PctMax float64
}

func GetScreen(market string, f ScreenFilter, limit int) ([]model.Record, error) {
	fs := resolveFS(market)
	if fs == "" {
		fs = resolveFS("沪深京A股")
	}

	// Fields for screening: quote fields + f23(PB) + f173(ROE)
	fields := strings.Join(model.QuoteFieldKeys, ",") + ",f23,f173"

	allRecords := make([]model.Record, 0)
	pageSize := 500

	// Fetch first page to get total count
	firstResp, err := client.FetchQuoteList(fs, fields, 1, pageSize)
	if err != nil {
		return nil, err
	}
	for _, item := range firstResp.Diff {
		r := parseScreenRecord(item)
		if matchScreenFilter(r, f) {
			allRecords = append(allRecords, r)
		}
	}

	total := firstResp.Total
	remaining := (total + pageSize - 1) / pageSize // total pages
	if remaining > 1 {
		// Fetch remaining pages concurrently
		type pageResult struct {
			records []model.Record
			err     error
		}
		results := make([]pageResult, remaining-1)
		var wg sync.WaitGroup
		sem := make(chan struct{}, client.MaxConcurrency)
		for p := 2; p <= remaining; p++ {
			wg.Add(1)
			go func(page int) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				resp, fetchErr := client.FetchQuoteList(fs, fields, page, pageSize)
				if fetchErr != nil {
					results[page-2] = pageResult{err: fetchErr}
					return
				}
				var recs []model.Record
				for _, item := range resp.Diff {
					r := parseScreenRecord(item)
					if matchScreenFilter(r, f) {
						recs = append(recs, r)
					}
				}
				results[page-2] = pageResult{records: recs}
			}(p)
		}
		wg.Wait()
		for _, res := range results {
			if res.err != nil {
				return nil, res.err
			}
			allRecords = append(allRecords, res.records...)
		}
	}

	// Sort by pct descending
	sort.Slice(allRecords, func(i, j int) bool {
		vi, _ := allRecords[i]["pct"].(float64)
		vj, _ := allRecords[j]["pct"].(float64)
		return vi > vj
	})

	if limit > 0 && len(allRecords) > limit {
		allRecords = allRecords[:limit]
	}

	return allRecords, nil
}

func parseScreenRecord(item map[string]interface{}) model.Record {
	r := make(model.Record)
	// Quote fields
	for fk, ok := range model.QuoteFields {
		if v, exists := item[fk]; exists {
			r[ok] = v
		}
	}
	// Extra fields
	if v, ok := item["f23"]; ok {
		r["pb"] = v
	}
	if v, ok := item["f173"]; ok {
		r["roe"] = v
	}
	if mktNum, ok := item["f13"]; ok {
		if code, ok2 := item["f12"]; ok2 {
			r["secid"] = fmt.Sprintf("%v.%v", mktNum, code)
		}
	}
	// Convert market cap from yuan to 亿元
	if v, ok := r["mkt_cap"]; ok {
		switch val := v.(type) {
		case float64:
			r["cap"] = val / 1e8
		case int64:
			r["cap"] = float64(val) / 1e8
		}
	}
	return r
}

func matchScreenFilter(r model.Record, f ScreenFilter) bool {
	if f.PEMin > 0 || f.PEMax > 0 {
		v := toFloat(r["pe"])
		if v == nil {
			return false
		}
		if f.PEMin > 0 && *v < f.PEMin {
			return false
		}
		if f.PEMax > 0 && *v > f.PEMax {
			return false
		}
	}
	if f.PBMin > 0 || f.PBMax > 0 {
		v := toFloat(r["pb"])
		if v == nil {
			return false
		}
		if f.PBMin > 0 && *v < f.PBMin {
			return false
		}
		if f.PBMax > 0 && *v > f.PBMax {
			return false
		}
	}
	if f.ROEMin > 0 {
		v := toFloat(r["roe"])
		if v == nil || *v < f.ROEMin {
			return false
		}
	}
	if f.CapMin > 0 || f.CapMax > 0 {
		v := toFloat(r["cap"])
		if v == nil {
			return false
		}
		if f.CapMin > 0 && *v < f.CapMin {
			return false
		}
		if f.CapMax > 0 && *v > f.CapMax {
			return false
		}
	}
	if f.PctMin != 0 || f.PctMax != 0 {
		v := toFloat(r["pct"])
		if v == nil {
			return false
		}
		if f.PctMin != 0 && *v < f.PctMin {
			return false
		}
		if f.PctMax != 0 && *v > f.PctMax {
			return false
		}
	}
	return true
}

func toFloat(v interface{}) *float64 {
	switch val := v.(type) {
	case float64:
		return &val
	case int64:
		f := float64(val)
		return &f
	case json.Number:
		f, _ := val.Float64()
		return &f
	default:
		return nil
	}
}
