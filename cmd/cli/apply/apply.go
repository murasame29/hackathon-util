package apply

import (
	"github.com/murasame29/hackathon-util/internal/application/apply"
	"github.com/spf13/cobra"
)

func NewApplyCommand() *cobra.Command {
	o := apply.NewApplyOptions()

	cmd := &cobra.Command{
		Use:   `apply  (-f FILE | -s "sheetID" -r "range ("sheet!a:b")`,
		Short: "Create channels and roles, and grant roles to members.",
		Long:  "Create channels and roles, and grant roles to members.",
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
