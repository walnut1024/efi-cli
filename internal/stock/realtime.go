package stock

import (
	"fmt"
	"sort"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var realtimeSchema = cliutil.CommandSchema{
	Command: "stock realtime",
	Entity:  "stock_realtime_list",
	Supports: map[string]interface{}{
		"format":      []string{"csv", "json", "table"},
		"raw":         true,
		"schema":      true,
		"list_fields": true,
		"sort":        true,
	},
	DefaultFields: quoteSchema.DefaultFields,
	Fields:        quoteSchema.Fields,
}

func GetRealtime(market string) ([]model.Record, error) {
	fs := resolveFS(market)
	if fs == "" {
		return nil, cliutil.NewError(cliutil.ErrInvalidArg, "参数错误", market, "stock_realtime", fmt.Errorf("unknown market"), nil)
	}

	// First page to get total
	firstResp, err := client.FetchQuoteList(fs, quoteFields, 1, 200)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", market, "stock_realtime", err, nil)
	}

	total := firstResp.Total
	pz := len(firstResp.Diff)
	if pz == 0 {
		return nil, nil
	}

	pages := total / pz
	if total%pz != 0 {
		pages++
	}

	allRecords := model.ParseQuoteRecords(firstResp.Diff)

	if pages > 1 {
		type result struct {
			records []model.Record
			err     error
			page    int
		}
		ch := make(chan result, pages-1)
		sem := make(chan struct{}, client.MaxConcurrency)
		pageResults := make([]result, 0, pages-1)

		for p := 2; p <= pages; p++ {
			go func(page int) {
				sem <- struct{}{}
				resp, e := client.FetchQuoteList(fs, quoteFields, page, pz)
				<-sem
				if e != nil {
					ch <- result{err: e, page: page}
					return
				}
				ch <- result{records: model.ParseQuoteRecords(resp.Diff), page: page}
			}(p)
		}

		for i := 0; i < pages-1; i++ {
			r := <-ch
			if r.err != nil {
				return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", market, "stock_realtime", fmt.Errorf("page %d: %w", r.page, r.err), nil)
			}

			pageResults = append(pageResults, r)
		}

		sort.Slice(pageResults, func(i, j int) bool {
			return pageResults[i].page < pageResults[j].page
		})
		for _, pageResult := range pageResults {
			allRecords = append(allRecords, pageResult.records...)
		}
	}

	return allRecords, nil
}

func GetRealtimeRaw(market string) ([]byte, error) {
	fs := resolveFS(market)
	if fs == "" {
		return nil, cliutil.NewError(cliutil.ErrInvalidArg, "参数错误", market, "stock_realtime", fmt.Errorf("unknown market"), nil)
	}
	_, raw, err := client.FetchQuoteListRaw(fs, quoteFields, 1, 200)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", market, "stock_realtime", err, nil)
	}
	return raw, nil
}

func resolveFS(market string) string {
	if market == "" {
		market = "沪深京A股"
	}
	if v, ok := model.FSDict[market]; ok {
		return v
	}
	return ""
}
