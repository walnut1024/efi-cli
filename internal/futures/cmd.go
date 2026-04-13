package futures

import (
	"github.com/spf13/cobra"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var (
	futuresOutput cliutil.OutputConfig
	futuresBeg    string
	futuresEnd    string
	futuresKlt    int
	futuresFqt    int
	futuresMax    int
)

func formatOrExit(records []model.Record) { futuresOutput.FormatOrExit(records) }

func AddCommands(root *cobra.Command) {
	futuresCmd := &cobra.Command{Use: "futures", Short: "Futures data commands"}
	root.AddCommand(futuresCmd)

	realtimeCmd := &cobra.Command{
		Use:   "realtime",
		Short: "Get realtime futures quotes",
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetRealtime()
			if err != nil {
				cliutil.PrintErrorAndExit(err, futuresOutput)
			}
			formatOrExit(records)
		},
	}
	futuresOutput.AddFlags(realtimeCmd)
	futuresCmd.AddCommand(realtimeCmd)

	historyCmd := &cobra.Command{
		Use:   "history [secid]",
		Short: "Get futures K-line history (secid e.g. 113.ag1024)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetHistory(args[0], futuresBeg, futuresEnd, futuresKlt, futuresFqt)
			if err != nil {
				cliutil.PrintErrorAndExit(err, futuresOutput)
			}
			formatOrExit(records)
		},
	}
	historyCmd.Flags().IntVar(&futuresKlt, "klt", 101, "Period: 1/5/15/30/60min, 101=day, 102=week, 103=month")
	historyCmd.Flags().IntVar(&futuresFqt, "fqt", 1, "Adjust: 0=none, 1=forward, 2=backward")
	historyCmd.Flags().StringVar(&futuresBeg, "beg", "19000101", "Start date YYYYMMDD")
	historyCmd.Flags().StringVar(&futuresEnd, "end", "20500101", "End date YYYYMMDD")
	futuresOutput.AddFlags(historyCmd)
	futuresCmd.AddCommand(historyCmd)

	dealCmd := &cobra.Command{
		Use:   "deal [secid]",
		Short: "Get futures deal details (secid e.g. 115.ZCM)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetDealDetail(args[0], futuresMax)
			if err != nil {
				cliutil.PrintErrorAndExit(err, futuresOutput)
			}
			formatOrExit(records)
		},
	}
	dealCmd.Flags().IntVar(&futuresMax, "max", 10000, "Max records count")
	futuresOutput.AddFlags(dealCmd)
	futuresCmd.AddCommand(dealCmd)
}
