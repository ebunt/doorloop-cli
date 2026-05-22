package cli

import (
	"github.com/spf13/cobra"
)

func newUsersCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "users",
		Short:  "DoorLoop users (staff)",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newUsersCreateCmd(flags))
	cmd.AddCommand(newUsersDeleteCmd(flags))
	cmd.AddCommand(newUsersGetCmd(flags))
	cmd.AddCommand(newUsersListCmd(flags))
	cmd.AddCommand(newUsersUpdateCmd(flags))
	return cmd
}
