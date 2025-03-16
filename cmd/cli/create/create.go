package create

import (
	"github.com/spf13/cobra"
)

func NewCreateCommands() *cobra.Command {
	createCmdRoot := &cobra.Command{
		Use: "create",
	}

	createCmdRoot.AddCommand(
		NewCreateChannelsCommand(),
		NewCreateRolesCommand(),
	)

	return createCmdRoot
}
