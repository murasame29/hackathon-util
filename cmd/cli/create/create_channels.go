package create

import (
	"github.com/murasame29/hackathon-util/internal/application/create"
	"github.com/spf13/cobra"
)

func NewCreateChannelsCommand() *cobra.Command {
	o := create.NewCreateChannelsOptions()

	cmd := &cobra.Command{
		Use:   "create channels (-f FILE | -u URL)",
		Short: "create to channels",
		Long:  "create to channels (会議,雑談,VC) in discord server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&o.URL, "u", o.URL, "set url")
	cmd.Flags().StringVar(&o.FilePath, "f", o.FilePath, "set file")

	return cmd
}
