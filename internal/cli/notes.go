package cli

import (
	"github.com/spf13/cobra"
)

func newNotesCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "notes",
		Short:  "Notes attached to records",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newNotesCreateCmd(flags))
	cmd.AddCommand(newNotesDeleteCmd(flags))
	cmd.AddCommand(newNotesGetCmd(flags))
	cmd.AddCommand(newNotesListCmd(flags))
	cmd.AddCommand(newNotesUpdateCmd(flags))
	return cmd
}
