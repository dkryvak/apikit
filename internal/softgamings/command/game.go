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

func newGameCommand(module cmdkit.ModuleName) *cli.Command {
	return &cli.Command{
		Name:        "game",
		Usage:       "Games API endpoints",
		Description: "Endpoints for listing games and generating game URLs",
		Commands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "GET /{KEY}/Game/FullList/ (list games)",
				Flags:  []cli.Flag{flag.Config, flag.Out},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleGameList),
			},
			{
				Name:   "launch",
				Usage:  "GET /{KEY}/User/AuthHTML/ (launch game)",
				Flags:  []cli.Flag{flag.Config, flag.Out, flag.Query, flag.QuerySchema},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleGameLaunch),
			},
			{
				Name:   "demo-launch",
				Usage:  "GET /{KEY}/User/AuthHTML/ (launch game)",
				Flags:  []cli.Flag{flag.Config, flag.Out, flag.Query, flag.QuerySchema},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleGameDemoLaunch),
			},
		},
	}
}

func handleGameList(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.SoftgamingsConfig) error {
	serverIp, err := ip.GetCurrentIp(ctx)
	if err != nil {
		return err
	}

	transactionId := strings.ReplaceAll(uuid.New().String(), "-", "")
	hash := crypto.Md5Checksum(strings.Join([]string{"Game/FullList", serverIp, transactionId, cfg.Key, cfg.Pwd}, "/"))

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.GET,
		Path:   fmt.Sprintf("/%s/Game/FullList/", cfg.Key),
		Query: map[string]string{
			"TID":  transactionId,
			"Hash": hash,
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "game_list", command, metadata, cmdkit.FileTypeJSON)
}

func handleGameLaunch(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.SoftgamingsConfig) error {
	query, err := cmdkit.PrintSchemaOrParseQuery(
		command,
		console.PrintGameLaunchQuerySchema,
		func(query *dto.GameLaunchQuery) {},
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

	if query.UserIP == "" {
		query.UserIP = serverIp
	}

	query.TID = strings.ReplaceAll(uuid.New().String(), "-", "")
	query.Hash = crypto.Md5Checksum(strings.Join([]string{
		"User/AuthHTML", serverIp, query.TID, cfg.Key, query.Login, query.Password, query.System, cfg.Pwd,
	}, "/"))

	queryMap, err := mapping.StructToStringMap(query)
	if err != nil {
		return err
	}

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.GET,
		Path:   fmt.Sprintf("/%s/User/AuthHTML/", cfg.Key),
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

	return cmdkit.WriteOutIfNeeded(responseBody, "game_launch", command, metadata, cmdkit.FileTypeHTML)
}

func handleGameDemoLaunch(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.SoftgamingsConfig) error {
	query, err := cmdkit.PrintSchemaOrParseQuery(
		command,
		console.PrintGameDemoLaunchQuerySchema,
		func(query *dto.GameDemoLaunchQuery) {
			query.Login = "$DemoUser$"
			query.Password = "Demo"
			query.Demo = 1
		},
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

	if query.UserIP == "" {
		query.UserIP = serverIp
	}

	query.TID = strings.ReplaceAll(uuid.New().String(), "-", "")
	query.Hash = crypto.Md5Checksum(strings.Join([]string{
		"User/AuthHTML", serverIp, query.TID, cfg.Key, query.Login, query.Password, query.System, cfg.Pwd,
	}, "/"))

	queryMap, err := mapping.StructToStringMap(query)
	if err != nil {
		return err
	}

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.GET,
		Path:   fmt.Sprintf("/%s/User/AuthHTML/", cfg.Key),
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

	return cmdkit.WriteOutIfNeeded(responseBody, "game_demo_launch", command, metadata, cmdkit.FileTypeHTML)
}
