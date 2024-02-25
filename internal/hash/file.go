package hash

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// File returns the SHA256 hash of the file at the given path.
func File(path string) (string, error) {
	hash := crypto.SHA256.New()

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer f.Close()

	if _, err := io.Copy(hash, f); err != nil {
		return "", fmt.Errorf("failed to hash %s: %w", path, err)
	}

	b := hash.Sum(nil)
	return hex.EncodeToString(b), nil
}
