package command

import (
	"apikit/internal/kit/command"

	"github.com/urfave/cli/v3"
)

const Module = command.ModuleName("redgenn")

// NewConfigCommand returns the mode-neutral `redgenn` command exposing only
// config management (apikit redgenn config ...).
func NewConfigCommand() *cli.Command {
	return &cli.Command{
		Name:        "redgenn",
		Usage:       "Manage Redgenn API configs",
		Description: "Create and manage saved Redgenn configs",
		Commands: []*cli.Command{
			newConfigCommand(Module),
		},
	}
}

// NewCallCommand returns the `redgenn` command exposing only call endpoints.
// It is mounted under both `kube` and `local` namespaces.
func NewCallCommand() *cli.Command {
	return &cli.Command{
		Name:        "redgenn",
		Usage:       "Run a Redgenn API request",
		Description: "Use subcommands to call specific endpoints",
		Commands: []*cli.Command{
			{
				Name:        "call",
				Usage:       "Run an API request",
				Description: "Use subcommands to call specific endpoints",
				Commands: []*cli.Command{
					newGameCommand(Module),
					newFreegamesCommand(Module),
				},
			},
		},
	}
}
