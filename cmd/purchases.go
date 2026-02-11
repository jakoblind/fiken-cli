package cmd

import (
	"fmt"
	"net/url"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var purchasesCmd = &cobra.Command{
	Use:   "purchases",
	Short: "Manage purchases",
	Long:  "List and manage purchases/expenses.",
}

var purchasesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List purchases",
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
		params.Set("pageSize", "25")

		endpoint := fmt.Sprintf(api.EndpointPurchases, slug)

		var purchases []api.Purchase
		page := 0
		for {
			params.Set("page", fmt.Sprintf("%d", page))
			var pagePurchases []api.Purchase
			pagination, err := client.GetWithParams(endpoint, params, &pagePurchases)
			if err != nil {
				return fmt.Errorf("fetching purchases: %w", err)
			}
			purchases = append(purchases, pagePurchases...)

			if pagination == nil || page+1 >= pagination.PageCount || len(pagePurchases) == 0 {
				break
			}
			page++
			// Only fetch first few pages by default
			if page >= 4 {
				break
			}
		}

		if jsonOutput {
			return output.PrintJSON(purchases)
		}

		if len(purchases) == 0 {
			output.PrintInfo("No purchases found.")
			return nil
		}

		table := output.NewTable("ID", "DATE", "KIND", "PAID", "AMOUNT", "IDENTIFIER")
		for _, p := range purchases {
			paid := "No"
			if p.Paid {
				paid = "Yes"
			}
			// Sum net amounts from lines
			var totalNet int64
			for _, l := range p.Lines {
				totalNet += l.NetAmount
			}
			table.AddRow(
				fmt.Sprintf("%d", p.PurchaseId),
				p.Date,
				p.Kind,
				paid,
				output.FormatAmount(totalNet),
				p.Identifier,
			)
		}
		table.Print()

		fmt.Printf("\n%d purchases\n", len(purchases))
		return nil
	},
}

var purchasesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a purchase",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not yet implemented. Use the Fiken web UI to create purchases for now.")
		return nil
	},
}

func init() {
	purchasesCmd.AddCommand(purchasesListCmd)
	purchasesCmd.AddCommand(purchasesCreateCmd)
	rootCmd.AddCommand(purchasesCmd)
}
