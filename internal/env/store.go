package env

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Load reads an env by name, applies defaults and validates it.
func Load(name string) (*Env, error) {
	path, err := FilePath(name)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("env '%s' not found. Use 'apikit env create' to create it", name)
		}
		return nil, fmt.Errorf("failed to read env file: %w", err)
	}

	var e Env
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("invalid env %q: %w", name, err)
	}

	// Trust the filename as the source of truth for the name.
	e.Name = name
	e.ApplyDefaults()

	if err := e.Validate(); err != nil {
		return nil, fmt.Errorf("invalid env %q: %w", name, err)
	}

	return &e, nil
}

// Save writes an env to the store (applying defaults first).
func Save(e *Env) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("cannot create env dir: %w", err)
	}

	e.ApplyDefaults()
	if err := e.Validate(); err != nil {
		return err
	}

	path, err := FilePath(e.Name)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal env: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("cannot write env %q: %w", e.Name, err)
	}

	return nil
}

func Exists(name string) bool {
	path, err := FilePath(name)
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

func Delete(name string) error {
	path, err := FilePath(name)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("env '%s' not found", name)
		}
		return fmt.Errorf("cannot delete env %q: %w", name, err)
	}
	return nil
}

// ListNames returns all stored env names (sorted by the OS directory order).
func ListNames() ([]string, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".json") {
			names = append(names, strings.TrimSuffix(name, ".json"))
		}
	}

	return names, nil
}
