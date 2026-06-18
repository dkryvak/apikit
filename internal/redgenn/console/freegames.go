package console

import "apikit/internal/kit/console"

func PrintFreegamesBetLevelBodySchema() {
	console.PrintSchema(console.Schema{
		Title: "Freegames: bet-levels request body",
		Intro: "Request body must be a JSON object.",
		Fields: []console.SchemaField{
			{
				Name:        "login",
				Required:    true,
				Type:        "string",
				Description: "Admin login. Issued by the integration manager",
			},
			{
				Name:        "password",
				Required:    true,
				Type:        "string",
				Description: "Admin password.Issued by the integration manager.",
			},
			{
				Name:        "cm",
				Required:    true,
				Type:        "string",
				Description: "Command code. Required value - \"bet_level\"",
			},
			{
				Name:        "wlcode",
				Required:    true,
				Type:        "string",
				Description: "Player’s site ID. Only single wlcode value can be specified.",
			},
			{
				Name:        "currency",
				Required:    true,
				Type:        "string",
				Description: "Currency code. Multiple values can be specified, separated by commas.",
				Example:     "EUR",
			},
			{
				Name:        "game",
				Required:    true,
				Type:        "string",
				Description: "Game ID. Multiple values can be specified, separated by commas. One of the parameters game or producer must be passed",
				Example:     "bgaming_ufo_pyramids",
			},
			{
				Name:        "producer",
				Required:    false,
				Type:        "string",
				Description: "Producer name. Only single producer value can be specified. One of the parameters game or producer must be passed",
				Example:     "bgaming",
			},
		},
		Example: `{
  "currency": "EUR",
  "game": "bgaming_ufo_pyramids",
}`,
		Notes: []string{
			`Fields "login", "password", "cm" and "wlcode" are required by API, but CLI injects them from config.`,
		},
	})
}
