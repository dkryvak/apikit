package command

import (
	cmdkit "apikit/internal/kit/command"
	"apikit/internal/kit/flag"
	httpclient "apikit/internal/kit/http/client"
	httptypes "apikit/internal/kit/http/types"
	"apikit/internal/kit/middleware"
	"apikit/internal/redgenn/config"
	"apikit/internal/redgenn/console"
	"apikit/internal/redgenn/dto"
	"context"

	"github.com/urfave/cli/v3"
)

func newFreegamesCommand(module cmdkit.ModuleName) *cli.Command {
	flags := []cli.Flag{flag.Config, flag.Out, flag.Body, flag.BodySchema}
	return &cli.Command{
		Name:  "freegames",
		Usage: "Freegames API endpoints",
		Description: "Grant rewards (freespins, freebets, prizes) to one or multiple players.\n" +
			"First-time-player bonuses are supported by most providers.\n" +
			"Some rewards may require interacting with the game before a free bet can be granted.",
		Commands: []*cli.Command{
			{
				Name:   "bet-levels",
				Usage:  "POST /admservice (list allowed bet amounts)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesBetLevel),
			},
		},
	}
}

func handleFreegamesBetLevel(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.RedgennConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintFreegamesBetLevelBodySchema,
		func(request *dto.FreegamesBetLevelsRequest) {
			request.Login = cfg.Login
			request.Password = cfg.Password
			request.Command = "bet_level"
			request.WlCode = cfg.WlCode
		},
	)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.POST,
		Path:   "/admservice",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "freegames_bet_levels", command, metadata, cmdkit.FileTypeJSON)
}
