package command

import (
	"apikit/internal/alea/config"
	"apikit/internal/alea/console"
	"apikit/internal/alea/crypto"
	"apikit/internal/alea/dto"
	"apikit/internal/alea/gql"
	cmdkit "apikit/internal/kit/command"
	"apikit/internal/kit/flag"
	"apikit/internal/kit/graphql"
	httpclient "apikit/internal/kit/http/client"
	httptypes "apikit/internal/kit/http/types"
	"apikit/internal/kit/middleware"
	"context"
	"encoding/json"
	"fmt"

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
				Name:   "grant",
				Usage:  "POST /api/v1/free-spins (grant a free-spin bonus)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesGrant),
			},
			{
				Name:   "cancel",
				Usage:  "DELETE /api/v1/free-spins/{casinoBonusId} (cancel a free-spin bonus)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesCancel),
			},
			{
				Name:   "bet-amounts",
				Usage:  "POST /api/graphql · freeSpinsCurrencies (list free-spin cost levels per currency)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesBetAmounts),
			},
		},
	}
}

// aleaBonusHeaders builds the common headers for the signed Alea bonus REST API.
func aleaBonusHeaders(cfg *config.AleaConfig, digest string) map[string]string {
	h := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": "Bearer " + cfg.EnvironmentToken,
		"Digest":        digest,
		"Alea-CasinoId": cfg.CasinoId,
	}
	if cfg.EnvironmentId != "" {
		h["Alea-EnvironmentId"] = cfg.EnvironmentId
	}
	return h
}

func handleFreegamesGrant(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.AleaConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody[dto.FreegamesGrantRequest](
		command,
		console.PrintFreegamesGrantBodySchema,
		nil,
	)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	// Sign and send the exact same bytes so the Digest matches the body.
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	response, err := httpclient.NewHttpClient(cfg.BonusApiHost).Do(ctx, &httptypes.Request{
		Method:  httptypes.POST,
		Path:    "/api/v1/free-spins",
		Headers: aleaBonusHeaders(cfg, crypto.Digest(cfg.SecretKey, string(payload))),
		Body:    payload,
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "freegames_grant", command, metadata, cmdkit.FileTypeJSON)
}

func handleFreegamesCancel(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.AleaConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody[dto.FreegamesCancelRequest](
		command,
		console.PrintFreegamesCancelBodySchema,
		nil,
	)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	response, err := httpclient.NewHttpClient(cfg.BonusApiHost).Do(ctx, &httptypes.Request{
		Method:  httptypes.DELETE,
		Path:    "/api/v1/free-spins/" + body.CasinoBonusId,
		Headers: aleaBonusHeaders(cfg, crypto.Digest(cfg.SecretKey, body.CasinoBonusId)),
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "freegames_cancel", command, metadata, cmdkit.FileTypeJSON)
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
