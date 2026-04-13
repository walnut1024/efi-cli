package bond

import (
	"github.com/spf13/cobra"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var (
	bondOutput cliutil.OutputConfig
	bondBeg    string
	bondEnd    string
	bondKlt    int
	bondFqt    int
)

func formatOrExit(records []model.Record) { bondOutput.FormatOrExit(records) }

func AddCommands(root *cobra.Command) {
	bondCmd := &cobra.Command{Use: "bond", Short: "Bond data commands"}
	root.AddCommand(bondCmd)

	infoCmd := &cobra.Command{
		Use:   "info [code]",
		Short: "Get bond base info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			record, err := GetBaseInfo(args[0])
			if err != nil {
				cliutil.PrintErrorAndExit(err, bondOutput)
			}
			formatOrExit([]model.Record{record})
		},
	}
	bondOutput.AddFlags(infoCmd)
	bondCmd.AddCommand(infoCmd)

	historyCmd := &cobra.Command{
		Use:   "history [secid]",
		Short: "Get bond K-line history (secid e.g. 113.ag1024)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetHistory(args[0], bondBeg, bondEnd, bondKlt, bondFqt)
			if err != nil {
				cliutil.PrintErrorAndExit(err, bondOutput)
			}
			formatOrExit(records)
		},
	}
	historyCmd.Flags().IntVar(&bondKlt, "klt", 101, "Period: 1/5/15/30/60min, 101=day, 102=week, 103=month")
	historyCmd.Flags().IntVar(&bondFqt, "fqt", 1, "Adjust: 0=none, 1=forward, 2=backward")
	historyCmd.Flags().StringVar(&bondBeg, "beg", "19000101", "Start date YYYYMMDD")
	historyCmd.Flags().StringVar(&bondEnd, "end", "20500101", "End date YYYYMMDD")
	bondOutput.AddFlags(historyCmd)
	bondCmd.AddCommand(historyCmd)

	realtimeCmd := &cobra.Command{
		Use:   "realtime",
		Short: "Get realtime bond quotes",
		Run: func(cmd *cobra.Command, args []string) {
			records, err := GetRealtime()
			if err != nil {
				cliutil.PrintErrorAndExit(err, bondOutput)
			}
			formatOrExit(records)
		},
	}
	bondOutput.AddFlags(realtimeCmd)
	bondCmd.AddCommand(realtimeCmd)
}
