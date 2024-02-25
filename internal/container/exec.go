package container

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const waitDuration = 75 * time.Millisecond

// Exec executes a command in a container and waits for it to finish.
func Exec(ctx context.Context, api client.APIClient, containerID string, cmd ...string) error {
	execID, err := api.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd: cmd,
	})
	if err != nil {
		return fmt.Errorf("failed to create exec: %w", err)
	}

	if err = api.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{}); err != nil {
		return fmt.Errorf("failed to start exec: %w", err)
	}

	for {
		inspect, err := api.ContainerExecInspect(ctx, execID.ID)
		if err != nil {
			return fmt.Errorf("failed to inspect exec: %w", err)
		}
		if !inspect.Running {
			if inspect.ExitCode != 0 {
				return fmt.Errorf("exec %s returned with exit code %d", cmd, inspect.ExitCode)
			}
			return nil
		}
		time.Sleep(waitDuration)
	}
}
