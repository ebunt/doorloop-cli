package cli

import (
	"github.com/spf13/cobra"
)

func newPropertiesCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "properties",
		Short:  "Physical properties — residential or commercial",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newPropertiesGetCmd(flags))
	cmd.AddCommand(newPropertiesListCmd(flags))
	return cmd
}
