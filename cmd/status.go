package cmd

import (
	"fmt"
	"net/url"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Dashboard overview",
	Long:  "Show a dashboard overview with pending items and key metrics.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		if jsonOutput {
			return statusJSON(client, slug)
		}

		fmt.Printf("📊 Dashboard for: %s\n", slug)
		fmt.Println(repeatStr("─", 50))

		// Inbox
		fmt.Print("\n📥 Inbox: ")
		inboxParams := url.Values{"pageSize": {"1"}}
		var inboxDocs []api.InboxDocument
		pagination, err := client.GetWithParams(fmt.Sprintf(api.EndpointInbox, slug), inboxParams, &inboxDocs)
		if err != nil {
			fmt.Printf("error (%v)\n", err)
		} else if pagination != nil {
			fmt.Printf("%d documents\n", pagination.ResultCount)
		} else {
			fmt.Printf("%d documents\n", len(inboxDocs))
		}

		// Unpaid purchases
		fmt.Print("🛒 Purchases: ")
		purchaseParams := url.Values{"pageSize": {"1"}}
		var purchases []api.Purchase
		pagination, err = client.GetWithParams(fmt.Sprintf(api.EndpointPurchases, slug), purchaseParams, &purchases)
		if err != nil {
			fmt.Printf("error (%v)\n", err)
		} else if pagination != nil {
			fmt.Printf("%d total\n", pagination.ResultCount)
		} else {
			fmt.Printf("%d total\n", len(purchases))
		}

		// Bank accounts
		fmt.Print("🏦 Bank accounts: ")
		var bankAccounts []api.BankAccount
		_, err = client.Get(fmt.Sprintf(api.EndpointBankAccounts, slug), &bankAccounts)
		if err != nil {
			fmt.Printf("error (%v)\n", err)
		} else {
			fmt.Printf("%d accounts\n", len(bankAccounts))
			for _, ba := range bankAccounts {
				if !ba.Inactive {
					fmt.Printf("   %s (%s) - %s\n", ba.Name, ba.AccountCode, ba.BankAccountNumber)
				}
			}
		}

		// Contacts
		fmt.Print("👥 Contacts: ")
		contactParams := url.Values{"pageSize": {"1"}}
		var contacts []api.Contact
		pagination, err = client.GetWithParams(fmt.Sprintf(api.EndpointContacts, slug), contactParams, &contacts)
		if err != nil {
			fmt.Printf("error (%v)\n", err)
		} else if pagination != nil {
			fmt.Printf("%d total\n", pagination.ResultCount)
		} else {
			fmt.Printf("%d total\n", len(contacts))
		}

		fmt.Println()
		return nil
	},
}

type statusData struct {
	Company        string `json:"company"`
	InboxCount     int    `json:"inbox_count"`
	PurchaseCount  int    `json:"purchase_count"`
	BankAccounts   int    `json:"bank_accounts"`
	ContactCount   int    `json:"contact_count"`
}

func statusJSON(client *api.Client, slug string) error {
	data := statusData{Company: slug}

	inboxParams := url.Values{"pageSize": {"1"}}
	var inboxDocs []api.InboxDocument
	pagination, err := client.GetWithParams(fmt.Sprintf(api.EndpointInbox, slug), inboxParams, &inboxDocs)
	if err == nil && pagination != nil {
		data.InboxCount = pagination.ResultCount
	}

	purchaseParams := url.Values{"pageSize": {"1"}}
	var purchases []api.Purchase
	pagination, err = client.GetWithParams(fmt.Sprintf(api.EndpointPurchases, slug), purchaseParams, &purchases)
	if err == nil && pagination != nil {
		data.PurchaseCount = pagination.ResultCount
	}

	var bankAccounts []api.BankAccount
	_, err = client.Get(fmt.Sprintf(api.EndpointBankAccounts, slug), &bankAccounts)
	if err == nil {
		data.BankAccounts = len(bankAccounts)
	}

	contactParams := url.Values{"pageSize": {"1"}}
	var contacts []api.Contact
	pagination, err = client.GetWithParams(fmt.Sprintf(api.EndpointContacts, slug), contactParams, &contacts)
	if err == nil && pagination != nil {
		data.ContactCount = pagination.ResultCount
	}

	return output.PrintJSON(data)
}

func repeatStr(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
