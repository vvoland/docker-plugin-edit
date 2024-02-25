package main

import (
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
	"github.com/vvoland/docker-plugin-edit/internal/app"
)

func Command(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <volume> <file>",
		Short: "Edit a text file inside a named volume",
		Args:  cobra.ExactArgs(2),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			volumeName := args[0]
			path := args[1]

			return app.Edit(ctx, dockerCli, volumeName, path)
		},
	}

	return cmd
}

func main() {
	plugin.Run(Command, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "Paweł Gronowski",
		Version:       "0.1.0",
	})
}
