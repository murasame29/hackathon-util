package create

import (
	"github.com/murasame29/hackathon-util/internal/application/create"
	"github.com/spf13/cobra"
)

func NewCreateRolesCommand() *cobra.Command {
	o := create.NewCreateRolesOptions()

	cmd := &cobra.Command{
		Use:   `role (-f FILE | -s "sheetID" -r "range ("sheet!a:b")`,
		Short: "create to role",
		Long:  "create to role (会議,雑談,VC) in discord server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}

			if err := o.Validate(); err != nil {
				return err
			}

			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&o.SheetID, "s", o.SheetID, "set sheetID")
	cmd.Flags().StringVar(&o.Range, "r", o.Range, "set range")
	cmd.Flags().StringVar(&o.FilePath, "f", o.FilePath, "set file")

	return cmd
}
