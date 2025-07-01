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
		Long:  "create to role in discord server",
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

	cmd.Flags().StringVarP(&o.SheetID, "sheet-id", "s", o.SheetID, "set sheetID")
	cmd.Flags().StringVarP(&o.Range, "range", "r", o.Range, "set range")
	cmd.Flags().StringVarP(&o.FilePath, "file", "f", o.FilePath, "set file")
	cmd.Flags().StringVarP(&o.EnvFilePath, "env-file", "e", o.FilePath, "set env file (e.g. -e .env)")

	return cmd
}
