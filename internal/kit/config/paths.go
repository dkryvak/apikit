package config

import (
	"apikit/internal/kit/command"
	"apikit/internal/meta"
	"fmt"
	"os"
	"path/filepath"
)

func Dir(module command.ModuleName) (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot resolve user config dir: %w", err)
	}
	return filepath.Join(dir, meta.AppDirName(), module, "config"), nil
}

func FilePath(alias string, module command.ModuleName) (string, error) {
	dir, err := Dir(module)
	if err != nil {
		return "", err
	}
	filename := alias + ".json"
	return filepath.Join(dir, filename), nil
}
