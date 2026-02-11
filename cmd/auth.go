package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jakoblind/fiken-cli/auth"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  "Manage your Fiken API token for authentication.",
}

var authTokenCmd = &cobra.Command{
	Use:   "token [token]",
	Short: "Set or show the API token",
	Long: `Set your Fiken Personal API Token for authentication.

Get your token from: https://fiken.no/innstillinger/api
If called without arguments, shows whether a token is configured.

To avoid exposing the token in process listings, pipe it via stdin:
  echo "mytoken" | fiken auth token --stdin
  fiken auth token --stdin < tokenfile`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fromStdin, _ := cmd.Flags().GetBool("stdin")

		if len(args) == 0 && !fromStdin {
			// Show token status
			if auth.TokenExists() {
				token, err := auth.LoadToken()
				if err != nil {
					return err
				}
				token = strings.TrimSpace(token)
				output.PrintSuccess(fmt.Sprintf("Token configured: %s", maskToken(token)))
			} else {
				output.PrintInfo("No token configured. Run 'fiken auth token <token>' to authenticate.")
			}
			return nil
		}

		var token string
		if fromStdin {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				token = strings.TrimSpace(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("reading token from stdin: %w", err)
			}
		} else {
			token = strings.TrimSpace(args[0])
		}

		if token == "" {
			return fmt.Errorf("token cannot be empty")
		}

		if err := auth.SaveToken(token); err != nil {
			return fmt.Errorf("saving token: %w", err)
		}

		output.PrintSuccess("Token saved to ~/.config/fiken/token")
		return nil
	},
}

// maskToken safely masks a token for display, handling short tokens.
func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored token",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.RemoveToken(); err != nil {
			return err
		}
		output.PrintSuccess("Token removed")
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !auth.TokenExists() {
			output.PrintError("Not authenticated. Run 'fiken auth token <token>' to set up.")
			return nil
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		var companies []interface{}
		_, err = client.Get("/companies", &companies)
		if err != nil {
			output.PrintError(fmt.Sprintf("Token is invalid or expired: %v", err))
			return nil
		}

		output.PrintSuccess(fmt.Sprintf("Authenticated. Access to %d company(ies).", len(companies)))
		return nil
	},
}

func init() {
	authTokenCmd.Flags().Bool("stdin", false, "Read token from stdin instead of command arguments")
	authCmd.AddCommand(authTokenCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
	rootCmd.AddCommand(authCmd)
}
