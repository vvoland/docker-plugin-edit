package main

import (
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types/volume"
	"github.com/spf13/cobra"
	"github.com/vvoland/docker-plugin-edit/internal/app"
)

func Command(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <volume> [file]",
		Short: "Edit a text file inside a named volume. If no file is specified, list the contents of the volume.",
		Args:  cobra.MinimumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			ctx := cmd.Context()
			if len(args) == 0 {
				api := dockerCli.Client()
				list, err := api.VolumeList(ctx, volume.ListOptions{})
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				names := make([]string, 0, len(list.Volumes))
				for _, v := range list.Volumes {
					names = append(names, v.Name)
				}
				return names, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			volumeName := args[0]
			if len(args) <= 1 {
				return app.List(ctx, dockerCli, volumeName)
			}
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
