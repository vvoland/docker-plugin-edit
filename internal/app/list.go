package app

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/cli/cli/command"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
)

func List(ctx context.Context, cli command.Cli, volumeName string) error {
	api := cli.Client()

	trve := true
	create, err := api.ContainerCreate(ctx,
		&containertypes.Config{
			Image: "busybox",
			Tty:   true,
			Cmd:   []string{"tree", "/volume"},
		},
		&containertypes.HostConfig{
			Init: &trve,
			Mounts: []mount.Mount{
				{Type: mount.TypeVolume, Source: volumeName, Target: "/volume"},
			},
		},
		&network.NetworkingConfig{},
		nil, "",
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	cid := create.ID
	defer func() {
		api.ContainerRemove(ctx, cid, containertypes.RemoveOptions{Force: true})
	}()

	if err := api.ContainerStart(ctx, cid, containertypes.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	resp, err := api.ContainerLogs(ctx, cid, containertypes.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get logs: %w", err)
	}
	defer resp.Close()

	_, err = io.Copy(os.Stdout, resp)
	return err
}
