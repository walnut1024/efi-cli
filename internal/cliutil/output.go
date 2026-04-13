package cliutil

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/walnut1024/efi-cli/internal/model"
	"github.com/walnut1024/efi-cli/internal/output"
)

type OutputConfig struct {
	Format     string
	Fields     string
	Limit      int
	NoHeader   bool
	Pretty     bool
	Compact    bool
	Schema     bool
	Raw        bool
	ListFields bool
	Sort       string
	Asc        bool
	Desc       bool
}

func (c *OutputConfig) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.Format, "format", "csv", "Output format: csv (default), json, table")
	cmd.Flags().StringVar(&c.Fields, "fields", "", "Comma-separated fields to output")
	cmd.Flags().IntVar(&c.Limit, "limit", 50, "Max rows to output")
	cmd.Flags().BoolVar(&c.NoHeader, "no-header", false, "Omit header row (csv/table)")
	cmd.Flags().BoolVar(&c.Pretty, "pretty", false, "Pretty-print JSON output")
	cmd.Flags().BoolVar(&c.Compact, "compact", false, "Compact JSON output")
}

func (c *OutputConfig) AddInspectFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&c.Schema, "schema", false, "Print command schema")
	cmd.Flags().BoolVar(&c.Raw, "raw", false, "Print raw upstream response")
	cmd.Flags().BoolVar(&c.ListFields, "list-fields", false, "List available fields")
}

func (c *OutputConfig) AddSortFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.Sort, "sort", "", "Sort by field")
	cmd.Flags().BoolVar(&c.Asc, "asc", false, "Sort ascending")
	cmd.Flags().BoolVar(&c.Desc, "desc", false, "Sort descending")
}

func (c OutputConfig) FormatOrExit(records []model.Record) {
	c.FormatWithSchemaOrExit(records, nil)
}

func (c OutputConfig) FormatWithSchemaOrExit(records []model.Record, schema *CommandSchema) {
	fields := c.Fields
	if schema != nil && strings.EqualFold(fields, "all") {
		fields = strings.Join(schema.FieldNames(), ",")
	}
	records, _ = output.SelectRecords(records, fields, c.Limit)

	if c.useJSONOutput() {
		if err := output.WriteJSON(os.Stdout, records, c.Pretty); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	if err := output.Format(records, c.Format, fields, c.Limit, c.NoHeader); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (c OutputConfig) HandleSchemaOrListFields(schema *CommandSchema) bool {
	if schema == nil {
		return false
	}
	if c.Schema {
		if err := output.WriteJSON(os.Stdout, schema, c.Pretty || !c.Compact); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return true
	}
	if c.ListFields {
		names := schema.FieldNames()
		if c.useJSONOutput() {
			if err := output.WriteJSON(os.Stdout, names, c.Pretty); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		} else {
			for _, name := range names {
				fmt.Fprintln(os.Stdout, name)
			}
		}
		return true
	}
	return false
}

func (c OutputConfig) PrintRawOrExit(raw []byte) bool {
	if !c.Raw {
		return false
	}
	if c.useJSONOutput() {
		var v interface{}
		if err := json.Unmarshal(raw, &v); err == nil {
			if err := output.WriteJSON(os.Stdout, v, c.Pretty || !c.Compact); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return true
		}
	}
	if _, err := os.Stdout.Write(raw); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(raw) == 0 || raw[len(raw)-1] != '\n' {
		fmt.Fprintln(os.Stdout)
	}
	return true
}

func (c OutputConfig) useJSONOutput() bool {
	return c.Pretty || c.Compact || c.Format == "json"
}

func (c OutputConfig) ApplySort(records []model.Record) ([]model.Record, error) {
	if err := c.ValidateSort(); err != nil {
		return nil, err
	}
	if c.Sort == "" {
		return records, nil
	}

	desc := c.Desc || (!c.Asc && !c.Desc)
	sorted := append([]model.Record(nil), records...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return compareRecordField(sorted[i], sorted[j], c.Sort, desc)
	})
	return sorted, nil
}

func (c OutputConfig) ValidateSort() error {
	if c.Sort == "" {
		if c.Asc || c.Desc {
			return NewError(ErrInvalidArg, "参数错误", "", "sort", fmt.Errorf("--sort is required when using --asc or --desc"), nil)
		}
		return nil
	}
	if c.Asc && c.Desc {
		return NewError(ErrInvalidArg, "参数错误", c.Sort, "sort", fmt.Errorf("cannot use --asc and --desc together"), nil)
	}
	return nil
}

func compareRecordField(left, right model.Record, field string, desc bool) bool {
	lv, lok := left[field]
	rv, rok := right[field]

	if !lok || isNilValue(lv) {
		return false
	}
	if !rok || isNilValue(rv) {
		return true
	}

	lf, lokNum := toFloat64(lv)
	rf, rokNum := toFloat64(rv)
	if lokNum && rokNum {
		if lf == rf {
			return false
		}
		if desc {
			return lf > rf
		}
		return lf < rf
	}

	ls := fmt.Sprintf("%v", lv)
	rs := fmt.Sprintf("%v", rv)
	if ls == rs {
		return false
	}
	if desc {
		return ls > rs
	}
	return ls < rs
}

func isNilValue(v interface{}) bool {
	return v == nil
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(n), 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}
