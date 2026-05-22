package cli

import (
	"github.com/spf13/cobra"
)

func newLeasePaymentsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "lease-payments",
		Short:  "Rent payments recorded against leases",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newLeasePaymentsCreateCmd(flags))
	cmd.AddCommand(newLeasePaymentsDeleteCmd(flags))
	cmd.AddCommand(newLeasePaymentsGetCmd(flags))
	cmd.AddCommand(newLeasePaymentsListCmd(flags))
	cmd.AddCommand(newLeasePaymentsUpdateCmd(flags))
	return cmd
}
