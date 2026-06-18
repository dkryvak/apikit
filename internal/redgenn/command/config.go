package command

import (
	cmdkit "apikit/internal/kit/command"
	cfgkit "apikit/internal/kit/config"
	"apikit/internal/kit/middleware"
	"apikit/internal/kit/validator"
	"apikit/internal/meta"
	"apikit/internal/redgenn/config"
	"apikit/internal/redgenn/console"
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

func newConfigCommand(module cmdkit.ModuleName) *cli.Command {
	return &cli.Command{
		Name:        "config",
		Usage:       "Manage API configurations",
		Description: "Create and manage saved API configs (base URL, credentials, and environment-specific settings)",
		Commands: []*cli.Command{
			{
				Name:        "create",
				Usage:       "Create a new config interactively",
				Description: "Starts an interactive wizard and saves the config under a chosen alias",
				Action:      middleware.WithMetadata(module, handleConfigCreate),
			},
			{
				Name:        "import",
				Usage:       "Import config from JSON",
				Description: "Create a new configuration by importing a JSON. Supports reading from STDIN",
				ArgsUsage:   "<alias> <json>",
				Action:      middleware.WithMetadata(module, handleConfigImport),
			},
			{
				Name:        "list",
				Usage:       "List all saved configs",
				Description: "Prints all saved config aliases",
				Action:      middleware.WithMetadata(module, handleConfigList),
			},
			{
				Name:        "show",
				Usage:       "Show details of a config",
				Description: "Displays all stored values for the selected config alias",
				ArgsUsage:   "<alias>",
				Action:      middleware.WithMetadata(module, handleConfigShow),
			},
			{
				Name:        "delete",
				Usage:       "Delete a config",
				Description: "Removes the config by alias from local storage",
				ArgsUsage:   "<alias>",
				Action:      middleware.WithMetadata(module, handleConfigDelete),
			},
		},
	}
}

func handleConfigCreate(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata) error {
	fmt.Println("Creating new API configuration")
	fmt.Println()

	cfg, err := console.PromptConfig(metadata.Module)
	if err != nil {
		return err
	}

	if err = cfgkit.Save(cfg, metadata.Module); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	console.PrintConfig(cfg)

	configPath, _ := cfgkit.FilePath(cfg.Alias(), metadata.Module)
	fmt.Printf("✓ Config '%s' saved to %s\n\n", cfg.Alias(), configPath)
	fmt.Printf("Usage: api call <endpoint-name> --config %s\n", cfg.Alias())

	return nil
}

func handleConfigImport(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata) error {
	fmt.Println("Importing API configuration")
	fmt.Println()

	args := command.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("usage: %s %s config import <alias> <json>", meta.AppName, metadata.Module)
	}
	alias := args[0]
	rawJson := strings.Join(args[1:], " ")

	var cfg config.RedgennConfig
	if err := cmdkit.ParseJson(rawJson, &cfg); err != nil {
		return err
	}

	cfg.SetAlias(alias)

	if err := validator.Struct(cfg); err != nil {
		return err
	}

	if err := cfgkit.Save(&cfg, metadata.Module); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	console.PrintConfig(&cfg)

	configPath, _ := cfgkit.FilePath(cfg.Alias(), metadata.Module)
	fmt.Printf("✓ Config '%s' saved to %s\n\n", cfg.Alias(), configPath)
	fmt.Printf("Usage: api call <endpoint-name> --config %s\n", cfg.Alias())

	return nil
}

func handleConfigList(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata) error {
	aliases, err := cfgkit.ListAliases(metadata.Module)
	if err != nil {
		return err
	}

	if len(aliases) == 0 {
		fmt.Println("No configs found.")
		fmt.Println("Create one with: api config create")
		return nil
	}

	configDir, _ := cfgkit.Dir(metadata.Module)
	fmt.Printf("Configs stored in: %s\n\n", configDir)
	fmt.Println("Available configs:")
	for _, alias := range aliases {
		fmt.Printf("  - %s\n", alias)
	}

	fmt.Println()
	fmt.Println("Use 'api config show <alias>' to view details")

	return nil
}

func handleConfigShow(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata) error {
	if command.NArg() < 1 {
		return fmt.Errorf("config alias required. Usage: api config show <alias>")
	}

	alias := command.Args().First()

	var cfg = &config.RedgennConfig{}
	if err := cfgkit.Load(alias, metadata.Module, cfg); err != nil {
		return err
	}

	console.PrintConfig(cfg)

	return nil
}

func handleConfigDelete(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata) error {
	if command.NArg() < 1 {
		return fmt.Errorf("config alias required. Usage: api config delete <alias>")
	}

	alias := command.Args().First()

	fmt.Printf("Are you sure you want to delete config '%s'? (y/n): ", alias)
	var response string
	_, _ = fmt.Scanln(&response)

	if response != "y" && response != "yes" {
		fmt.Println("Deletion cancelled")
		return nil
	}

	if err := cfgkit.Delete(alias, metadata.Module); err != nil {
		return err
	}

	fmt.Printf("✓ Config '%s' deleted\n", alias)

	return nil
}
