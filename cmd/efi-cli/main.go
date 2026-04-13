package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/walnut1024/efi-cli/internal/bond"
	"github.com/walnut1024/efi-cli/internal/fund"
	"github.com/walnut1024/efi-cli/internal/futures"
	"github.com/walnut1024/efi-cli/internal/index"
	"github.com/walnut1024/efi-cli/internal/searchcmd"
	"github.com/walnut1024/efi-cli/internal/stock"
)

var rootCmd = &cobra.Command{
	Use:   "efi-cli",
	Short: "Financial data CLI for AI agents",
	Long:  "efi-cli - query Chinese/US/HK stock, fund, bond, futures data from public market APIs",
}

func main() {
	searchcmd.AddCommands(rootCmd)
	stock.AddCommands(rootCmd)
	fund.AddCommands(rootCmd)
	bond.AddCommands(rootCmd)
	futures.AddCommands(rootCmd)
	index.AddCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
