package cli

import (
	"github.com/spf13/cobra"
)

func newLeasesCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "leases",
		Short:  "Lease agreements linking tenants to units",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newLeasesGetCmd(flags))
	cmd.AddCommand(newLeasesListCmd(flags))
	cmd.AddCommand(newLeasesMoveInCmd(flags))
	cmd.AddCommand(newLeasesMoveOutCmd(flags))
	cmd.AddCommand(newLeasesTenantsCmd(flags))
	return cmd
}
