package cli

import (
	"github.com/spf13/cobra"
)

func newTasksCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "tasks",
		Short:  "Maintenance and operational tasks",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newTasksAgendaCmd(flags))
	cmd.AddCommand(newTasksCreateCmd(flags))
	cmd.AddCommand(newTasksDeleteCmd(flags))
	cmd.AddCommand(newTasksGetCmd(flags))
	cmd.AddCommand(newTasksListCmd(flags))
	cmd.AddCommand(newTasksPostUpdateCmd(flags))
	cmd.AddCommand(newTasksUpdateCmd(flags))
	return cmd
}
