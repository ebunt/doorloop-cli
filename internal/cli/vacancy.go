package cli

import (
	"fmt"
	"time"

	"doorloop-pp-cli/internal/store"

	"github.com/spf13/cobra"
)

type vacancyRow struct {
	UnitID       string `json:"unit_id"`
	UnitName     string `json:"unit_name"`
	PropertyID   string `json:"property_id"`
	PropertyName string `json:"property_name"`
	LastLeaseEnd string `json:"last_lease_end,omitempty"`
	DaysVacant   int    `json:"days_vacant"`
}

func newVacancyCmd(flags *rootFlags) *cobra.Command {
	var property string
	var dbPath string

	cmd := &cobra.Command{
		Use:   "vacancy",
		Short: "List vacant units — units with no active lease",
		Long: `Find every vacant unit across your portfolio. A unit is vacant when it has no
active lease. Shows days vacant and last lease end date.
Queries local SQLite — run 'sync' first.`,
		Example: `  doorloop-pp-cli vacancy
  doorloop-pp-cli vacancy --property prop_123 --json
  doorloop-pp-cli vacancy --json --agent`,
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
					u.id AS unit_id,
					COALESCE(json_extract(u.data, '$.name'), u.id) AS unit_name,
					COALESCE(json_extract(u.data, '$.propertyId'), json_extract(u.data, '$.property_id'), '') AS property_id,
					COALESCE(json_extract(u.data, '$.propertyName'), json_extract(u.data, '$.property_name'), '') AS property_name,
					(SELECT MAX(COALESCE(json_extract(l.data, '$.endDate'), json_extract(l.data, '$.end_date'), ''))
					 FROM resources l WHERE l.resource_type = 'leases'
					 AND (json_extract(l.data, '$.unitId') = u.id OR json_extract(l.data, '$.unit_id') = u.id)
					) AS last_lease_end
				FROM resources u
				WHERE u.resource_type = 'units'
				  AND NOT EXISTS (
					SELECT 1 FROM resources l WHERE l.resource_type = 'leases'
					AND (json_extract(l.data, '$.unitId') = u.id OR json_extract(l.data, '$.unit_id') = u.id)
					AND COALESCE(json_extract(l.data, '$.status'), '') = 'ACTIVE'
				  )`

			qargs := []any{}
			if property != "" {
				q += ` AND (json_extract(u.data, '$.propertyId') = ? OR json_extract(u.data, '$.property_id') = ?)`
				qargs = append(qargs, property, property)
			}
			q += ` ORDER BY last_lease_end ASC`

			rows, err := db.Query(q, qargs...)
			if err != nil {
				return fmt.Errorf("querying vacancy data: %w", err)
			}
			defer rows.Close()

			today := time.Now()
			var results []vacancyRow
			for rows.Next() {
				var r vacancyRow
				var lastEnd *string
				if err := rows.Scan(&r.UnitID, &r.UnitName, &r.PropertyID, &r.PropertyName, &lastEnd); err != nil {
					continue
				}
				if lastEnd != nil && *lastEnd != "" {
					r.LastLeaseEnd = *lastEnd
					if len(*lastEnd) >= 10 {
						if t, err := time.Parse("2006-01-02", (*lastEnd)[:10]); err == nil {
							r.DaysVacant = int(today.Sub(t).Hours() / 24)
						}
					}
				}
				results = append(results, r)
			}
			if err := rows.Err(); err != nil {
				return fmt.Errorf("reading vacancy results: %w", err)
			}

			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "[]")
				return nil
			}
			return printJSONFiltered(cmd.OutOrStdout(), results, flags)
		},
	}

	cmd.Flags().StringVar(&property, "property", "", "Filter by property ID")
	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: ~/.local/share/doorloop-pp-cli/data.db)")
	return cmd
}
