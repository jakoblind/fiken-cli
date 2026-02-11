package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// PrintJSON outputs data as formatted JSON.
func PrintJSON(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// Table helps build and print tabular output.
type Table struct {
	headers []string
	rows    [][]string
}

// NewTable creates a new table with the given headers.
func NewTable(headers ...string) *Table {
	return &Table{
		headers: headers,
	}
}

// AddRow adds a row to the table.
func (t *Table) AddRow(values ...string) {
	t.rows = append(t.rows, values)
}

// Print outputs the table to stdout.
func (t *Table) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// Print headers
	fmt.Fprintln(w, strings.Join(t.headers, "\t"))
	// Print separator
	sep := make([]string, len(t.headers))
	for i, h := range t.headers {
		sep[i] = strings.Repeat("─", len(h))
	}
	fmt.Fprintln(w, strings.Join(sep, "\t"))
	// Print rows
	for _, row := range t.rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}

// FormatAmount converts cents to a human-readable amount string (e.g., 100000 -> "1 000,00").
func FormatAmount(cents int64) string {
	negative := cents < 0
	if negative {
		cents = -cents
	}

	kr := cents / 100
	ore := cents % 100

	// Format with space as thousands separator
	krStr := formatThousands(kr)

	result := fmt.Sprintf("%s,%02d", krStr, ore)
	if negative {
		result = "-" + result
	}
	return result
}

func formatThousands(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	s := fmt.Sprintf("%d", n)
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	return strings.Join(parts, " ")
}

// PrintSuccess prints a success message.
func PrintSuccess(msg string) {
	fmt.Printf("✓ %s\n", msg)
}

// PrintError prints an error message.
func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "✗ %s\n", msg)
}

// PrintInfo prints an informational message.
func PrintInfo(msg string) {
	fmt.Printf("ℹ %s\n", msg)
}
