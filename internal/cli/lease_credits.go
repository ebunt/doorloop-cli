package cli

import (
	"github.com/spf13/cobra"
)

func newLeaseCreditsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "lease-credits",
		Short:  "Credits applied to leases",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newLeaseCreditsCreateCmd(flags))
	cmd.AddCommand(newLeaseCreditsDeleteCmd(flags))
	cmd.AddCommand(newLeaseCreditsGetCmd(flags))
	cmd.AddCommand(newLeaseCreditsListCmd(flags))
	cmd.AddCommand(newLeaseCreditsUpdateCmd(flags))
	return cmd
}
