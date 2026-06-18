package command

import (
	"apikit/internal/kit/command"

	"github.com/urfave/cli/v3"
)

const Module = command.ModuleName("softgamings")

// NewConfigCommand returns the mode-neutral `softgamings` command exposing only
// config management (apikit softgamings config ...).
func NewConfigCommand() *cli.Command {
	return &cli.Command{
		Name:        "softgamings",
		Usage:       "Manage Softgamings API configs",
		Description: "Create and manage saved Softgamings configs",
		Commands: []*cli.Command{
			newConfigCommand(Module),
		},
	}
}

// NewCallCommand returns the `softgamings` command exposing only call endpoints.
// It is mounted under both `kube` and `local` namespaces.
func NewCallCommand() *cli.Command {
	return &cli.Command{
		Name:        "softgamings",
		Usage:       "Run a Softgamings API request",
		Description: "Use subcommands to call specific endpoints",
		Commands: []*cli.Command{
			{
				Name:        "call",
				Usage:       "Run an API request",
				Description: "Use subcommands to call specific endpoints",
				Commands: []*cli.Command{
					newGameCommand(Module),
					newFreegamesCommand(Module),
					newUserCommand(Module),
				},
			},
		},
	}
}
