package command

import (
	cmdkit "apikit/internal/kit/command"
	"apikit/internal/kit/flag"
	httpclient "apikit/internal/kit/http/client"
	httptypes "apikit/internal/kit/http/types"
	"apikit/internal/kit/ip"
	"apikit/internal/kit/mapping"
	"apikit/internal/kit/middleware"
	"apikit/internal/softgamings/config"
	"apikit/internal/softgamings/console"
	"apikit/internal/softgamings/crypto"
	"apikit/internal/softgamings/dto"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/urfave/cli/v3"
)

func newFreegamesCommand(module cmdkit.ModuleName) *cli.Command {
	return &cli.Command{
		Name:  "freegames",
		Usage: "Freegames API endpoints",
		Description: "Grant rewards (freespins, freebets, prizes) to one or multiple players.\n" +
			"First-time-player bonuses are supported by most providers.\n" +
			"Some rewards may require interacting with the game before a free bet can be granted.",
		Commands: []*cli.Command{
			{
				Name:   "grant",
				Usage:  "GET /{KEY}/Freerounds/Add/ (grant freegames bonus)",
				Flags:  []cli.Flag{flag.Config, flag.Out, flag.Query, flag.QuerySchema},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesGrant),
			},
			{
				Name:   "cancel",
				Usage:  "GET /{KEY}/Freerounds/Remove/ (cancel freegames bonus)",
				Flags:  []cli.Flag{flag.Config, flag.Out, flag.Query, flag.QuerySchema},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesCancel),
			},
		},
	}
}

func handleFreegamesGrant(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.SoftgamingsConfig) error {
	query, err := cmdkit.PrintSchemaOrParseQuery(
		command,
		console.PrintFreegamesGrantQuerySchema,
		func(query *dto.FreegamesGrantQuery) {},
	)
	if err != nil {
		return err
	}
	if query == nil {
		return nil
	}

	serverIp, err := ip.GetCurrentIp(ctx)
	if err != nil {
		return err
	}

	query.Hash = crypto.Md5Checksum(strings.Join([]string{
		query.Operator + "/Freerounds", serverIp, query.TID, cfg.Key, cfg.Pwd}, "/",
	))

	queryMap, err := mapping.StructToStringMap(query)
	if err != nil {
		return err
	}

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.GET,
		Path:   fmt.Sprintf("/%s/Freerounds/Add/", cfg.Key),
		Query:  queryMap,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return err
	}

	responseBody := response.Body
	if len(responseBody) > 2 {
		responseBody = responseBody[2:]
	}

	return cmdkit.WriteOutIfNeeded(responseBody, "freegames_grant", command, metadata, cmdkit.FileTypeJSON)
}

func handleFreegamesCancel(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.SoftgamingsConfig) error {
	query, err := cmdkit.PrintSchemaOrParseQuery(
		command,
		console.PrintFreegamesCancelQuerySchema,
		func(query *dto.FreegamesCancelQuery) {},
	)
	if err != nil {
		return err
	}
	if query == nil {
		return nil
	}

	serverIp, err := ip.GetCurrentIp(ctx)
	if err != nil {
		return err
	}

	query.TID = strings.ReplaceAll(uuid.New().String(), "-", "")
	query.Hash = crypto.Md5Checksum(strings.Join([]string{
		query.Operator + "/Freerounds", serverIp, query.TID, cfg.Key, cfg.Pwd}, "/",
	))

	queryMap, err := mapping.StructToStringMap(query)
	if err != nil {
		return err
	}

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.GET,
		Path:   fmt.Sprintf("/%s/Freerounds/Remove/", cfg.Key),
		Query:  queryMap,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return err
	}

	responseBody := response.Body
	if len(responseBody) > 2 {
		responseBody = responseBody[2:]
	}

	return cmdkit.WriteOutIfNeeded(responseBody, "freegames_cancel", command, metadata, cmdkit.FileTypeJSON)
}
