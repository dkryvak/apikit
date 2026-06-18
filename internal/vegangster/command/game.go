package command

import (
	cmdkit "apikit/internal/kit/command"
	"apikit/internal/kit/flag"
	httpclient "apikit/internal/kit/http/client"
	httptypes "apikit/internal/kit/http/types"
	"apikit/internal/kit/middleware"
	"apikit/internal/vegangster/config"
	"apikit/internal/vegangster/console"
	"apikit/internal/vegangster/crypto"
	"apikit/internal/vegangster/dto"
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
				Usage:  "POST /operator/v1/game/list (list games)",
				Flags:  []cli.Flag{flag.Config, flag.Out},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleGameList),
			},
			{
				Name:   "launch",
				Usage:  "POST /operator/v1/game/url (create game session URL)",
				Flags:  []cli.Flag{flag.Config, flag.Out, flag.Body, flag.BodySchema},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleGameLaunch),
			},
			{
				Name:   "demo-launch",
				Usage:  "POST /operator/v1/game/demo/url (create demo game session URL)",
				Flags:  []cli.Flag{flag.Config, flag.Out, flag.Body, flag.BodySchema},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleGameDemoLaunch),
			},
		},
	}
}

func handleGameList(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.VegangsterConfig) error {
	body := dto.GameListRequest{
		OperatorID: cfg.OperatorId,
		BrandID:    cfg.BrandId,
	}

	signature, err := crypto.CreateSignature(body, cfg.PrivateKey)
	if err != nil {
		return err
	}

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.POST,
		Path:   "/operator/v1/game/list",
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"X-API-Signature": signature,
		},
		Body: body,
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "game_list", command, metadata, cmdkit.FileTypeJSON)
}

func handleGameLaunch(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.VegangsterConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintGameLaunchBodySchema,
		func(request *dto.GameLaunchRequest) {
			request.OperatorID = cfg.OperatorId
			request.BrandID = cfg.BrandId
		},
	)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	signature, err := crypto.CreateSignature(body, cfg.PrivateKey)
	if err != nil {
		return err
	}

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.POST,
		Path:   "/operator/v1/game/url",
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"X-API-Signature": signature,
		},
		Body: body,
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "game_launch", command, metadata, cmdkit.FileTypeJSON)
}

func handleGameDemoLaunch(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.VegangsterConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintGameDemoLaunchBodySchema,
		func(request *dto.GameDemoLaunchRequest) {
			request.OperatorID = cfg.OperatorId
			request.BrandID = cfg.BrandId
		},
	)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	signature, err := crypto.CreateSignature(body, cfg.PrivateKey)
	if err != nil {
		return err
	}

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.POST,
		Path:   "/operator/v1/game/demo/url",
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"X-API-Signature": signature,
		},
		Body: body,
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "game_demo_launch", command, metadata, cmdkit.FileTypeJSON)
}
