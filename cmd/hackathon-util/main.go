package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/murasame29/hackathon-util/internal/command/create"
	cmddelete "github.com/murasame29/hackathon-util/internal/command/delete"
	"github.com/murasame29/hackathon-util/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	// JSON handler – structured logs parseable by New Relic / Datadog etc.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	_ = godotenv.Load()
}

func main() {
	var manifestFile string

	root := &cobra.Command{
		Use:   "hackathon-util",
		Short: "Hackathon Discord utility – create or delete team channels/roles from a Google Sheet",
	}

	// -f / --file is a persistent flag shared by all sub-commands.
	root.PersistentFlags().StringVarP(&manifestFile, "file", "f", "", "path to YAML manifest (required)")
	_ = root.MarkPersistentFlagRequired("file")

	root.AddCommand(newCreateCmd(&manifestFile))
	root.AddCommand(newDeleteCmd(&manifestFile))

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// newCreateCmd builds the "create" sub-command.
func newCreateCmd(manifestFile *string) *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Discord roles, categories, and channels from a Google Sheet",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(*manifestFile)
			if err != nil {
				return err
			}

			botToken := os.Getenv("DISCORD_BOT_TOKEN")
			if botToken == "" {
				return fmt.Errorf("DISCORD_BOT_TOKEN environment variable is required")
			}

			return create.Run(create.Config{
				BotToken: botToken,
				DryRun:   dryRun,
				Cfg:      cfg,
			})
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without applying them (reads from Discord/Sheets API but skips writes)")
	return cmd
}

// newDeleteCmd builds the "delete" sub-command.
func newDeleteCmd(manifestFile *string) *cobra.Command {
	var dryRun bool
	var removeAllMembers bool

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete Discord roles, categories, and channels listed in a Google Sheet",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(*manifestFile)
			if err != nil {
				return err
			}

			botToken := os.Getenv("DISCORD_BOT_TOKEN")
			if botToken == "" {
				return fmt.Errorf("DISCORD_BOT_TOKEN environment variable is required")
			}

			return cmddelete.Run(cmddelete.Config{
				BotToken:         botToken,
				DryRun:           dryRun,
				RemoveAllMembers: removeAllMembers,
				Cfg:              cfg,
			})
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without applying them (reads from Discord/Sheets API but skips writes)")
	cmd.Flags().BoolVar(&removeAllMembers, "remove-all-members", false, "also remove the participants role from all members")

	return cmd
}
