package cli

import (
	"fmt"

	"doorloop-pp-cli/internal/store"

	"github.com/spf13/cobra"
)

type cashflowRow struct {
	PropertyID   string  `json:"property_id"`
	PropertyName string  `json:"property_name"`
	Income       float64 `json:"income"`
	Expenses     float64 `json:"expenses"`
	Net          float64 `json:"net"`
	FromDate     string  `json:"from_date,omitempty"`
	ToDate       string  `json:"to_date,omitempty"`
}

func newCashflowCmd(flags *rootFlags) *cobra.Command {
	var property string
	var fromDate string
	var toDate string
	var dbPath string

	cmd := &cobra.Command{
		Use:   "cashflow",
		Short: "Income, expenses, and net cash flow per property",
		Long: `Compute income (lease payments) minus expenses (expenses + vendor bills) per property
for any date range. Queries synced local SQLite — run 'sync' first.`,
		Example: `  doorloop-pp-cli cashflow
  doorloop-pp-cli cashflow --from 2026-05-01 --to 2026-05-31 --json
  doorloop-pp-cli cashflow --property prop_123 --json --agent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbPath == "" {
				dbPath = defaultDBPath("doorloop-pp-cli")
			}
			db, err := store.OpenReadOnly(dbPath)
			if err != nil {
				return fmt.Errorf("opening local database (run 'sync' first): %w", err)
			}
			defer db.Close()

			// Income: sum lease_payments grouped by property via leases join
			incomeQ := `
				SELECT
					COALESCE(json_extract(l.data, '$.propertyId'), json_extract(l.data, '$.property_id'), '') AS property_id,
					COALESCE(SUM(CAST(COALESCE(json_extract(lp.data, '$.amount'), '0') AS REAL)), 0) AS income
				FROM resources lp
				JOIN resources l ON l.resource_type = 'leases'
					AND l.id = COALESCE(json_extract(lp.data, '$.leaseId'), json_extract(lp.data, '$.lease_id'))
				WHERE lp.resource_type = 'lease_payments'`
			incomeArgs := []any{}
			if fromDate != "" {
				incomeQ += ` AND COALESCE(json_extract(lp.data, '$.date'), '') >= ?`
				incomeArgs = append(incomeArgs, fromDate)
			}
			if toDate != "" {
				incomeQ += ` AND COALESCE(json_extract(lp.data, '$.date'), '') <= ?`
				incomeArgs = append(incomeArgs, toDate)
			}
			if property != "" {
				incomeQ += ` AND (json_extract(l.data, '$.propertyId') = ? OR json_extract(l.data, '$.property_id') = ?)`
				incomeArgs = append(incomeArgs, property, property)
			}
			incomeQ += ` GROUP BY property_id`

			incomeByProperty := map[string]float64{}
			rows, err := db.Query(incomeQ, incomeArgs...)
			if err != nil {
				return fmt.Errorf("querying income: %w", err)
			}
			for rows.Next() {
				var pid string
				var income float64
				if err := rows.Scan(&pid, &income); err == nil && pid != "" {
					incomeByProperty[pid] += income
				}
			}
			rows.Close()
			if err := rows.Err(); err != nil {
				return fmt.Errorf("reading income: %w", err)
			}

			// Expenses: sum expenses + vendor_bills grouped by property_id
			expQ := `
				SELECT
					COALESCE(json_extract(data, '$.propertyId'), json_extract(data, '$.property_id'), '') AS property_id,
					COALESCE(SUM(CAST(COALESCE(json_extract(data, '$.amount'), '0') AS REAL)), 0) AS expenses
				FROM resources
				WHERE resource_type IN ('expenses', 'vendor_bills')`
			expArgs := []any{}
			if fromDate != "" {
				expQ += ` AND COALESCE(json_extract(data, '$.date'), '') >= ?`
				expArgs = append(expArgs, fromDate)
			}
			if toDate != "" {
				expQ += ` AND COALESCE(json_extract(data, '$.date'), '') <= ?`
				expArgs = append(expArgs, toDate)
			}
			if property != "" {
				expQ += ` AND (json_extract(data, '$.propertyId') = ? OR json_extract(data, '$.property_id') = ?)`
				expArgs = append(expArgs, property, property)
			}
			expQ += ` GROUP BY property_id`

			expByProperty := map[string]float64{}
			rows, err = db.Query(expQ, expArgs...)
			if err != nil {
				return fmt.Errorf("querying expenses: %w", err)
			}
			for rows.Next() {
				var pid string
				var exp float64
				if err := rows.Scan(&pid, &exp); err == nil && pid != "" {
					expByProperty[pid] += exp
				}
			}
			rows.Close()
			if err := rows.Err(); err != nil {
				return fmt.Errorf("reading expenses: %w", err)
			}

			// Collect all property IDs across both query sets
			seen := map[string]bool{}
			for pid := range incomeByProperty {
				seen[pid] = true
			}
			for pid := range expByProperty {
				seen[pid] = true
			}

			// Look up property names from synced data
			nameByID := map[string]string{}
			nameRows, err := db.Query(`SELECT id, COALESCE(json_extract(data, '$.name'), id) FROM resources WHERE resource_type = 'properties'`)
			if err == nil {
				for nameRows.Next() {
					var pid, name string
					if err := nameRows.Scan(&pid, &name); err == nil {
						nameByID[pid] = name
					}
				}
				nameRows.Close()
			}

			var results []cashflowRow
			for pid := range seen {
				if property != "" && pid != property {
					continue
				}
				income := incomeByProperty[pid]
				exp := expByProperty[pid]
				name := nameByID[pid]
				if name == "" {
					name = pid
				}
				results = append(results, cashflowRow{
					PropertyID:   pid,
					PropertyName: name,
					Income:       income,
					Expenses:     exp,
					Net:          income - exp,
					FromDate:     fromDate,
					ToDate:       toDate,
				})
			}

			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "[]")
				return nil
			}
			return printJSONFiltered(cmd.OutOrStdout(), results, flags)
		},
	}

	cmd.Flags().StringVar(&fromDate, "from", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&toDate, "to", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&property, "property", "", "Filter by property ID")
	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: ~/.local/share/doorloop-pp-cli/data.db)")
	return cmd
}
