package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Manage contacts",
}

var contactsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List contacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not yet implemented.")
		return nil
	},
}

func init() {
	contactsCmd.AddCommand(contactsListCmd)
	rootCmd.AddCommand(contactsCmd)
}
