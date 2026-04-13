package searchcmd

import (
	"github.com/spf13/cobra"
	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var searchOutput cliutil.OutputConfig

var searchSchema = cliutil.CommandSchema{
	Command: "search",
	Entity:  "security_search_result",
	Supports: map[string]interface{}{
		"format":      []string{"csv", "json", "table"},
		"raw":         true,
		"schema":      true,
		"list_fields": true,
	},
	DefaultFields: []string{"code", "name", "market", "quote_id", "unified_code", "pinyin"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "market", Type: "string", Desc: "市场名称"},
		{Name: "mkt_num", Type: "string", Desc: "市场编号"},
		{Name: "quote_id", Type: "string", Desc: "secid"},
		{Name: "unified_code", Type: "string", Desc: "统一代码"},
		{Name: "pinyin", Type: "string", Desc: "拼音"},
		{Name: "inner_code", Type: "string", Desc: "内部代码"},
	},
}

func AddCommands(root *cobra.Command) {
	searchCmd := &cobra.Command{
		Use:   "search [keyword]",
		Short: "Search securities and return candidate list",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if searchOutput.HandleSchemaOrListFields(&searchSchema) {
				return
			}

			var (
				resp *client.SearchResponse
				err  error
			)
			if searchOutput.Raw {
				resp, err = client.SearchQuotesFresh(args[0], searchOutput.Limit)
			} else {
				resp, err = client.SearchQuotes(args[0], searchOutput.Limit)
			}
			if err != nil {
				cliutil.PrintErrorAndExit(err, searchOutput)
			}
			if resp == nil {
				searchOutput.FormatWithSchemaOrExit(nil, &searchSchema)
				return
			}
			if searchOutput.PrintRawOrExit(resp.Raw) {
				return
			}
			searchOutput.FormatWithSchemaOrExit(toRecords(resp.Candidates), &searchSchema)
		},
	}
	searchOutput.AddFlags(searchCmd)
	searchOutput.AddInspectFlags(searchCmd)
	root.AddCommand(searchCmd)
}

func toRecords(candidates []*client.Quote) []model.Record {
	records := make([]model.Record, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate == nil {
			continue
		}
		records = append(records, model.Record{
			"code":         candidate.Code,
			"name":         candidate.Name,
			"market":       model.GetMarketName(candidate.MktNum),
			"mkt_num":      candidate.MktNum,
			"quote_id":     candidate.QuoteID,
			"unified_code": candidate.UnifiedCode,
			"pinyin":       candidate.Pinyin,
			"inner_code":   candidate.InnerCode,
		})
	}
	return records
}
