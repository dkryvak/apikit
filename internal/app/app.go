package app

import (
	env "apikit/internal/env/command"
	"apikit/internal/meta"
	redgenn "apikit/internal/redgenn/command"
	softgamings "apikit/internal/softgamings/command"
	vegangster "apikit/internal/vegangster/command"

	"github.com/urfave/cli/v3"
)

func NewApp(version string) *cli.Command {
	return &cli.Command{
		Name:                  meta.AppName,
		Usage:                 "A CLI tool for making API requests",
		Version:               version,
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			// Mode-neutral management: shared envs and per-module configs.
			env.NewRootCommand(),
			vegangster.NewConfigCommand(),
			redgenn.NewConfigCommand(),
			softgamings.NewConfigCommand(),
			// Execution namespaces: the mode (kube|local) decides how calls run.
			newKubeCommand(),
			newLocalCommand(),
		},
	}
}

// newKubeCommand wraps module call subtrees to run via the Kubernetes cluster
// (whitelisted egress IP). Self-remoting is wired in Phase 3.
func newKubeCommand() *cli.Command {
	return &cli.Command{
		Name:        "kube",
		Usage:       "Run API requests via the Kubernetes cluster (whitelisted IP)",
		Description: "Executes the call inside a cluster pod; cluster is taken from the config's bound env",
		Commands: []*cli.Command{
			vegangster.NewCallCommand(),
			redgenn.NewCallCommand(),
			softgamings.NewCallCommand(),
		},
	}
}

// newLocalCommand wraps module call subtrees to run directly from this machine.
func newLocalCommand() *cli.Command {
	return &cli.Command{
		Name:        "local",
		Usage:       "Run API requests directly from this machine",
		Description: "Executes the call locally against the config's apiHost (env is ignored)",
		Commands: []*cli.Command{
			vegangster.NewCallCommand(),
			redgenn.NewCallCommand(),
			softgamings.NewCallCommand(),
		},
	}
}
