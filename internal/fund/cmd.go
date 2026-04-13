package fund

import (
	"github.com/spf13/cobra"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var (
	fundOutput        cliutil.OutputConfig
	fundDate          string
	fundType          string
	fundSort          string
	fundScreenSizeMin float64
	fundScreenSizeMax float64
	fundScreenY1Min   float64
	fundScreenY3Min   float64
	fundScreenY5Min   float64
)

func AddCommands(root *cobra.Command) {
	fundCmd := &cobra.Command{
		Use:   "fund",
		Short: "Fund data commands",
	}
	root.AddCommand(fundCmd)

	// efi fund history <code>
	historyCmd := &cobra.Command{
		Use:   "history [code]",
		Short: "Get fund NAV history",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetHistory(args[0], 40000)
			if err != nil {
				cliutil.PrintErrorAndExit(err, fundOutput)
			}
			fundOutput.FormatOrExit(records)
		},
	}
	fundOutput.AddFlags(historyCmd)
	fundCmd.AddCommand(historyCmd)

	// efi fund quote <codes...>
	quoteCmd := &cobra.Command{
		Use:   "quote [codes]",
		Short: "Get realtime fund estimate",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetRealtime(args)
			if err != nil {
				cliutil.PrintErrorAndExit(err, fundOutput)
			}
			fundOutput.FormatOrExit(records)
		},
	}
	fundOutput.AddFlags(quoteCmd)
	fundCmd.AddCommand(quoteCmd)

	// efi fund info <code>
	infoCmd := &cobra.Command{
		Use:   "info [code]",
		Short: "Get fund basic info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			record, err := GetBaseInfo(args[0])
			if err != nil {
				cliutil.PrintErrorAndExit(err, fundOutput)
			}
			fundOutput.FormatOrExit([]model.Record{record})
		},
	}
	fundOutput.AddFlags(infoCmd)
	fundCmd.AddCommand(infoCmd)

	// efi fund position <code> [--date <date>]
	positionCmd := &cobra.Command{
		Use:   "position [code]",
		Short: "Get fund stock holdings",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetInvestPosition(args[0], fundDate)
			if err != nil {
				cliutil.PrintErrorAndExit(err, fundOutput)
			}
			fundOutput.FormatOrExit(records)
		},
	}
	positionCmd.Flags().StringVar(&fundDate, "date", "", "Report date YYYY-MM-DD")
	fundOutput.AddFlags(positionCmd)
	fundCmd.AddCommand(positionCmd)

	// efi fund period <code>
	periodCmd := &cobra.Command{
		Use:   "period [code]",
		Short: "Get fund period returns",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetPeriodChange(args[0])
			if err != nil {
				cliutil.PrintErrorAndExit(err, fundOutput)
			}
			fundOutput.FormatOrExit(records)
		},
	}
	fundOutput.AddFlags(periodCmd)
	fundCmd.AddCommand(periodCmd)

	// efi fund asset <code> [--date <date>]
	assetCmd := &cobra.Command{
		Use:   "asset [code]",
		Short: "Get fund asset allocation",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetAssetAllocation(args[0], fundDate)
			if err != nil {
				cliutil.PrintErrorAndExit(err, fundOutput)
			}
			fundOutput.FormatOrExit(records)
		},
	}
	assetCmd.Flags().StringVar(&fundDate, "date", "", "Report date YYYY-MM-DD")
	fundOutput.AddFlags(assetCmd)
	fundCmd.AddCommand(assetCmd)

	// efi fund rank [--type gp|hh|zq|all] [--sort 3y|2y|1y|6m|3m|1m|1w]
	rankCmd := &cobra.Command{
		Use:   "rank",
		Short: "Get fund ranking by performance",
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetFundRank(fundType, fundSort, fundOutput.Limit)
			if err != nil {
				cliutil.PrintErrorAndExit(err, fundOutput)
			}
			fundOutput.FormatOrExit(records)
		},
	}
	rankCmd.Flags().StringVar(&fundType, "type", "all", "Fund type: all, gp (stock), hh (hybrid), zq (bond)")
	rankCmd.Flags().StringVar(&fundSort, "sort", "3y", "Sort by: 1w, 1m, 3m, 6m, 1y, 2y, 3y, 5y")
	fundOutput.AddFlags(rankCmd)
	fundCmd.AddCommand(rankCmd)

	// efi fund screen
	screenCmd := &cobra.Command{
		Use:   "screen",
		Short: "Screen funds by type, size, and return thresholds",
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetFundScreen(fundType, fundScreenSizeMin, fundScreenSizeMax, fundScreenY1Min, fundScreenY3Min, fundScreenY5Min, fundSort, fundOutput.Limit)
			if err != nil {
				cliutil.PrintErrorAndExit(err, fundOutput)
			}
			fundOutput.FormatOrExit(records)
		},
	}
	screenCmd.Flags().StringVar(&fundType, "type", "all", "Fund type: all, gp (stock), hh (hybrid), zq (bond)")
	screenCmd.Flags().StringVar(&fundSort, "sort", "3y", "Sort by: 1w, 1m, 3m, 6m, 1y, 2y, 3y, 5y")
	screenCmd.Flags().Float64Var(&fundScreenSizeMin, "size-min", 0, "Min fund size (亿元)")
	screenCmd.Flags().Float64Var(&fundScreenSizeMax, "size-max", 0, "Max fund size (亿元, 0=unlimited)")
	screenCmd.Flags().Float64Var(&fundScreenY1Min, "1y-min", 0, "Min 1-year return (%)")
	screenCmd.Flags().Float64Var(&fundScreenY3Min, "3y-min", 0, "Min 3-year return (%)")
	screenCmd.Flags().Float64Var(&fundScreenY5Min, "5y-min", 0, "Min 5-year return (%)")
	fundOutput.AddFlags(screenCmd)
	fundCmd.AddCommand(screenCmd)
}
