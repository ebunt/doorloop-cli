package cli

import (
	"github.com/spf13/cobra"
)

func newOwnersCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "owners",
		Short:  "Property owners",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newOwnersCreateCmd(flags))
	cmd.AddCommand(newOwnersDeleteCmd(flags))
	cmd.AddCommand(newOwnersGetCmd(flags))
	cmd.AddCommand(newOwnersListCmd(flags))
	cmd.AddCommand(newOwnersUpdateCmd(flags))
	return cmd
}
