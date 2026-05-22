package cli

import (
	"fmt"
	"time"

	"doorloop-pp-cli/internal/store"

	"github.com/spf13/cobra"
)

type portfolioHealthRow struct {
	PortfolioID        string  `json:"portfolio_id"`
	PortfolioName      string  `json:"portfolio_name"`
	TotalUnits         int     `json:"total_units"`
	OccupiedUnits      int     `json:"occupied_units"`
	VacantUnits        int     `json:"vacant_units"`
	OccupancyRate      float64 `json:"occupancy_rate"`
	OutstandingBalance float64 `json:"outstanding_balance"`
	ExpiringSoon       int     `json:"expiring_soon"`
}

func newPortfolioHealthCmd(flags *rootFlags) *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:   "portfolio-health",
		Short: "One-row-per-portfolio: units, occupancy, balance, expiring leases",
		Long: `Aggregate per-portfolio health metrics: total units, occupancy rate, total
outstanding balance, and leases expiring within 30 days.
Queries local SQLite — run 'sync' first.`,
		Example: `  doorloop-pp-cli portfolio-health
  doorloop-pp-cli portfolio-health --json
  doorloop-pp-cli portfolio-health --json --agent`,
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
			todayStr := today.Format("2006-01-02")
			expiryThreshold := today.AddDate(0, 0, 30).Format("2006-01-02")

			q := `
				SELECT
					portfolio.id AS portfolio_id,
					COALESCE(json_extract(portfolio.data, '$.name'), portfolio.id) AS portfolio_name,
					COUNT(DISTINCT u.id) AS total_units,
					COUNT(DISTINCT CASE
						WHEN EXISTS (
							SELECT 1 FROM resources al WHERE al.resource_type = 'leases'
							AND (json_extract(al.data, '$.unitId') = u.id OR json_extract(al.data, '$.unit_id') = u.id)
							AND COALESCE(json_extract(al.data, '$.status'), '') = 'ACTIVE'
						) THEN u.id END) AS occupied_units,
					COALESCE(SUM(CAST(COALESCE(json_extract(l.data, '$.outstandingBalance'), json_extract(l.data, '$.outstanding_balance'), '0') AS REAL)), 0) AS outstanding_balance,
					COUNT(DISTINCT CASE
						WHEN COALESCE(json_extract(l.data, '$.status'), '') = 'ACTIVE'
						 AND COALESCE(json_extract(l.data, '$.endDate'), json_extract(l.data, '$.end_date'), '') >= ?
						 AND COALESCE(json_extract(l.data, '$.endDate'), json_extract(l.data, '$.end_date'), '') <= ?
						THEN l.id END) AS expiring_soon
				FROM resources portfolio
				JOIN resources p ON p.resource_type = 'properties'
					AND (json_extract(p.data, '$.portfolioId') = portfolio.id OR json_extract(p.data, '$.portfolio_id') = portfolio.id)
				JOIN resources u ON u.resource_type = 'units'
					AND (json_extract(u.data, '$.propertyId') = p.id OR json_extract(u.data, '$.property_id') = p.id)
				LEFT JOIN resources l ON l.resource_type = 'leases'
					AND (json_extract(l.data, '$.unitId') = u.id OR json_extract(l.data, '$.unit_id') = u.id)
					AND COALESCE(json_extract(l.data, '$.status'), '') = 'ACTIVE'
				WHERE portfolio.resource_type = 'portfolios'
				GROUP BY portfolio.id
				ORDER BY portfolio_name ASC`

			rows, err := db.Query(q, todayStr, expiryThreshold)
			if err != nil {
				return fmt.Errorf("querying portfolio health: %w", err)
			}
			defer rows.Close()

			var results []portfolioHealthRow
			for rows.Next() {
				var r portfolioHealthRow
				if err := rows.Scan(&r.PortfolioID, &r.PortfolioName, &r.TotalUnits, &r.OccupiedUnits, &r.OutstandingBalance, &r.ExpiringSoon); err != nil {
					continue
				}
				r.VacantUnits = r.TotalUnits - r.OccupiedUnits
				if r.TotalUnits > 0 {
					r.OccupancyRate = float64(r.OccupiedUnits) / float64(r.TotalUnits) * 100
				}
				results = append(results, r)
			}
			if err := rows.Err(); err != nil {
				return fmt.Errorf("reading portfolio health: %w", err)
			}

			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "[]")
				return nil
			}
			return printJSONFiltered(cmd.OutOrStdout(), results, flags)
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: ~/.local/share/doorloop-pp-cli/data.db)")
	return cmd
}
