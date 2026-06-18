package middleware

import (
	cmdkit "apikit/internal/kit/command"
	cfgkit "apikit/internal/kit/config"
	"apikit/internal/kube"
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func WithMetadata(
	module cmdkit.ModuleName,
	fn func(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata) error,
) cli.ActionFunc {
	return func(ctx context.Context, command *cli.Command) error {
		return fn(ctx, command, cmdkit.Metadata{Module: module})
	}
}

func WithMetadataAndConfig[T cfgkit.Config](
	module cmdkit.ModuleName,
	newCfg func() T,
	fn func(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg T) error,
) cli.ActionFunc {
	return func(ctx context.Context, command *cli.Command) error {
		alias := command.String("config")
		if alias == "" {
			return fmt.Errorf("required flag \"config\" not set")
		}

		cfg := newCfg()
		if err := cfgkit.Load(alias, module, cfg); err != nil {
			return err
		}

		// Under the `kube` namespace, run the call inside the cluster instead of
		// locally — unless this is just a schema request (purely local info).
		if kube.IsRemote() && !isSchemaRequest(command) {
			return kube.SelfRemote(ctx, command, module, alias, cfg)
		}

		return fn(ctx, command, cmdkit.Metadata{Module: module, Alias: alias}, cfg)
	}
}

func isSchemaRequest(command *cli.Command) bool {
	return command.Bool("body-schema") || command.Bool("query-schema")
}
