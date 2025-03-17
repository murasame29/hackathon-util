package delete

import (
	deleteApp "github.com/murasame29/hackathon-util/internal/application/delete"
	"github.com/spf13/cobra"
)

func NewDeleteRolesCommand() *cobra.Command {
	o := deleteApp.NewDeleteRolesOptions()

	cmd := &cobra.Command{
		Use:   `role (-f FILE | -s "sheetID" -r "range ("sheet!a:b")`,
		Short: "delete to role",
		Long:  "delete to role (会議,雑談,VC) in discord server",
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
