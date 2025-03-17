package create

import (
	"github.com/spf13/cobra"
)

func NewCommands() *cobra.Command {
	createCmdRoot := &cobra.Command{
		Use: "create",
	}

	createCmdRoot.AddCommand(
		NewCreateChannelsCommand(),
		NewCreateRolesCommand(),
	)

	return createCmdRoot
}
