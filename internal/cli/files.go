package cli

import (
	"github.com/spf13/cobra"
)

func newFilesCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "files",
		Short:  "Documents and file attachments",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newFilesCreateCmd(flags))
	cmd.AddCommand(newFilesDeleteCmd(flags))
	cmd.AddCommand(newFilesGetCmd(flags))
	cmd.AddCommand(newFilesListCmd(flags))
	cmd.AddCommand(newFilesUpdateCmd(flags))
	return cmd
}
