package session

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/vvoland/docker-plugin-edit/internal/container"
)

type Session struct {
	tmpDir      string
	containerID string

	api client.APIClient

	containerVolumePath string
	containerWorkPath   string
	hostWorkFile        string
}

func New(ctx context.Context, api client.APIClient, volumeName, volumePath string) (_ *Session, outErr error) {
	var sess Session

	tmpDir, err := os.MkdirTemp("", "docker-edit-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary dir: %w", err)
	}
	sess.tmpDir = tmpDir

	defer func() {
		if outErr != nil {
			sess.Close()
		}
	}()

	trve := true
	create, err := api.ContainerCreate(ctx,
		&containertypes.Config{
			Image: "busybox",
			Cmd:   []string{"sleep", "infinity"},
		},
		&containertypes.HostConfig{
			Init:       &trve,
			AutoRemove: true,
			Mounts: []mount.Mount{
				{Type: mount.TypeBind, Source: tmpDir, Target: "/work"},
				{Type: mount.TypeVolume, Source: volumeName, Target: "/volume"},
			},
		},
		&network.NetworkingConfig{},
		nil, "",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}
	sess.containerID = create.ID

	if err := api.ContainerStart(ctx, sess.containerID, containertypes.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	sess.api = api

	if err := sess.openFile(ctx, volumePath); err != nil {
		return nil, err
	}

	return &sess, nil
}

func (s *Session) openFile(ctx context.Context, volumeFilePath string) error {
	uid := os.Geteuid()
	gid := os.Getegid()

	filename := filepath.Base(volumeFilePath)

	s.containerVolumePath = filepath.Join("/volume", volumeFilePath)
	s.containerWorkPath = filepath.Join("/work", filename)

	copyAndOwn := strings.Join([]string{
		"</dev/null", ">>" + s.containerVolumePath, // Make sure file exists
		"cp", s.containerVolumePath, s.containerWorkPath, // Copy file to work path
		"&& chmod +w ", s.containerWorkPath}, // Make work file writable
		" ")
	if uid != -1 && gid != -1 {
		copyAndOwn += strings.Join([]string{"&& chown ", fmt.Sprintf("%d:%d", uid, gid), s.containerWorkPath}, " ")
	}
	if err := container.Exec(ctx, s.api, s.containerID, "sh", "-c", copyAndOwn); err != nil {
		return fmt.Errorf("failed to create exec: %w", err)
	}

	s.hostWorkFile = filepath.Join(s.tmpDir, filename)
	return nil
}

func (s *Session) HostWorkFile() string {
	return s.hostWorkFile
}

func (s *Session) Commit(ctx context.Context) error {
	return container.Exec(ctx, s.api, s.containerID, "sh", "-c",
		strings.Join([]string{
			"cat", "<" + s.containerWorkPath, ">" + s.containerVolumePath,
		}, " "))
}

func (s *Session) Close() error {
	var err error
	if s.containerID != "" {
		err = s.api.ContainerRemove(context.Background(), s.containerID, containertypes.RemoveOptions{Force: true})
	}

	if s.tmpDir != "" {
		err = os.RemoveAll(s.tmpDir)
	}

	return err
}
