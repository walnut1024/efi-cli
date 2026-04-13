package index

import (
	"github.com/spf13/cobra"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var (
	indexOutput cliutil.OutputConfig
	indexBeg    string
	indexEnd    string
	indexKlt    int
)

func formatOrExit(records []model.Record) { indexOutput.FormatOrExit(records) }

func AddCommands(root *cobra.Command) {
	indexCmd := &cobra.Command{
		Use:   "index",
		Short: "Index data commands (上证指数, 沪深300, etc.)",
	}
	root.AddCommand(indexCmd)

	// efi index quote <names/codes...>
	quoteCmd := &cobra.Command{
		Use:   "quote [names or codes]",
		Short: "Get realtime index quote",
		Long:  "Supports: 上证指数, 深证成指, 创业板指, 沪深300, 中证500, 中证1000, 上证50, 科创50, 恒生指数\nOr secid format: 1.000001, 0.399001, etc.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetQuote(args)
			if err != nil {
				cliutil.PrintErrorAndExit(err, indexOutput)
			}
			formatOrExit(records)
		},
	}
	indexOutput.AddFlags(quoteCmd)
	indexCmd.AddCommand(quoteCmd)

	// efi index history <name or code>
	historyCmd := &cobra.Command{
		Use:   "history [name or code]",
		Short: "Get index K-line history",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetHistory(args[0], indexBeg, indexEnd, indexKlt)
			if err != nil {
				cliutil.PrintErrorAndExit(err, indexOutput)
			}
			formatOrExit(records)
		},
	}
	historyCmd.Flags().StringVar(&indexBeg, "beg", "19000101", "Start date YYYYMMDD")
	historyCmd.Flags().StringVar(&indexEnd, "end", "20500101", "End date YYYYMMDD")
	historyCmd.Flags().IntVar(&indexKlt, "klt", 101, "Period: 1/5/15/30/60min, 101=day, 102=week, 103=month")
	indexOutput.AddFlags(historyCmd)
	indexCmd.AddCommand(historyCmd)
}
