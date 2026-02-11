package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var salesCmd = &cobra.Command{
	Use:   "sales",
	Short: "Manage sales",
}

var salesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sales",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not yet implemented.")
		return nil
	},
}

func init() {
	salesCmd.AddCommand(salesListCmd)
	rootCmd.AddCommand(salesCmd)
}
