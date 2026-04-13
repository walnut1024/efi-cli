package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/walnut1024/efi-cli/internal/model"
)

func Format(records []model.Record, format, fieldsStr string, limit int, noHeader bool) error {
	records, fieldOrder := SelectRecords(records, fieldsStr, limit)

	switch format {
	case "json":
		return formatJSON(records)
	case "csv":
		return formatCSV(records, noHeader, fieldOrder)
	case "table":
		return formatTable(records, noHeader, fieldOrder)
	default:
		return formatCSV(records, noHeader, fieldOrder)
	}
}

func SelectRecords(records []model.Record, fieldsStr string, limit int) ([]model.Record, []string) {
	if limit > 0 && len(records) > limit {
		records = records[:limit]
	}

	var fieldOrder []string
	if fieldsStr != "" {
		fieldOrder = parseFieldOrder(fieldsStr)
		records = filterFields(records, fieldOrder)
	}
	return records, fieldOrder
}

func parseFieldOrder(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, len(parts))
	for i, p := range parts {
		result[i] = strings.TrimSpace(p)
	}
	return result
}

func filterFields(records []model.Record, order []string) []model.Record {
	fieldSet := make(map[string]bool, len(order))
	for _, f := range order {
		fieldSet[f] = true
	}

	filtered := make([]model.Record, 0, len(records))
	for _, r := range records {
		nr := make(model.Record, len(order))
		for k, v := range r {
			if fieldSet[k] {
				nr[k] = v
			}
		}
		filtered = append(filtered, nr)
	}
	return filtered
}

func formatJSON(records []model.Record) error {
	return WriteJSON(os.Stdout, records, false)
}

func WriteJSON(w io.Writer, v interface{}, pretty bool) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(v)
}

func formatCSV(records []model.Record, noHeader bool, fieldOrder []string) error {
	if len(records) == 0 {
		return nil
	}
	w := csv.NewWriter(os.Stdout)
	keys := orderedKeys(records[0], fieldOrder)

	if !noHeader {
		if err := w.Write(keys); err != nil {
			return err
		}
	}
	for _, r := range records {
		row := make([]string, len(keys))
		for i, k := range keys {
			if v, ok := r[k]; ok {
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func formatTable(records []model.Record, noHeader bool, fieldOrder []string) error {
	if len(records) == 0 {
		return nil
	}
	keys := orderedKeys(records[0], fieldOrder)
	return renderTable(os.Stdout, records, keys, noHeader)
}

func orderedKeys(r model.Record, fieldOrder []string) []string {
	if len(fieldOrder) > 0 {
		keys := make([]string, 0, len(r))
		for _, k := range fieldOrder {
			if _, ok := r[k]; ok {
				keys = append(keys, k)
			}
		}
		// Add any remaining keys not in fieldOrder
		for k := range r {
			found := false
			for _, f := range fieldOrder {
				if f == k {
					found = true
					break
				}
			}
			if !found {
				keys = append(keys, k)
			}
		}
		return keys
	}

	defaultOrder := inferDefaultFieldOrder(r)
	if len(defaultOrder) > 0 {
		keys := make([]string, 0, len(r))
		seen := make(map[string]struct{}, len(r))
		for _, k := range defaultOrder {
			if _, ok := r[k]; ok {
				keys = append(keys, k)
				seen[k] = struct{}{}
			}
		}
		rest := make([]string, 0, len(r)-len(keys))
		for k := range r {
			if _, ok := seen[k]; !ok {
				rest = append(rest, k)
			}
		}
		sort.Strings(rest)
		keys = append(keys, rest...)
		return keys
	}

	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func inferDefaultFieldOrder(r model.Record) []string {
	switch {
	case hasKeys(r, "date", "open", "close", "high", "low"):
		return model.DefaultKlineOutputFields
	case hasKeys(r, "date", "main_net", "small_net", "med_net"):
		return model.DefaultBillOutputFields
	case hasKeys(r, "date", "fin_balance", "loan_balance"):
		return append([]string{}, marginPreferredOrder...)
	case hasKeys(r, "date", "fund_inflow", "net_deal_amt"):
		return append([]string{}, connectHistoryPreferredOrder...)
	case hasKeys(r, "date", "score", "rank"):
		return append([]string{}, commentPreferredOrder...)
	case hasKeys(r, "date", "stock_count", "free_shares"):
		return append([]string{}, restrictedPreferredOrder...)
	case hasKeys(r, "rank", "code"):
		return append([]string{}, hotRankPreferredOrder...)
	case hasKeys(r, "date", "mutual_type_name"):
		return append([]string{}, connectSummaryPreferredOrder...)
	case hasKeys(r, "time", "sh_connect"):
		return append([]string{}, connectRealtimePreferredOrder...)
	case hasKeys(r, "code", "name", "price", "pct", "seal_amount"):
		return append([]string{}, poolZTPreferredOrder...)
	case hasKeys(r, "code", "name", "price", "pct", "total_cap", "up_count"):
		return append([]string{}, sectorListPreferredOrder...)
	case hasKeys(r, "code", "name", "price", "pct", "amplitude"):
		return append([]string{}, sectorMemberPreferredOrder...)
	case hasKeys(r, "code", "name", "price", "pct"):
		return append([]string{}, quotePreferredOrder...)
	case hasKeys(r, "code", "name", "nav"):
		return append([]string{}, fundPreferredOrder...)
	case hasKeys(r, "code", "time", "price", "vol"):
		return append([]string{}, dealPreferredOrder...)
	case hasKeys(r, "code", "name", "pe", "pb"):
		return append([]string{}, baseInfoPreferredOrder...)
	default:
		return append([]string{}, commonPreferredOrder...)
	}
}

func hasKeys(r model.Record, keys ...string) bool {
	for _, key := range keys {
		if _, ok := r[key]; !ok {
			return false
		}
	}
	return true
}

var commonPreferredOrder = []string{
	"code", "name", "date", "time", "period",
}

var quotePreferredOrder = []string{
	"code", "name", "price", "pct", "chg",
	"open", "high", "low", "pre_close",
	"vol", "amount", "turnover", "vol_ratio",
	"pe", "pb", "roe",
	"mkt_cap", "float_cap",
	"industry", "board_name", "board_code",
	"latest_trade_date", "update_ts",
	"secid", "mkt_num",
}

var fundPreferredOrder = []string{
	"code", "name", "date", "nav_date", "estimate_time",
	"nav", "acc_nav", "estimate_pct", "pct",
	"week", "month", "m3", "m6", "y1", "y2", "y3", "y5",
	"rank", "total", "company", "estab_date",
	"stock_pct", "bond_pct", "cash_pct", "other_pct", "total_size",
	"fund_code", "stock_code", "stock_name", "change", "action", "sector",
}

var dealPreferredOrder = []string{
	"code", "time", "price", "vol", "num", "pre_close",
}

var baseInfoPreferredOrder = []string{
	"code", "name", "industry", "board_name", "board_code",
	"price", "pct", "pe", "pb", "roe",
	"net_profit", "gross_margin", "net_margin",
	"mkt_cap", "float_cap",
}

var marginPreferredOrder = []string{
	"date", "fin_balance", "loan_balance", "fin_buy_amt", "loan_sell_amt",
	"investor_num", "total_guarantee", "avg_guarantee_ratio", "index_close", "index_pct",
}

var connectHistoryPreferredOrder = []string{
	"date", "fund_inflow", "net_deal_amt", "buy_amt", "sell_amt",
	"hold_market_cap", "quota_balance", "lead_stock", "lead_code",
}

var connectSummaryPreferredOrder = []string{
	"date", "mutual_type_name", "direction", "board_type", "index_name",
}

var connectRealtimePreferredOrder = []string{
	"time", "sh_connect", "sz_connect", "north_total",
}

var commentPreferredOrder = []string{
	"code", "name", "price", "pct", "turnover", "pe",
	"score", "rank", "attention", "main_cost", "institution_participation", "rising", "date",
}

var hotRankPreferredOrder = []string{
	"rank", "code", "name", "price", "pct",
}

var restrictedPreferredOrder = []string{
	"date", "stock_count", "free_shares", "market_cap", "index_close", "index_pct",
}

var poolZTPreferredOrder = []string{
	"code", "name", "pct", "price", "amount", "float_cap", "total_cap",
	"turnover", "seal_amount", "consecutive", "break_count", "industry",
	"first_seal_time", "last_seal_time",
}

var sectorListPreferredOrder = []string{
	"code", "name", "pct", "chg", "price", "total_cap", "turnover",
	"up_count", "down_count", "lead_stock", "lead_code",
}

var sectorMemberPreferredOrder = []string{
	"code", "name", "pct", "price", "chg", "vol", "amt",
	"amplitude", "high", "low", "open", "pre_close", "turnover", "pe", "pb",
}

func renderTable(w io.Writer, records []model.Record, keys []string, noHeader bool) error {
	tw := table.NewWriter()

	style := table.StyleLight // value copy of the shared style
	style.Options.DrawBorder = false
	style.Options.SeparateRows = false
	style.Options.SeparateColumns = false
	style.Options.SeparateFooter = false
	style.Options.SeparateHeader = true
	box := style.Box // value copy
	box.PaddingLeft = " "
	box.PaddingRight = " "
	style.Box = box
	hdrFmt := style.Format // value copy
	hdrFmt.Header = text.FormatDefault
	style.Format = hdrFmt
	tw.SetStyle(style)

	if !noHeader {
		header := make(table.Row, len(keys))
		for i, key := range keys {
			header[i] = prettifyHeader(key)
		}
		tw.AppendHeader(header)
	}

	rows := make([]table.Row, 0, len(records))
	for _, record := range records {
		row := make(table.Row, len(keys))
		for i, key := range keys {
			row[i] = displayValue(record[key])
		}
		rows = append(rows, row)
	}
	tw.AppendRows(rows)
	tw.SetColumnConfigs(buildColumnConfigs(records, keys))
	_, err := fmt.Fprint(w, tw.Render())
	return err
}

func buildColumnConfigs(records []model.Record, keys []string) []table.ColumnConfig {
	configs := make([]table.ColumnConfig, 0, len(keys))
	for i, key := range keys {
		align := text.AlignLeft
		if isNumericColumn(records, key) {
			align = text.AlignRight
		}
		configs = append(configs, table.ColumnConfig{
			Number:      i + 1,
			Align:       align,
			AlignHeader: text.AlignLeft,
			WidthMax:    columnWidthLimit(key),
		})
	}
	return configs
}

func isNumericColumn(records []model.Record, key string) bool {
	seen := false
	for _, record := range records {
		v, ok := record[key]
		if !ok || v == nil {
			continue
		}
		seen = true
		switch v.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			continue
		default:
			return false
		}
	}
	return seen
}

func columnWidthLimit(key string) int {
	switch key {
	case "date":
		return 10
	case "time":
		return 8
	case "name", "industry", "board_name", "idx_name":
		return 18
	case "reason", "desc", "explain", "plan":
		return 28
	default:
		return 16
	}
}

func prettifyHeader(key string) string {
	if label, ok := headerLabels[key]; ok {
		return label
	}
	replacer := strings.NewReplacer("_", " ")
	return replacer.Replace(key)
}

var headerLabels = map[string]string{
	"code":                      "code",
	"name":                      "name",
	"date":                      "date",
	"time":                      "time",
	"open":                      "open",
	"close":                     "close",
	"high":                      "high",
	"low":                       "low",
	"price":                     "price",
	"pct":                       "pct%",
	"chg":                       "chg",
	"vol":                       "vol",
	"num":                       "trades",
	"amount":                    "amount",
	"turnover":                  "turn%",
	"amplitude":                 "amp%",
	"vol_ratio":                 "volr",
	"pre_close":                 "prev",
	"mkt_cap":                   "mkt cap",
	"float_cap":                 "float cap",
	"main_net":                  "main net",
	"small_net":                 "small net",
	"med_net":                   "med net",
	"big_net":                   "big net",
	"huge_net":                  "huge net",
	"main_pct":                  "main%",
	"small_pct":                 "small%",
	"med_pct":                   "med%",
	"big_pct":                   "big%",
	"huge_pct":                  "huge%",
	"nav":                       "nav",
	"acc_nav":                   "acc nav",
	"nav_date":                  "nav date",
	"estimate_time":             "est time",
	"estimate_pct":              "est%",
	"company":                   "company",
	"industry":                  "industry",
	"board_name":                "board",
	"board_code":                "board code",
	"latest_trade_date":         "trade date",
	"update_ts":                 "updated",
	"fin_balance":               "fin bal",
	"loan_balance":              "loan bal",
	"fin_buy_amt":               "fin buy",
	"loan_sell_amt":             "loan sell",
	"investor_num":              "investors",
	"total_guarantee":           "guarantee",
	"avg_guarantee_ratio":       "avg ratio",
	"fund_inflow":               "inflow",
	"net_deal_amt":              "net amt",
	"hold_market_cap":           "hold cap",
	"quota_balance":             "quota bal",
	"seal_amount":               "seal amt",
	"consecutive":               "cons",
	"break_count":               "breaks",
	"up_count":                  "up",
	"down_count":                "down",
	"lead_stock":                "leader",
	"lead_code":                 "lead code",
	"total_cap":                 "total cap",
	"score":                     "score",
	"attention":                 "attn",
	"main_cost":                 "main cost",
	"institution_participation": "inst part",
	"rank":                      "rank",
	"mutual_type_name":          "channel",
	"direction":                 "dir",
	"board_type":                "board",
	"sh_connect":                "sh conn",
	"sz_connect":                "sz conn",
	"north_total":               "north total",
	"stock_count":               "stocks",
	"free_shares":               "free shares",
	"free_date":                 "free date",
	"free_ratio":                "free ratio",
	"total_ratio":               "total ratio",
	"actual_free_shares":        "actual free",
	"batch_holder_num":          "holders",
}

func displayValue(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return "-"
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%.0f", val)
		}
		return fmt.Sprintf("%.2f", val)
	case float32:
		f := float64(val)
		if f == float64(int64(f)) {
			return fmt.Sprintf("%.0f", f)
		}
		return fmt.Sprintf("%.2f", f)
	default:
		return fmt.Sprintf("%v", val)
	}
}
