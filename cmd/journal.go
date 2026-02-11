package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var journalCmd = &cobra.Command{
	Use:   "journal",
	Short: "Manage journal entries",
}

var journalListCmd = &cobra.Command{
	Use:   "list",
	Short: "List journal entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not yet implemented.")
		return nil
	},
}

func init() {
	journalCmd.AddCommand(journalListCmd)
	rootCmd.AddCommand(journalCmd)
}
