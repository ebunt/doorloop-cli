package cli

import (
	"github.com/spf13/cobra"
)

func newUnitsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "units",
		Short:  "Individual rental units within properties",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newUnitsGetCmd(flags))
	cmd.AddCommand(newUnitsListCmd(flags))
	return cmd
}
