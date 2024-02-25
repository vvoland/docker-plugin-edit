package app

import (
	"fmt"
	"os"
	"os/exec"
)

func getEditor() (string, error) {
	editor := os.Getenv("EDITOR")
	if editor != "" {
		return editor, nil
	}

	for _, e := range []string{"vim", "vi", "nano"} {
		path, _ := exec.LookPath(e)
		if path != "" {
			return path, nil
		}
	}

	return "", fmt.Errorf("no suitable editor found")
}
