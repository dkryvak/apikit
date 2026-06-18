package env

import (
	"apikit/internal/meta"
	"fmt"
	"os"
	"path/filepath"
)

// Dir returns the shared env store directory: ~/.apikit/env
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot resolve user home dir: %w", err)
	}
	return filepath.Join(home, meta.AppDirName(), "env"), nil
}

// FilePath returns the path to a single env file: ~/.apikit/env/<name>.json
func FilePath(name string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".json"), nil
}
