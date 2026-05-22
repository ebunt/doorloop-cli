package cli

import (
	"github.com/spf13/cobra"
)

func newAccountsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "accounts",
		Short:  "Chart-of-accounts entries",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newAccountsGetCmd(flags))
	cmd.AddCommand(newAccountsListCmd(flags))
	return cmd
}
