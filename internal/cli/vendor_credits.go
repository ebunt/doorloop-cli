package cli

import (
	"github.com/spf13/cobra"
)

func newVendorCreditsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "vendor-credits",
		Short:  "Credits from vendors",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newVendorCreditsCreateCmd(flags))
	cmd.AddCommand(newVendorCreditsDeleteCmd(flags))
	cmd.AddCommand(newVendorCreditsGetCmd(flags))
	cmd.AddCommand(newVendorCreditsListCmd(flags))
	cmd.AddCommand(newVendorCreditsUpdateCmd(flags))
	return cmd
}
