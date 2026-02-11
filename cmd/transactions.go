package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Manage transactions",
}

var transactionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List transactions",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not yet implemented.")
		return nil
	},
}

func init() {
	transactionsCmd.AddCommand(transactionsListCmd)
	rootCmd.AddCommand(transactionsCmd)
}
