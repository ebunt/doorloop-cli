package cli

import (
	"github.com/spf13/cobra"
)

func newLeaseReturnedPaymentsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "lease-returned-payments",
		Short:  "Returned/bounced payments on leases",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newLeaseReturnedPaymentsCreateCmd(flags))
	cmd.AddCommand(newLeaseReturnedPaymentsDeleteCmd(flags))
	cmd.AddCommand(newLeaseReturnedPaymentsGetCmd(flags))
	cmd.AddCommand(newLeaseReturnedPaymentsListCmd(flags))
	cmd.AddCommand(newLeaseReturnedPaymentsUpdateCmd(flags))
	return cmd
}
