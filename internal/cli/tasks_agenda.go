package cli

import (
	"fmt"
	"time"

	"doorloop-pp-cli/internal/store"

	"github.com/spf13/cobra"
)

type taskAgendaItem struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	DueDate      string `json:"due_date"`
	Status       string `json:"status"`
	Urgency      string `json:"urgency"`
	PropertyID   string `json:"property_id,omitempty"`
	AssigneeName string `json:"assignee_name,omitempty"`
}

func newTasksAgendaCmd(flags *rootFlags) *cobra.Command {
	var property string
	var assignee string
	var dbPath string

	cmd := &cobra.Command{
		Use:   "agenda",
		Short: "Today's tasks bucketed as OVERDUE, DUE_TODAY, and DUE_THIS_WEEK",
		Long: `Show tasks with due dates bucketed by urgency: OVERDUE, DUE_TODAY, and DUE_THIS_WEEK.
Filter by property or assignee. Queries local SQLite — run 'sync' first.`,
		Example: `  doorloop-pp-cli tasks agenda
  doorloop-pp-cli tasks agenda --property prop_123 --json
  doorloop-pp-cli tasks agenda --assignee user_789 --json --agent`,
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
			weekStr := today.AddDate(0, 0, 7).Format("2006-01-02")

			q := `
				SELECT
					id,
					COALESCE(json_extract(data, '$.title'), '') AS title,
					COALESCE(json_extract(data, '$.dueDate'), json_extract(data, '$.due_date'), '') AS due_date,
					COALESCE(json_extract(data, '$.status'), '') AS status,
					COALESCE(json_extract(data, '$.propertyId'), json_extract(data, '$.property_id'), '') AS property_id,
					COALESCE(json_extract(data, '$.assigneeName'), json_extract(data, '$.assignee_name'), '') AS assignee_name
				FROM resources
				WHERE resource_type = 'tasks'
				  AND COALESCE(json_extract(data, '$.dueDate'), json_extract(data, '$.due_date'), '') != ''
				  AND COALESCE(json_extract(data, '$.dueDate'), json_extract(data, '$.due_date'), '') <= ?
				  AND COALESCE(json_extract(data, '$.status'), '') NOT IN ('COMPLETED', 'CLOSED', 'completed', 'closed')`

			qargs := []any{weekStr}
			if property != "" {
				q += ` AND (json_extract(data, '$.propertyId') = ? OR json_extract(data, '$.property_id') = ?)`
				qargs = append(qargs, property, property)
			}
			if assignee != "" {
				q += ` AND (json_extract(data, '$.assigneeId') = ? OR json_extract(data, '$.assignee_id') = ?
				            OR json_extract(data, '$.assigneeName') = ? OR json_extract(data, '$.assignee_name') = ?)`
				qargs = append(qargs, assignee, assignee, assignee, assignee)
			}
			q += ` ORDER BY due_date ASC`

			rows, err := db.Query(q, qargs...)
			if err != nil {
				return fmt.Errorf("querying tasks: %w", err)
			}
			defer rows.Close()

			var results []taskAgendaItem
			for rows.Next() {
				var item taskAgendaItem
				if err := rows.Scan(&item.ID, &item.Title, &item.DueDate, &item.Status, &item.PropertyID, &item.AssigneeName); err != nil {
					continue
				}
				dueStr := item.DueDate
				if len(dueStr) >= 10 {
					dueStr = dueStr[:10]
				}
				switch {
				case dueStr < todayStr:
					item.Urgency = "OVERDUE"
				case dueStr == todayStr:
					item.Urgency = "DUE_TODAY"
				default:
					item.Urgency = "DUE_THIS_WEEK"
				}
				results = append(results, item)
			}
			if err := rows.Err(); err != nil {
				return fmt.Errorf("reading tasks: %w", err)
			}

			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "[]")
				return nil
			}
			return printJSONFiltered(cmd.OutOrStdout(), results, flags)
		},
	}

	cmd.Flags().StringVar(&property, "property", "", "Filter by property ID")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Filter by assignee ID or name")
	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: ~/.local/share/doorloop-pp-cli/data.db)")
	return cmd
}
