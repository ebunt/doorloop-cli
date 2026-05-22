package cli

import (
	"github.com/spf13/cobra"
)

func newTenantsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "tenants",
		Short:  "Tenants — active lease tenants and prospects",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newTenantsCreateCmd(flags))
	cmd.AddCommand(newTenantsDeleteCmd(flags))
	cmd.AddCommand(newTenantsGetCmd(flags))
	cmd.AddCommand(newTenantsListCmd(flags))
	cmd.AddCommand(newTenantsUpdateCmd(flags))
	return cmd
}
