package bind

import (
	"github.com/spf13/cobra"
)

func NewCommands() *cobra.Command {
	bindCmdRoot := &cobra.Command{
		Use: "bind",
	}

	bindCmdRoot.AddCommand(
		NewBindRoleCommand(),
	)

	return bindCmdRoot
}
