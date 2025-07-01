package delete

import "github.com/spf13/cobra"

func NewCommands() *cobra.Command {
	deleteCmdRoot := &cobra.Command{
		Use: "delete",
	}

	deleteCmdRoot.AddCommand(
		NewDeleteChannelsCommand(),
		NewDeleteRolesCommand(),
	)

	return deleteCmdRoot
}
