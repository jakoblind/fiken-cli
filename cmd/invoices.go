package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var invoicesCmd = &cobra.Command{
	Use:   "invoices",
	Short: "Manage invoices",
}

var invoicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List invoices",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not yet implemented.")
		return nil
	},
}

func init() {
	invoicesCmd.AddCommand(invoicesListCmd)
	rootCmd.AddCommand(invoicesCmd)
}
