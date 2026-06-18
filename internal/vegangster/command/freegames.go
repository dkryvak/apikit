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
				Name:   "grant",
				Usage:  "POST /operator/v1/freegames/grant (grant freegames bonus)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesGrant),
			},
			{
				Name:   "cancel",
				Usage:  "POST /operator/v1/freegames/cancel (cancel freegames bonus)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesCancel),
			},
			{
				Name:   "bet-amounts",
				Usage:  "POST /operator/v1/freegames/bet-amounts (list allowed bet amounts)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesBetAmounts),
			},
			{
				Name:   "bet-amounts",
				Usage:  "POST /operator/v1/freegames/bet-amount/list (list allowed bet amounts by games)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesBetAmountList),
			},
		},
	}
}

func handleFreegamesGrant(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.VegangsterConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintFreegamesGrantBodySchema,
		func(request *dto.FreegamesGrantRequest) {
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
		Path:   "/operator/v1/freegames/grant",
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"X-API-Signature": signature,
		},
		Body: body,
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "freegames_grant", command, metadata, cmdkit.FileTypeJSON)
}

func handleFreegamesCancel(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.VegangsterConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintFreegamesCancelBodySchema,
		func(request *dto.FreegamesCancelRequest) {
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
		Path:   "/operator/v1/freegames/cancel",
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"X-API-Signature": signature,
		},
		Body: body,
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "freegames_cancel", command, metadata, cmdkit.FileTypeJSON)
}

func handleFreegamesBetAmounts(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.VegangsterConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintFreegamesBetAmountsBodySchema,
		func(request *dto.FreegamesBetAmountsRequest) {
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
		Path:   "/operator/v1/freegames/bet-amounts",
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"X-API-Signature": signature,
		},
		Body: body,
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "freegames_bet_amounts", command, metadata, cmdkit.FileTypeJSON)
}

func handleFreegamesBetAmountList(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.VegangsterConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintFreegamesBetAmountListBodySchema,
		func(request *dto.FreegamesBetAmountListRequest) {
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
		Path:   "/operator/v1/freegames/bet-amount/list",
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"X-API-Signature": signature,
		},
		Body: body,
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "freegames_bet_amount_list", command, metadata, cmdkit.FileTypeJSON)
}
