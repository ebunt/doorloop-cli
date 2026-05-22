package cli

import (
	"fmt"

	"doorloop-pp-cli/internal/store"

	"github.com/spf13/cobra"
)

type ledgerEntry struct {
	Date           string  `json:"date"`
	Type           string  `json:"type"`
	Description    string  `json:"description"`
	Amount         float64 `json:"amount"`
	RunningBalance float64 `json:"running_balance"`
}

func newLedgerCmd(flags *rootFlags) *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:   "ledger <lease-id>",
		Short: "Bank-statement-style ledger for a lease: charges, payments, credits",
		Long: `View a chronological ledger for any lease — every charge, payment, credit, and
returned payment with a running balance. Queries local SQLite — run 'sync' first.`,
		Example: `  doorloop-pp-cli ledger lease_456
  doorloop-pp-cli ledger lease_456 --json
  doorloop-pp-cli ledger lease_456 --json --agent`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			leaseID := args[0]
			if dbPath == "" {
				dbPath = defaultDBPath("doorloop-pp-cli")
			}
			db, err := store.OpenReadOnly(dbPath)
			if err != nil {
				return fmt.Errorf("opening local database (run 'sync' first): %w", err)
			}
			defer db.Close()

			// UNION all four financial tables for this lease, sorted chronologically
			q := `
				SELECT 'charge' AS type,
					COALESCE(json_extract(data, '$.date'), '') AS date,
					CAST(COALESCE(json_extract(data, '$.amount'), '0') AS REAL) AS amount,
					COALESCE(json_extract(data, '$.description'), '') AS description
				FROM resources WHERE resource_type = 'lease_charges'
				  AND (json_extract(data, '$.leaseId') = ? OR json_extract(data, '$.lease_id') = ?)
				UNION ALL
				SELECT 'payment',
					COALESCE(json_extract(data, '$.date'), ''),
					-CAST(COALESCE(json_extract(data, '$.amount'), '0') AS REAL),
					''
				FROM resources WHERE resource_type = 'lease_payments'
				  AND (json_extract(data, '$.leaseId') = ? OR json_extract(data, '$.lease_id') = ?)
				UNION ALL
				SELECT 'credit',
					COALESCE(json_extract(data, '$.date'), ''),
					-CAST(COALESCE(json_extract(data, '$.amount'), '0') AS REAL),
					COALESCE(json_extract(data, '$.description'), '')
				FROM resources WHERE resource_type = 'lease_credits'
				  AND (json_extract(data, '$.leaseId') = ? OR json_extract(data, '$.lease_id') = ?)
				UNION ALL
				SELECT 'returned_payment',
					COALESCE(json_extract(data, '$.date'), ''),
					CAST(COALESCE(json_extract(data, '$.amount'), '0') AS REAL),
					''
				FROM resources WHERE resource_type = 'lease_returned_payments'
				  AND (json_extract(data, '$.leaseId') = ? OR json_extract(data, '$.lease_id') = ?)
				ORDER BY date ASC, type ASC`

			rows, err := db.Query(q,
				leaseID, leaseID, // charges
				leaseID, leaseID, // payments
				leaseID, leaseID, // credits
				leaseID, leaseID, // returned_payments
			)
			if err != nil {
				return fmt.Errorf("querying ledger: %w", err)
			}
			defer rows.Close()

			var entries []ledgerEntry
			var balance float64
			for rows.Next() {
				var e ledgerEntry
				if err := rows.Scan(&e.Type, &e.Date, &e.Amount, &e.Description); err != nil {
					continue
				}
				balance += e.Amount
				e.RunningBalance = balance
				entries = append(entries, e)
			}
			if err := rows.Err(); err != nil {
				return fmt.Errorf("reading ledger: %w", err)
			}

			if len(entries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "[]")
				return nil
			}
			return printJSONFiltered(cmd.OutOrStdout(), entries, flags)
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: ~/.local/share/doorloop-pp-cli/data.db)")
	return cmd
}
