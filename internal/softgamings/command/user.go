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

func newUserCommand(module cmdkit.ModuleName) *cli.Command {
	return &cli.Command{
		Name:        "user",
		Usage:       "User API endpoints",
		Description: "Endpoints for retrieving and managing user data",
		Commands: []*cli.Command{
			{
				Name:   "get",
				Usage:  "GET /{KEY}/User/GetUserData/",
				Flags:  []cli.Flag{flag.Config, flag.Out, flag.Query, flag.QuerySchema},
				Action: middleware.WithMetadataAndConfig(module, config.New, handleGetUser),
			},
		},
	}
}

func handleGetUser(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.SoftgamingsConfig) error {
	query, err := cmdkit.PrintSchemaOrParseQuery(
		command,
		console.PrintGetUserQuerySchema,
		func(query *dto.GetUserQuery) {},
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
		"User/GetUserData", serverIp, query.TID, cfg.Key, query.Login, cfg.Pwd,
	}, "/"))

	queryMap, err := mapping.StructToStringMap(query)
	if err != nil {
		return err
	}

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.GET,
		Path:   fmt.Sprintf("/%s/User/GetUserData/", cfg.Key),
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

	return cmdkit.WriteOutIfNeeded(responseBody, "user_get", command, metadata, cmdkit.FileTypeJSON)
}
