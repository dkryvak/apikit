package command

import (
	cmdkit "apikit/internal/kit/command"
	"apikit/internal/kit/flag"
	httpclient "apikit/internal/kit/http/client"
	httptypes "apikit/internal/kit/http/types"
	"apikit/internal/kit/mapping"
	"apikit/internal/kit/middleware"
	"apikit/internal/redgenn/config"
	"apikit/internal/redgenn/console"
	"apikit/internal/redgenn/dto"
	"context"

	"github.com/urfave/cli/v3"
)

func newGameCommand(module cmdkit.ModuleName) *cli.Command {
	return &cli.Command{
		Name:        "game",
		Usage:       "Games API endpoints",
		Description: "Endpoints for listing games and generating game URLs",
		Commands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "POST /admservice (list games)",
				Flags:  []cli.Flag{flag.Config, flag.Out, flag.Body, flag.BodySchema},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleGameList),
			},
			{
				Name:   "launch",
				Usage:  "GET /game-launch-url (create game session URL)",
				Flags:  []cli.Flag{flag.Config, flag.Out, flag.Body, flag.BodySchema},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleGameLaunch),
			},
			{
				Name:   "demo-launch",
				Usage:  "GET /game-launch-url (create demo game session URL)",
				Flags:  []cli.Flag{flag.Config, flag.Out, flag.Body, flag.BodySchema},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleGameDemoLaunch),
			},
		},
	}
}

func handleGameList(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.RedgennConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintGameListBodySchema,
		func(request *dto.GameListRequest) {
			request.Login = cfg.Login
			request.Password = cfg.Password
			request.Command = "games"
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

	return cmdkit.WriteOutIfNeeded(response.Body, "game_list", command, metadata, cmdkit.FileTypeJSON)
}

func handleGameLaunch(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.RedgennConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintGameLaunchBodySchema,
		func(request *dto.GameLaunchRequest) {
			request.Mode = "real"
			request.WlCode = cfg.WlCode
		},
	)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	queryParams, err := mapping.StructToStringMap(body)
	if err != nil {
		return err
	}

	response, err := httpclient.NewHttpClient(cfg.LaunchApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.GET,
		Path:   "/game-launch-url",
		Query:  queryParams,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "game_launch", command, metadata, cmdkit.FileTypeJSON)
}

func handleGameDemoLaunch(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.RedgennConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintGameDemoLaunchBodySchema,
		func(request *dto.GameLaunchRequest) {
			request.Mode = "demo"
			request.WlCode = cfg.WlCode
			request.Token = "demo"
		},
	)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	queryParams, err := mapping.StructToStringMap(body)
	if err != nil {
		return err
	}

	response, err := httpclient.NewHttpClient(cfg.LaunchApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.GET,
		Path:   "/game-launch-url",
		Query:  queryParams,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "game_demo_launch", command, metadata, cmdkit.FileTypeJSON)
}
