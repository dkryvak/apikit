package command

import (
	"apikit/internal/kit/command"

	"github.com/urfave/cli/v3"
)

const Module = command.ModuleName("alea")

// NewConfigCommand returns the mode-neutral `alea` command exposing only
// config management (apikit alea config ...).
func NewConfigCommand() *cli.Command {
	return &cli.Command{
		Name:        "alea",
		Usage:       "Manage Alea API configs",
		Description: "Create and manage saved Alea configs",
		Commands: []*cli.Command{
			newConfigCommand(Module),
		},
	}
}

// NewCallCommand returns the `alea` command exposing only call endpoints.
// It is mounted under both `kube` and `local` namespaces.
func NewCallCommand() *cli.Command {
	return &cli.Command{
		Name:        "alea",
		Usage:       "Run an Alea API request",
		Description: "Use subcommands to call specific endpoints",
		Commands: []*cli.Command{
			{
				Name:        "call",
				Usage:       "Run an API request",
				Description: "Use subcommands to call specific endpoints",
				Commands: []*cli.Command{
					newFreegamesCommand(Module),
				},
			},
		},
	}
}
