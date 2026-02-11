package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/auth"
	"github.com/spf13/cobra"
)

var (
	jsonOutput     bool
	noInput        bool
	company        string
	keyringBackend string
)

var rootCmd = &cobra.Command{
	Use:   "fiken",
	Short: "Fiken.no accounting API client",
	Long: `fiken is a command-line client for the Fiken.no accounting API.

Manage your Norwegian business accounting from the terminal:
companies, purchases, invoices, bank accounts, and more.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Apply keyring backend setting before any command runs.
		auth.KeyringBackend = keyringBackend
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&noInput, "no-input", false, "Non-interactive mode")
	rootCmd.PersistentFlags().StringVar(&company, "company", "", "Company slug (auto-detected if only one)")
	rootCmd.PersistentFlags().StringVar(&keyringBackend, "keyring-backend", "auto",
		"Keyring backend: auto, secret-service, keychain, wincred, pass, file")

	// Support FIKEN_KEYRING_BACKEND env var as default.
	if env := os.Getenv("FIKEN_KEYRING_BACKEND"); env != "" {
		keyringBackend = env
	}
}

// getClient creates an API client using the stored token.
func getClient() (*api.Client, error) {
	token, err := auth.LoadToken()
	if err != nil {
		return nil, err
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("token is empty. Run 'fiken auth token <token>' to set up authentication")
	}
	return api.NewClient(token), nil
}

// resolveCompany determines which company to use.
// Priority: --company flag > config default > auto-detect (if only one).
func resolveCompany(client *api.Client) (string, error) {
	if company != "" {
		return company, nil
	}

	cfg, err := auth.LoadConfig()
	if err == nil && cfg.DefaultCompany != "" {
		return cfg.DefaultCompany, nil
	}

	// Auto-detect: fetch companies and use the only one if there's just one.
	var companies []api.Company
	_, err = client.Get(api.EndpointCompanies, &companies)
	if err != nil {
		return "", fmt.Errorf("fetching companies: %w", err)
	}

	switch len(companies) {
	case 0:
		return "", fmt.Errorf("no companies found on this account")
	case 1:
		return companies[0].Slug, nil
	default:
		names := make([]string, len(companies))
		for i, c := range companies {
			names[i] = fmt.Sprintf("  %s (%s)", c.Name, c.Slug)
		}
		return "", fmt.Errorf("multiple companies found. Use --company to select one:\n%s", strings.Join(names, "\n"))
	}
}
