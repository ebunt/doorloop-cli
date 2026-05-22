package cli

import (
	"github.com/spf13/cobra"
)

func newVendorsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "vendors",
		Short:  "Vendors and contractors",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newVendorsCreateCmd(flags))
	cmd.AddCommand(newVendorsDeleteCmd(flags))
	cmd.AddCommand(newVendorsGetCmd(flags))
	cmd.AddCommand(newVendorsListCmd(flags))
	cmd.AddCommand(newVendorsUpdateCmd(flags))
	return cmd
}
