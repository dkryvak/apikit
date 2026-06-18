package file

import (
	"fmt"
	"os"
	"strings"
)

func WriteBytes(content []byte, out string, filename string) error {
	// "-" means stdout: emit the raw body so the output is pipeable and can be
	// captured by the kube self-remote orchestrator. No formatting, no file.
	if strings.TrimSpace(out) == "-" {
		_, err := os.Stdout.Write(content)
		return err
	}

	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("filename is empty")
	}

	filePath, err := resolveOutPath(out, filename)
	if err != nil {
		return err
	}

	// pretty by extension (json/xml), plain otherwise
	content = formatByFileExt(content, filePath)

	// create directories once (single responsibility)
	if err := ensureParentDir(filePath); err != nil {
		return err
	}

	if err := os.WriteFile(filePath, content, 0o644); err != nil {
		return fmt.Errorf("failed to write file %q: %w", filePath, err)
	}

	return nil
}
