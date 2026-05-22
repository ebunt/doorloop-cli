package cli

import (
	"github.com/spf13/cobra"
)

func newExpensesCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "expenses",
		Short:  "Property operating expenses",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newExpensesCreateCmd(flags))
	cmd.AddCommand(newExpensesDeleteCmd(flags))
	cmd.AddCommand(newExpensesGetCmd(flags))
	cmd.AddCommand(newExpensesListCmd(flags))
	cmd.AddCommand(newExpensesUpdateCmd(flags))
	return cmd
}
