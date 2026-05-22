package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newLeasesListCmd(flags *rootFlags) *cobra.Command {
	var flagFilterGroup string
	var flagFilterProperty string
	var flagFilterOwner string
	var flagFilterText string
	var flagFilterStartDateFrom string
	var flagFilterStartDateTo string
	var flagFilterEndDateFrom string
	var flagFilterEndDateTo string
	var flagFilterPropertyClass string
	var flagFilterUnit string
	var flagFilterTenant string
	var flagFilterOutstandingBalanceGreaterThan string
	var flagFilterStatus string
	var flagFilterTerm string

	cmd := &cobra.Command{
		Use:         "list",
		Short:       "",
		Example:     "  doorloop-pp-cli leases list",
		Annotations: map[string]string{"pp:endpoint": "leases.list", "pp:method": "GET", "pp:path": "/leases", "mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := flags.newClient()
			if err != nil {
				return err
			}

			path := "/leases"
			params := map[string]string{}
			if flagFilterGroup != "" {
				params["filter_group"] = fmt.Sprintf("%v", flagFilterGroup)
			}
			if flagFilterProperty != "" {
				params["filter_property"] = fmt.Sprintf("%v", flagFilterProperty)
			}
			if flagFilterOwner != "" {
				params["filter_owner"] = fmt.Sprintf("%v", flagFilterOwner)
			}
			if flagFilterText != "" {
				params["filter_text"] = fmt.Sprintf("%v", flagFilterText)
			}
			if flagFilterStartDateFrom != "" {
				params["filter_start_date_from"] = fmt.Sprintf("%v", flagFilterStartDateFrom)
			}
			if flagFilterStartDateTo != "" {
				params["filter_start_date_to"] = fmt.Sprintf("%v", flagFilterStartDateTo)
			}
			if flagFilterEndDateFrom != "" {
				params["filter_end_date_from"] = fmt.Sprintf("%v", flagFilterEndDateFrom)
			}
			if flagFilterEndDateTo != "" {
				params["filter_end_date_to"] = fmt.Sprintf("%v", flagFilterEndDateTo)
			}
			if flagFilterPropertyClass != "" {
				params["filter_propertyClass"] = fmt.Sprintf("%v", flagFilterPropertyClass)
			}
			if flagFilterUnit != "" {
				params["filter_unit"] = fmt.Sprintf("%v", flagFilterUnit)
			}
			if flagFilterTenant != "" {
				params["filter_tenant"] = fmt.Sprintf("%v", flagFilterTenant)
			}
			if flagFilterOutstandingBalanceGreaterThan != "" {
				params["filter_outstandingBalanceGreaterThan"] = fmt.Sprintf("%v", flagFilterOutstandingBalanceGreaterThan)
			}
			if flagFilterStatus != "" {
				params["filter_status"] = fmt.Sprintf("%v", flagFilterStatus)
			}
			if flagFilterTerm != "" {
				params["filter_term"] = fmt.Sprintf("%v", flagFilterTerm)
			}
			data, prov, err := resolveRead(cmd.Context(), c, flags, "leases", false, path, params, nil)
			if err != nil {
				return classifyAPIError(err, flags)
			}
			// Print provenance to stderr for human-facing output only.
			// Machine-format flags (--json, --csv, --compact, --quiet, --plain,
			// --select) and piped stdout suppress this line; the JSON envelope
			// already carries meta.source for those consumers.
			// SYNC: keep this gate aligned with command_promoted.go.tmpl.
			if wantsHumanTable(cmd.OutOrStdout(), flags) {
				var countItems []json.RawMessage
				_ = json.Unmarshal(data, &countItems)
				printProvenance(cmd, len(countItems), prov)
			}
			// For JSON output, wrap with provenance envelope before passing through flags.
			// --select wins over --compact when both are set; --compact only runs when
			// no explicit fields were requested. Explicit format flags (--csv, --quiet,
			// --plain) opt out of the auto-JSON path so piped consumers that asked for
			// a non-JSON format reach the standard pipeline below.
			if flags.asJSON || (!isTerminal(cmd.OutOrStdout()) && !flags.csv && !flags.quiet && !flags.plain) {
				filtered := data
				if flags.selectFields != "" {
					filtered = filterFields(filtered, flags.selectFields)
				} else if flags.compact {
					filtered = compactFields(filtered)
				}
				wrapped, wrapErr := wrapWithProvenance(filtered, prov)
				if wrapErr != nil {
					return wrapErr
				}
				return printOutput(cmd.OutOrStdout(), wrapped, true)
			}
			// For all other output modes (table, csv, plain, quiet), use the standard pipeline
			if wantsHumanTable(cmd.OutOrStdout(), flags) {
				var items []map[string]any
				if json.Unmarshal(data, &items) == nil && len(items) > 0 {
					if err := printAutoTable(cmd.OutOrStdout(), items); err != nil {
						return err
					}
					if len(items) >= 25 {
						fmt.Fprintf(os.Stderr, "\nShowing %d results. To narrow: add --limit, --json --select, or filter flags.\n", len(items))
					}
					return nil
				}
			}
			return printOutputWithFlags(cmd.OutOrStdout(), data, flags)
		},
	}
	cmd.Flags().StringVar(&flagFilterGroup, "filter-group", "", "Filter by portfolio ID")
	cmd.Flags().StringVar(&flagFilterProperty, "filter-property", "", "Filter by property ID")
	cmd.Flags().StringVar(&flagFilterOwner, "filter-owner", "", "Filter by owner ID")
	cmd.Flags().StringVar(&flagFilterText, "filter-text", "", "Filter by lease name")
	cmd.Flags().StringVar(&flagFilterStartDateFrom, "filter-start-date-from", "", "Start date range from (YYYY-MM-DD)")
	cmd.Flags().StringVar(&flagFilterStartDateTo, "filter-start-date-to", "", "Start date range to (YYYY-MM-DD)")
	cmd.Flags().StringVar(&flagFilterEndDateFrom, "filter-end-date-from", "", "End date range from (YYYY-MM-DD)")
	cmd.Flags().StringVar(&flagFilterEndDateTo, "filter-end-date-to", "", "End date range to (YYYY-MM-DD)")
	cmd.Flags().StringVar(&flagFilterPropertyClass, "filter-property-class", "", "Filter by property class (RESIDENTIAL or COMMERCIAL)")
	cmd.Flags().StringVar(&flagFilterUnit, "filter-unit", "", "Filter by unit ID")
	cmd.Flags().StringVar(&flagFilterTenant, "filter-tenant", "", "Filter by tenant ID")
	cmd.Flags().StringVar(&flagFilterOutstandingBalanceGreaterThan, "filter-outstanding-balance-greater-than", "", "Filter by minimum outstanding balance")
	cmd.Flags().StringVar(&flagFilterStatus, "filter-status", "", "Filter by status (ACTIVE or INACTIVE)")
	cmd.Flags().StringVar(&flagFilterTerm, "filter-term", "", "Filter by term (Rollover or AtWill)")

	return cmd
}
