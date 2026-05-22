package cli

import (
	"github.com/spf13/cobra"
)

func newVendorBillsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "vendor-bills",
		Short:  "Bills from vendors against properties",
		Hidden: true,
		RunE:   parentNoSubcommandRunE(flags),
	}

	cmd.AddCommand(newVendorBillsCreateCmd(flags))
	cmd.AddCommand(newVendorBillsDeleteCmd(flags))
	cmd.AddCommand(newVendorBillsGetCmd(flags))
	cmd.AddCommand(newVendorBillsListCmd(flags))
	cmd.AddCommand(newVendorBillsUpdateCmd(flags))
	return cmd
}
