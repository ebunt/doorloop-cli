package cli

import (
	"fmt"
	"strings"
	"time"

	"doorloop-pp-cli/internal/store"

	"github.com/spf13/cobra"
)

type delinquencyRow struct {
	TenantID           string  `json:"tenant_id"`
	TenantName         string  `json:"tenant_name"`
	Email              string  `json:"email"`
	Phone              string  `json:"phone"`
	LeaseID            string  `json:"lease_id"`
	UnitID             string  `json:"unit_id"`
	PropertyID         string  `json:"property_id"`
	OutstandingBalance float64 `json:"outstanding_balance"`
	DaysPastDue        int     `json:"days_past_due"`
}

func newDelinquencyCmd(flags *rootFlags) *cobra.Command {
	var minBalance float64
	var property string
	var dbPath string

	cmd := &cobra.Command{
		Use:   "delinquency",
		Short: "List tenants with outstanding balances and contact info",
		Long: `Query synced leases for tenants with an outstanding balance above the threshold.
Joins leases and tenants from local SQLite — run 'sync' first.`,
		Example: `  doorloop-pp-cli delinquency
  doorloop-pp-cli delinquency --min-balance 100 --json
  doorloop-pp-cli delinquency --property prop_123 --json --agent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbPath == "" {
				dbPath = defaultDBPath("doorloop-pp-cli")
			}
			db, err := store.OpenReadOnly(dbPath)
			if err != nil {
				return fmt.Errorf("opening local database (run 'sync' first): %w", err)
			}
			defer db.Close()

			q := `
				SELECT
					l.id,
					COALESCE(json_extract(t.data, '$.firstName'), '') || ' ' || COALESCE(json_extract(t.data, '$.lastName'), '') AS tenant_name,
					COALESCE(json_extract(t.data, '$.email'), '') AS email,
					COALESCE(json_extract(t.data, '$.phoneNumber'), json_extract(t.data, '$.phone'), '') AS phone,
					COALESCE(t.id, '') AS tenant_id,
					COALESCE(json_extract(l.data, '$.unitId'), json_extract(l.data, '$.unit_id'), '') AS unit_id,
					COALESCE(json_extract(l.data, '$.propertyId'), json_extract(l.data, '$.property_id'), '') AS property_id,
					CAST(COALESCE(json_extract(l.data, '$.outstandingBalance'), json_extract(l.data, '$.outstanding_balance'), '0') AS REAL) AS outstanding_balance,
					COALESCE(json_extract(l.data, '$.startDate'), json_extract(l.data, '$.start_date'), '') AS start_date
				FROM resources l
				LEFT JOIN resources t ON t.resource_type = 'tenants'
					AND (json_extract(t.data, '$.leaseId') = l.id OR json_extract(t.data, '$.lease_id') = l.id)
				WHERE l.resource_type = 'leases'
				  AND COALESCE(json_extract(l.data, '$.status'), '') = 'ACTIVE'
				  AND CAST(COALESCE(json_extract(l.data, '$.outstandingBalance'), json_extract(l.data, '$.outstanding_balance'), '0') AS REAL) > ?`

			qargs := []any{minBalance}
			if property != "" {
				q += ` AND (json_extract(l.data, '$.propertyId') = ? OR json_extract(l.data, '$.property_id') = ?)`
				qargs = append(qargs, property, property)
			}
			q += ` ORDER BY outstanding_balance DESC`

			rows, err := db.Query(q, qargs...)
			if err != nil {
				return fmt.Errorf("querying delinquency data: %w", err)
			}
			defer rows.Close()

			today := time.Now()
			var results []delinquencyRow
			for rows.Next() {
				var r delinquencyRow
				var startDate string
				if err := rows.Scan(&r.LeaseID, &r.TenantName, &r.Email, &r.Phone, &r.TenantID, &r.UnitID, &r.PropertyID, &r.OutstandingBalance, &startDate); err != nil {
					continue
				}
				r.TenantName = strings.TrimSpace(r.TenantName)
				if len(startDate) >= 10 {
					if t, err := time.Parse("2006-01-02", startDate[:10]); err == nil {
						r.DaysPastDue = int(today.Sub(t).Hours() / 24)
					}
				}
				results = append(results, r)
			}
			if err := rows.Err(); err != nil {
				return fmt.Errorf("reading delinquency results: %w", err)
			}

			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "[]")
				return nil
			}
			return printJSONFiltered(cmd.OutOrStdout(), results, flags)
		},
	}

	cmd.Flags().Float64Var(&minBalance, "min-balance", 0, "Minimum outstanding balance threshold (default: 0, shows all with any balance)")
	cmd.Flags().StringVar(&property, "property", "", "Filter by property ID")
	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: ~/.local/share/doorloop-pp-cli/data.db)")
	return cmd
}
