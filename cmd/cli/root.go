package cli

import (
	"log"

	bindCmd "github.com/murasame29/hackathon-util/cmd/cli/bind"
	createCmd "github.com/murasame29/hackathon-util/cmd/cli/create"
	deleteCmd "github.com/murasame29/hackathon-util/cmd/cli/delete"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "hackathon util tools",
	Short: "",
	Long:  "",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.AddCommand(
		bindCmd.NewCommands(),
		createCmd.NewCommands(),
		deleteCmd.NewCommands(),
	)
}
