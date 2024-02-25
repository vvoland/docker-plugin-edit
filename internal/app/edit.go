package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/cli/cli/command"
	"github.com/vvoland/docker-plugin-edit/internal/hash"
	"github.com/vvoland/docker-plugin-edit/internal/session"
)

func Edit(ctx context.Context, cli command.Cli, volumeName string, path string) error {
	api := cli.Client()

	if _, err := api.VolumeInspect(ctx, volumeName); err != nil {
		return fmt.Errorf("failed to inspect volume %s: %w", volumeName, err)
	}

	sess, err := session.New(ctx, api, volumeName, path)
	if err != nil {
		return err
	}
	defer sess.Close()

	changed, err := editOnHost(sess.HostWorkFile())
	if err != nil {
		return err
	}
	if !changed {
		fmt.Fprintln(cli.Out(), "No changes")
		return nil
	}

	fmt.Fprintf(cli.Out(), "Save? [Y/n]: ")
	var answer string
	fmt.Scanln(&answer)
	if strings.ToLower(answer) == "n" {
		fmt.Fprintln(cli.Out(), "Dropping changes")
		return nil
	}

	if err := sess.Commit(ctx); err != nil {
		return err
	}

	fmt.Fprintln(cli.Out(), "Done!")
	return nil
}

func editOnHost(path string) (changed bool, _ error) {
	before, err := hash.File(path)
	if err != nil {
		return false, err
	}

	editor, err := getEditor()
	if err != nil {
		return false, err
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed to run editor: %w", err)
	}

	after, err := hash.File(path)
	if err != nil {
		return false, err
	}

	return before != after, nil
}
