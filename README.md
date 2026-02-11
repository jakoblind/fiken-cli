# fiken-cli

A command-line client for the [Fiken.no](https://fiken.no) accounting API. Manage your Norwegian business accounting from the terminal.

## Features

- 🏢 List and manage companies
- 📊 Chart of accounts and balances
- 🛒 View purchases and expenses
- 📥 EHF inbox management
- 🏦 Bank account overview
- 📋 Dashboard with key metrics
- 🔄 JSON output for scripting
- ⚡ Built-in rate limiting and pagination

## Installation

### From source

```bash
git clone https://github.com/jakoblind/fiken-cli.git
cd fiken-cli
make install
```

### Build locally

```bash
make build
./fiken --help
```

## Quick Start

### 1. Get your API token

Go to [Fiken API Settings](https://fiken.no/innstillinger/api) and create a Personal API Token.

### 2. Authenticate

```bash
fiken auth token <your-token>
```

### 3. List your companies

```bash
fiken companies
```

### 4. Set a default company (optional)

```bash
fiken companies default <company-slug>
```

### 5. View your dashboard

```bash
fiken status
```

## Command Reference

### Authentication

```bash
fiken auth token <token>    # Save API token
fiken auth token            # Show token status
fiken auth status           # Verify token works
fiken auth logout           # Remove stored token
```

### Companies

```bash
fiken companies             # List all companies
fiken companies default     # Show default company
fiken companies default <slug>  # Set default company
```

### Accounts

```bash
fiken accounts              # List chart of accounts
fiken accounts --from 1000 --to 2000  # Filter by account range
```

### Balances

```bash
fiken balances              # List account balances
```

### Bank Accounts

```bash
fiken bank list             # List bank accounts
```

### Inbox (EHF)

```bash
fiken inbox                 # List all inbox documents
fiken inbox --status pending    # Filter by status
```

### Purchases

```bash
fiken purchases list        # List purchases
```

### Status Dashboard

```bash
fiken status                # Overview of pending items
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON (default: table) |
| `--no-input` | Non-interactive mode |
| `--company <slug>` | Select company (auto-detected if only one) |
| `--keyring-backend <backend>` | Keyring backend (default: `auto`) |

## Credential Storage

Credentials (API token, default company) are stored securely using your OS keyring via [99designs/keyring](https://github.com/99designs/keyring).

### Supported backends

| Backend | OS | Flag value |
|---------|----|------------|
| Secret Service (GNOME Keyring / KDE Wallet) | Linux | `secret-service` |
| Keychain | macOS | `keychain` |
| Windows Credential Manager | Windows | `wincred` |
| [pass](https://www.passwordstore.org/) | Linux/macOS | `pass` |
| Encrypted file | Any (fallback) | `file` |

By default (`auto`), the best available backend is used. The encrypted file backend is the last-resort fallback and will prompt for a password.

### Choosing a backend

```bash
# Use a specific backend
fiken --keyring-backend file auth token <token>

# Or set via environment variable
export FIKEN_KEYRING_BACKEND=file
fiken auth token <token>
```

### Migration from plaintext storage

If you previously stored your token in `~/.config/fiken/token`, it will be automatically migrated to the keyring on first use. The plaintext file is deleted after successful migration.

## API Details

- Base URL: `https://api.fiken.no/api/v2`
- Auth: Bearer token (Personal API Token)
- Amounts are in cents (øre): `100000` = `1 000,00 kr`
- Rate limit: max 4 requests/second (enforced by client)
- Pagination: automatic for large result sets

## Examples

### List companies as JSON

```bash
fiken companies --json
```

### Script: get all account codes

```bash
fiken accounts --json | jq '.[].code'
```

### Use with a specific company

```bash
fiken purchases list --company my-company-slug
```

### Quick status check

```bash
fiken status --json | jq '.inbox_count'
```

## Development

```bash
make build    # Build binary
make test     # Run tests
make fmt      # Format code
make lint     # Run linter
make clean    # Clean build artifacts
```

## License

MIT – see [LICENSE](LICENSE).
