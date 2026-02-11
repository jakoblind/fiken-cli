package api

const (
	BaseURL = "https://api.fiken.no/api/v2"

	// Company endpoints
	EndpointCompanies = "/companies"

	// Endpoints under /companies/{slug}
	EndpointAccounts        = "/companies/%s/accounts"
	EndpointAccountBalances = "/companies/%s/accountBalances"
	EndpointBankAccounts    = "/companies/%s/bankAccounts"
	EndpointInbox           = "/companies/%s/inbox"
	EndpointPurchases       = "/companies/%s/purchases"
	EndpointSales           = "/companies/%s/sales"
	EndpointInvoices        = "/companies/%s/invoices"
	EndpointJournalEntries  = "/companies/%s/journalEntries"
	EndpointTransactions    = "/companies/%s/transactions"
	EndpointContacts        = "/companies/%s/contacts"
)

// Pagination defaults
const (
	DefaultPageSize = 25
	MaxPageSize     = 100
)

// Response headers
const (
	HeaderPage        = "Fiken-Api-Page"
	HeaderPageSize    = "Fiken-Api-Page-Size"
	HeaderPageCount   = "Fiken-Api-Page-Count"
	HeaderResultCount = "Fiken-Api-Result-Count"
)
