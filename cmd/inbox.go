package cmd

import (
	"fmt"
	"net/url"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var inboxStatus string

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "List EHF inbox documents",
	Long:  "List documents in the EHF (electronic invoice) inbox.",
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
		if inboxStatus != "" {
			params.Set("status", inboxStatus)
		}
		params.Set("pageSize", "25")

		endpoint := fmt.Sprintf(api.EndpointInbox, slug)

		var documents []api.InboxDocument
		page := 0
		for {
			params.Set("page", fmt.Sprintf("%d", page))
			var pageDocs []api.InboxDocument
			pagination, err := client.GetWithParams(endpoint, params, &pageDocs)
			if err != nil {
				return fmt.Errorf("fetching inbox: %w", err)
			}
			documents = append(documents, pageDocs...)

			if pagination == nil || page+1 >= pagination.PageCount {
				break
			}
			page++
		}

		if jsonOutput {
			return output.PrintJSON(documents)
		}

		if len(documents) == 0 {
			output.PrintInfo("Inbox is empty.")
			return nil
		}

		table := output.NewTable("ID", "NAME", "FILENAME", "STATUS", "DATE")
		for _, d := range documents {
			table.AddRow(
				fmt.Sprintf("%d", d.DocumentId),
				d.Name,
				d.Filename,
				d.Status,
				d.CreatedDate.Format("2006-01-02"),
			)
		}
		table.Print()

		fmt.Printf("\n%d documents\n", len(documents))
		return nil
	},
}

func init() {
	inboxCmd.Flags().StringVar(&inboxStatus, "status", "", "Filter by status (pending, processed)")
	rootCmd.AddCommand(inboxCmd)
}
