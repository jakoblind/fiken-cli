package cmd

import (
	"fmt"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var bankCmd = &cobra.Command{
	Use:   "bank",
	Short: "Manage bank accounts",
}

var bankListCmd = &cobra.Command{
	Use:   "list",
	Short: "List bank accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointBankAccounts, slug)

		var bankAccounts []api.BankAccount
		_, err = client.Get(endpoint, &bankAccounts)
		if err != nil {
			return fmt.Errorf("fetching bank accounts: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(bankAccounts)
		}

		if len(bankAccounts) == 0 {
			output.PrintInfo("No bank accounts found.")
			return nil
		}

		table := output.NewTable("ID", "NAME", "ACCOUNT", "BANK ACCOUNT", "TYPE", "ACTIVE")
		for _, ba := range bankAccounts {
			active := "Yes"
			if ba.Inactive {
				active = "No"
			}
			table.AddRow(
				fmt.Sprintf("%d", ba.BankAccountId),
				ba.Name,
				ba.AccountCode,
				ba.BankAccountNumber,
				ba.Type,
				active,
			)
		}
		table.Print()

		return nil
	},
}

func init() {
	bankCmd.AddCommand(bankListCmd)
	rootCmd.AddCommand(bankCmd)
}
