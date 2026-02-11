package cmd

import (
	"fmt"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var balancesCmd = &cobra.Command{
	Use:   "balances",
	Short: "List account balances",
	Long:  "List account balances for the selected company.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointAccountBalances, slug)

		var balances []api.AccountBalance
		_, err = client.Get(endpoint, &balances)
		if err != nil {
			return fmt.Errorf("fetching balances: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(balances)
		}

		if len(balances) == 0 {
			output.PrintInfo("No account balances found.")
			return nil
		}

		table := output.NewTable("CODE", "NAME", "BALANCE")
		for _, b := range balances {
			table.AddRow(b.Account.Code, b.Account.Name, output.FormatAmount(b.Balance))
		}
		table.Print()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(balancesCmd)
}
