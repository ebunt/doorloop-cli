package cli

import (
	"github.com/spf13/cobra"
)

func newCommunicationLogsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "communication-logs",
		Short:  "Logged communications with tenants and owners",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newCommunicationLogsCreateCmd(flags))
	cmd.AddCommand(newCommunicationLogsDeleteCmd(flags))
	cmd.AddCommand(newCommunicationLogsGetCmd(flags))
	cmd.AddCommand(newCommunicationLogsListCmd(flags))
	cmd.AddCommand(newCommunicationLogsUpdateCmd(flags))
	return cmd
}
