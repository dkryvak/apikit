package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func resolveOutPath(out, filename string) (string, error) {
	raw := strings.TrimSpace(out)
	if raw == "" {
		return "", fmt.Errorf("out is empty")
	}

	// IMPORTANT: check trailing separator BEFORE Clean()
	trailingSep := strings.HasSuffix(raw, "/") || strings.HasSuffix(raw, "\\")

	out = filepath.Clean(raw)

	// "." means current dir
	if out == "." {
		return filepath.Join(".", filename), nil
	}

	// user explicitly indicated a directory (even if it doesn't exist yet)
	if trailingSep {
		return filepath.Join(out, filename), nil
	}

	// existing directory
	if st, err := os.Stat(out); err == nil && st.IsDir() {
		return filepath.Join(out, filename), nil
	}

	// otherwise treat as file path
	return out, nil
}

func ensureParentDir(filePath string) error {
	dir := filepath.Dir(filePath)

	// filepath.Dir("file.json") == "." -> nothing to create
	if dir == "." {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, err)
	}
	return nil
}
