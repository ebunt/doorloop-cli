package cli

import (
	"github.com/spf13/cobra"
)

func newPortfoliosCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "portfolios",
		Short:  "Portfolio groupings of properties",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newPortfoliosCreateCmd(flags))
	cmd.AddCommand(newPortfoliosDeleteCmd(flags))
	cmd.AddCommand(newPortfoliosGetCmd(flags))
	cmd.AddCommand(newPortfoliosListCmd(flags))
	cmd.AddCommand(newPortfoliosUpdateCmd(flags))
	return cmd
}
