package cmd

import (
	"fmt"
	"net/url"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var (
	accountsFromCode string
	accountsToCode   string
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List chart of accounts",
	Long:  "List the chart of accounts for the selected company.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		params := url.Values{}
		if accountsFromCode != "" {
			params.Set("fromAccount", accountsFromCode)
		}
		if accountsToCode != "" {
			params.Set("toAccount", accountsToCode)
		}
		params.Set("pageSize", "100")

		endpoint := fmt.Sprintf(api.EndpointAccounts, slug)

		var accounts []api.Account
		page := 0
		for {
			params.Set("page", fmt.Sprintf("%d", page))
			var pageAccounts []api.Account
			pagination, err := client.GetWithParams(endpoint, params, &pageAccounts)
			if err != nil {
				return fmt.Errorf("fetching accounts: %w", err)
			}
			accounts = append(accounts, pageAccounts...)

			if pagination == nil || page+1 >= pagination.PageCount {
				break
			}
			page++
		}

		if jsonOutput {
			return output.PrintJSON(accounts)
		}

		if len(accounts) == 0 {
			output.PrintInfo("No accounts found.")
			return nil
		}

		table := output.NewTable("CODE", "NAME")
		for _, a := range accounts {
			table.AddRow(a.Code, a.Name)
		}
		table.Print()

		fmt.Printf("\n%d accounts\n", len(accounts))
		return nil
	},
}

func init() {
	accountsCmd.Flags().StringVar(&accountsFromCode, "from", "", "Filter from account code")
	accountsCmd.Flags().StringVar(&accountsToCode, "to", "", "Filter to account code")
	rootCmd.AddCommand(accountsCmd)
}
