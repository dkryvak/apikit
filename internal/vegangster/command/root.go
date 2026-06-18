package command

import (
	"apikit/internal/kit/command"

	"github.com/urfave/cli/v3"
)

const Module = command.ModuleName("vegangster")

// NewConfigCommand returns the mode-neutral `vegangster` command exposing only
// config management (apikit vegangster config ...).
func NewConfigCommand() *cli.Command {
	return &cli.Command{
		Name:        "vegangster",
		Usage:       "Manage Vegangster API configs",
		Description: "Create and manage saved Vegangster configs",
		Commands: []*cli.Command{
			newConfigCommand(Module),
		},
	}
}

// NewCallCommand returns the `vegangster` command exposing only call endpoints.
// It is mounted under both `kube` and `local` namespaces.
func NewCallCommand() *cli.Command {
	return &cli.Command{
		Name:        "vegangster",
		Usage:       "Run a Vegangster API request",
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
