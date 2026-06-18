package command

import (
	"apikit/internal/env"
	cmdkit "apikit/internal/kit/command"
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

// NewRootCommand builds the mode-neutral `apikit env` command (CRUD only).
// Envs describe Kubernetes targets used by `apikit kube ...`; see docs/remote-design.md.
func NewRootCommand() *cli.Command {
	return &cli.Command{
		Name:        "env",
		Usage:       "Manage Kubernetes environments (clusters)",
		Description: "Create and manage shared envs (namespace, context, image) used by 'apikit kube ...'",
		Commands: []*cli.Command{
			{
				Name:        "create",
				Usage:       "Create a new env interactively",
				Description: "Starts an interactive wizard and saves the env under a chosen name",
				Action:      handleCreate,
			},
			{
				Name:        "import",
				Usage:       "Import env from JSON",
				Description: "Create a new env by importing JSON. Supports reading from STDIN with '-'",
				ArgsUsage:   "<name> <json>",
				Action:      handleImport,
			},
			{
				Name:        "list",
				Usage:       "List all saved envs",
				Description: "Prints all saved env names",
				Action:      handleList,
			},
			{
				Name:        "show",
				Usage:       "Show details of an env",
				Description: "Displays all stored values for the selected env",
				ArgsUsage:   "<name>",
				Action:      handleShow,
			},
			{
				Name:        "delete",
				Usage:       "Delete an env",
				Description: "Removes the env by name from local storage",
				ArgsUsage:   "<name>",
				Action:      handleDelete,
			},
		},
	}
}

func handleCreate(ctx context.Context, command *cli.Command) error {
	fmt.Println("Creating new environment")
	fmt.Println()

	e, err := env.PromptEnv()
	if err != nil {
		return err
	}

	if err := env.Save(e); err != nil {
		return fmt.Errorf("failed to save env: %w", err)
	}

	env.PrintEnv(e)

	path, _ := env.FilePath(e.Name)
	fmt.Printf("✓ Env '%s' saved to %s\n\n", e.Name, path)
	fmt.Printf("Usage: apikit kube <module> call <endpoint> --config <alias>  (config bound to env '%s')\n", e.Name)

	return nil
}

func handleImport(ctx context.Context, command *cli.Command) error {
	fmt.Println("Importing environment")
	fmt.Println()

	args := command.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("usage: apikit env import <name> <json>")
	}
	name := args[0]
	rawJson := strings.Join(args[1:], " ")

	var e env.Env
	if err := cmdkit.ParseJson(rawJson, &e); err != nil {
		return err
	}
	e.Name = name

	if err := env.Save(&e); err != nil {
		return fmt.Errorf("failed to save env: %w", err)
	}

	env.PrintEnv(&e)

	path, _ := env.FilePath(e.Name)
	fmt.Printf("✓ Env '%s' saved to %s\n", e.Name, path)

	return nil
}

func handleList(ctx context.Context, command *cli.Command) error {
	names, err := env.ListNames()
	if err != nil {
		return err
	}

	if len(names) == 0 {
		fmt.Println("No envs found.")
		fmt.Println("Create one with: apikit env create")
		return nil
	}

	dir, _ := env.Dir()
	fmt.Printf("Envs stored in: %s\n\n", dir)
	fmt.Println("Available envs:")
	for _, name := range names {
		fmt.Printf("  - %s\n", name)
	}

	fmt.Println()
	fmt.Println("Use 'apikit env show <name>' to view details")

	return nil
}

func handleShow(ctx context.Context, command *cli.Command) error {
	if command.NArg() < 1 {
		return fmt.Errorf("env name required. Usage: apikit env show <name>")
	}

	e, err := env.Load(command.Args().First())
	if err != nil {
		return err
	}

	env.PrintEnv(e)

	return nil
}

func handleDelete(ctx context.Context, command *cli.Command) error {
	if command.NArg() < 1 {
		return fmt.Errorf("env name required. Usage: apikit env delete <name>")
	}

	name := command.Args().First()

	fmt.Printf("Are you sure you want to delete env '%s'? (y/n): ", name)
	var response string
	_, _ = fmt.Scanln(&response)

	if response != "y" && response != "yes" {
		fmt.Println("Deletion cancelled")
		return nil
	}

	if err := env.Delete(name); err != nil {
		return err
	}

	fmt.Printf("✓ Env '%s' deleted\n", name)

	return nil
}
