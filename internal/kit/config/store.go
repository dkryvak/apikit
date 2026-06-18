package config

import (
	"apikit/internal/kit/command"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Load(alias string, module command.ModuleName, cfg Config) error {
	path, err := FilePath(alias, module)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config '%s' not found. Use 'api config create' to create it", alias)
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("invalid config %q: %w", alias, err)
	}

	cfg.SetAlias(alias)

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	for index, hook := range cfg.Hooks() {
		if err := hook(); err != nil {
			return fmt.Errorf("post-load hook #%d failed: %w", index, err)
		}
	}

	return nil
}

func Save(cfg Config, module command.ModuleName) error {
	dir, err := Dir(module)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("cannot create config dir: %w", err)
	}

	path, err := FilePath(cfg.Alias(), module)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal config: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("cannot write config %q: %w", cfg.Alias(), err)
	}

	return nil
}

func Exists(alias string, module command.ModuleName) bool {
	path, err := FilePath(alias, module)
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

func Delete(alias string, module command.ModuleName) error {
	path, err := FilePath(alias, module)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("cannot delete config %q: %w", alias, err)
	}
	return nil
}

func ListAliases(module command.ModuleName) ([]string, error) {
	dir, err := Dir(module)
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

	var aliases []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".json") {
			aliases = append(aliases, strings.TrimSuffix(name, ".json"))
		}
	}

	return aliases, nil
}
