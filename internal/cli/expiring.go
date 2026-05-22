package cli

import (
	"fmt"
	"strings"
	"time"

	"doorloop-pp-cli/internal/store"

	"github.com/spf13/cobra"
)

type expiringRow struct {
	LeaseID     string `json:"lease_id"`
	EndDate     string `json:"end_date"`
	DaysLeft    int    `json:"days_left"`
	UnitID      string `json:"unit_id"`
	PropertyID  string `json:"property_id"`
	TenantName  string `json:"tenant_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	HasProspect bool   `json:"has_prospect"`
}

func newExpiringCmd(flags *rootFlags) *cobra.Command {
	var days int
	var property string
	var dbPath string

	cmd := &cobra.Command{
		Use:   "expiring",
		Short: "List leases expiring within N days with tenant contact info",
		Long: `Find all active leases expiring within N days, with tenant contact info and
whether a prospect tenant already exists for the unit.
Queries local SQLite — run 'sync' first.`,
		Example: `  doorloop-pp-cli expiring
  doorloop-pp-cli expiring --days 60 --json
  doorloop-pp-cli expiring --property prop_123 --json --agent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbPath == "" {
				dbPath = defaultDBPath("doorloop-pp-cli")
			}
			db, err := store.OpenReadOnly(dbPath)
			if err != nil {
				return fmt.Errorf("opening local database (run 'sync' first): %w", err)
			}
			defer db.Close()

			today := time.Now()
			fromStr := today.Format("2006-01-02")
			toStr := today.AddDate(0, 0, days).Format("2006-01-02")

			q := `
				SELECT
					l.id AS lease_id,
					COALESCE(json_extract(l.data, '$.endDate'), json_extract(l.data, '$.end_date'), '') AS end_date,
					COALESCE(json_extract(l.data, '$.unitId'), json_extract(l.data, '$.unit_id'), '') AS unit_id,
					COALESCE(json_extract(l.data, '$.propertyId'), json_extract(l.data, '$.property_id'), '') AS property_id,
					COALESCE(json_extract(t.data, '$.firstName'), '') || ' ' || COALESCE(json_extract(t.data, '$.lastName'), '') AS tenant_name,
					COALESCE(json_extract(t.data, '$.email'), '') AS email,
					COALESCE(json_extract(t.data, '$.phoneNumber'), json_extract(t.data, '$.phone'), '') AS phone,
					CASE WHEN EXISTS (
						SELECT 1 FROM resources pt WHERE pt.resource_type = 'tenants'
						AND COALESCE(json_extract(pt.data, '$.type'), '') = 'PROSPECT_TENANT'
						AND (json_extract(pt.data, '$.unitId') = json_extract(l.data, '$.unitId')
						     OR json_extract(pt.data, '$.unit_id') = json_extract(l.data, '$.unit_id'))
					) THEN 1 ELSE 0 END AS has_prospect
				FROM resources l
				LEFT JOIN resources t ON t.resource_type = 'tenants'
					AND (json_extract(t.data, '$.leaseId') = l.id OR json_extract(t.data, '$.lease_id') = l.id)
					AND COALESCE(json_extract(t.data, '$.type'), '') != 'PROSPECT_TENANT'
				WHERE l.resource_type = 'leases'
				  AND COALESCE(json_extract(l.data, '$.status'), '') = 'ACTIVE'
				  AND COALESCE(json_extract(l.data, '$.endDate'), json_extract(l.data, '$.end_date'), '') >= ?
				  AND COALESCE(json_extract(l.data, '$.endDate'), json_extract(l.data, '$.end_date'), '') <= ?`

			qargs := []any{fromStr, toStr}
			if property != "" {
				q += ` AND (json_extract(l.data, '$.propertyId') = ? OR json_extract(l.data, '$.property_id') = ?)`
				qargs = append(qargs, property, property)
			}
			q += ` ORDER BY end_date ASC`

			rows, err := db.Query(q, qargs...)
			if err != nil {
				return fmt.Errorf("querying expiring leases: %w", err)
			}
			defer rows.Close()

			var results []expiringRow
			for rows.Next() {
				var r expiringRow
				var hasProspect int
				if err := rows.Scan(&r.LeaseID, &r.EndDate, &r.UnitID, &r.PropertyID, &r.TenantName, &r.Email, &r.Phone, &hasProspect); err != nil {
					continue
				}
				r.TenantName = strings.TrimSpace(r.TenantName)
				r.HasProspect = hasProspect == 1
				if len(r.EndDate) >= 10 {
					if t, err := time.Parse("2006-01-02", r.EndDate[:10]); err == nil {
						r.DaysLeft = int(t.Sub(today).Hours() / 24)
					}
				}
				results = append(results, r)
			}
			if err := rows.Err(); err != nil {
				return fmt.Errorf("reading expiring leases: %w", err)
			}

			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "[]")
				return nil
			}
			return printJSONFiltered(cmd.OutOrStdout(), results, flags)
		},
	}

	cmd.Flags().IntVar(&days, "days", 30, "Number of days ahead to look for expiring leases")
	cmd.Flags().StringVar(&property, "property", "", "Filter by property ID")
	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: ~/.local/share/doorloop-pp-cli/data.db)")
	return cmd
}
