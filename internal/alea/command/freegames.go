package command

import (
	"apikit/internal/alea/config"
	"apikit/internal/alea/console"
	"apikit/internal/alea/dto"
	"apikit/internal/alea/gql"
	cmdkit "apikit/internal/kit/command"
	"apikit/internal/kit/flag"
	"apikit/internal/kit/graphql"
	"apikit/internal/kit/middleware"
	"context"

	"github.com/urfave/cli/v3"
)

func newFreegamesCommand(module cmdkit.ModuleName) *cli.Command {
	flags := []cli.Flag{flag.Config, flag.Out, flag.Body, flag.BodySchema}
	return &cli.Command{
		Name:        "freegames",
		Usage:       "Freegames API endpoints",
		Description: "Free-spin related endpoints for Alea games.",
		Commands: []*cli.Command{
			{
				Name:   "bet-amounts",
				Usage:  "POST /api/graphql · freeSpinsCurrencies (list free-spin cost levels per currency)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesBetAmounts),
			},
		},
	}
}

func handleFreegamesBetAmounts(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.AleaConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintFreegamesBetAmountsBodySchema,
		func(v *dto.FreegamesBetAmountsVariables) {
			v.JurisdictionCode = cfg.JurisdictionCode
		},
	)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	headers := map[string]string{
		"Authorization": "Bearer " + cfg.EnvironmentToken,
		"Alea-CasinoId": cfg.CasinoId,
	}
	if cfg.EnvironmentId != "" {
		headers["Alea-EnvironmentId"] = cfg.EnvironmentId
	}

	response, err := graphql.Do(ctx, cfg.CustomerApiHost, headers, gql.FreeSpinsCurrencies, body)
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "freegames_bet_amounts", command, metadata, cmdkit.FileTypeJSON)
}
