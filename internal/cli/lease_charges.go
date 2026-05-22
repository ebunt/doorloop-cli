package cli

import (
	"github.com/spf13/cobra"
)

func newLeaseChargesCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "lease-charges",
		Short:  "Charges applied to leases (rent, fees, etc.)",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newLeaseChargesCreateCmd(flags))
	cmd.AddCommand(newLeaseChargesDeleteCmd(flags))
	cmd.AddCommand(newLeaseChargesGetCmd(flags))
	cmd.AddCommand(newLeaseChargesListCmd(flags))
	cmd.AddCommand(newLeaseChargesUpdateCmd(flags))
	return cmd
}
