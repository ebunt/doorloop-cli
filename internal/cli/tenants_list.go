package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newTenantsListCmd(flags *rootFlags) *cobra.Command {
	var flagFilterGroup string
	var flagFilterProperty string
	var flagFilterLease string
	var flagFilterText string
	var flagFilterUnit string
	var flagFilterType string

	cmd := &cobra.Command{
		Use:         "list",
		Short:       "",
		Example:     "  doorloop-pp-cli tenants list",
		Annotations: map[string]string{"pp:endpoint": "tenants.list", "pp:method": "GET", "pp:path": "/tenants", "mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := flags.newClient()
			if err != nil {
				return err
			}

			path := "/tenants"
			params := map[string]string{}
			if flagFilterGroup != "" {
				params["filter_group"] = fmt.Sprintf("%v", flagFilterGroup)
			}
			if flagFilterProperty != "" {
				params["filter_property"] = fmt.Sprintf("%v", flagFilterProperty)
			}
			if flagFilterLease != "" {
				params["filter_lease"] = fmt.Sprintf("%v", flagFilterLease)
			}
			if flagFilterText != "" {
				params["filter_text"] = fmt.Sprintf("%v", flagFilterText)
			}
			if flagFilterUnit != "" {
				params["filter_unit"] = fmt.Sprintf("%v", flagFilterUnit)
			}
			if flagFilterType != "" {
				params["filter_type"] = fmt.Sprintf("%v", flagFilterType)
			}
			data, prov, err := resolveRead(cmd.Context(), c, flags, "tenants", false, path, params, nil)
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
	cmd.Flags().StringVar(&flagFilterLease, "filter-lease", "", "Filter by lease ID")
	cmd.Flags().StringVar(&flagFilterText, "filter-text", "", "Filter by tenant name or email")
	cmd.Flags().StringVar(&flagFilterUnit, "filter-unit", "", "Filter by unit ID")
	cmd.Flags().StringVar(&flagFilterType, "filter-type", "", "Filter by type (LEASE_TENANT or PROSPECT_TENANT)")

	return cmd
}
