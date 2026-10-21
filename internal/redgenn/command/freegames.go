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
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
				Usage:  "GET /admservice · create_offer + add_offer (grant a free-spin offer)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesGrant),
			},
			{
				Name:   "cancel",
				Usage:  "GET /admservice · remove_offer (cancel a free-spin offer by gift id)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesCancel),
			},
			{
				Name:   "bet-amounts",
				Usage:  "POST /admservice (list allowed bet amounts)",
				Flags:  flags,
				Action: middleware.WithMetadataAndConfig(module, config.New, handleFreegamesBetAmounts),
			},
		},
	}
}

func xmlHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/xml",
		"Accept":       "application/xml",
	}
}

var redgennStatusRe = regexp.MustCompile(`status="([^"]*)"`)

// ensureRedgennOk fails fast when a provider step did not return status="ok".
func ensureRedgennOk(resp *httptypes.Response, step string) error {
	if resp.StatusCode >= 400 {
		return fmt.Errorf("redgenn %s failed: HTTP %d\n%s", step, resp.StatusCode, string(resp.Body))
	}
	if m := redgennStatusRe.FindSubmatch(resp.Body); m != nil && !strings.EqualFold(string(m[1]), "ok") {
		return fmt.Errorf("redgenn %s returned status %q:\n%s", step, string(m[1]), string(resp.Body))
	}
	return nil
}

func handleFreegamesGrant(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.RedgennConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody[dto.FreegamesGrantRequest](command, console.PrintFreegamesGrantBodySchema, nil)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	client := httpclient.NewHttpClient(cfg.ApiHost)

	// Step 1: create the offer.
	createQuery := map[string]string{
		"login":    cfg.Login,
		"password": cfg.Password,
		"wlcode":   cfg.WlCode,
		"cm":       "create_offer",
		"offer":    body.Offer,
		"game":     body.Game,
		"spins":    strconv.Itoa(body.Spins),
		"totalbet": strconv.Itoa(body.Totalbet),
		"currency": body.Currency,
	}
	if body.StartLive != "" {
		createQuery["start-live"] = body.StartLive
	}
	if body.EndDateAddOffer != "" {
		createQuery["end-date-add-offer"] = body.EndDateAddOffer
	}

	createResp, err := client.Do(ctx, &httptypes.Request{
		Method:  httptypes.GET,
		Path:    "/admservice",
		Query:   createQuery,
		Headers: xmlHeaders(),
	})
	if err != nil {
		return err
	}
	if err := ensureRedgennOk(createResp, "create_offer"); err != nil {
		return err
	}

	// Step 2: assign the offer to the player.
	addQuery := map[string]string{
		"login":    cfg.Login,
		"password": cfg.Password,
		"wlcode":   cfg.WlCode,
		"cm":       "add_offer",
		"offer":    body.Offer,
		"wlid":     body.Wlid,
		"type":     "real",
	}
	if body.Lifetime != "" {
		addQuery["lifetime"] = body.Lifetime
	}
	if body.CmId != "" {
		addQuery["cm_id"] = body.CmId
	}

	addResp, err := client.Do(ctx, &httptypes.Request{
		Method:  httptypes.GET,
		Path:    "/admservice",
		Query:   addQuery,
		Headers: xmlHeaders(),
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(addResp.Body, "freegames_grant", command, metadata, cmdkit.FileTypeXML)
}

func handleFreegamesCancel(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.RedgennConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody[dto.FreegamesCancelRequest](command, console.PrintFreegamesCancelBodySchema, nil)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	response, err := httpclient.NewHttpClient(cfg.ApiHost).Do(ctx, &httptypes.Request{
		Method: httptypes.GET,
		Path:   "/admservice",
		Query: map[string]string{
			"login":    cfg.Login,
			"password": cfg.Password,
			"cm":       "remove_offer",
			"gift_id":  body.GiftId,
		},
		Headers: xmlHeaders(),
	})
	if err != nil {
		return err
	}

	return cmdkit.WriteOutIfNeeded(response.Body, "freegames_cancel", command, metadata, cmdkit.FileTypeXML)
}

func handleFreegamesBetAmounts(ctx context.Context, command *cli.Command, metadata cmdkit.Metadata, cfg *config.RedgennConfig) error {
	body, err := cmdkit.PrintSchemaOrParseBody(
		command,
		console.PrintFreegamesBetAmountsBodySchema,
		func(request *dto.FreegamesBetAmountsRequest) {
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

	return cmdkit.WriteOutIfNeeded(response.Body, "freegames_bet_amounts", command, metadata, cmdkit.FileTypeJSON)
}
