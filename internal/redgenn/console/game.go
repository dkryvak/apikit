package console

import "apikit/internal/kit/console"

func PrintGameListBodySchema() {
	console.PrintSchema(console.Schema{
		Title: "Game: list request body",
		Intro: "Actual game list can be obtained through SpinByte API using the “cm=games” command",
		Fields: []console.SchemaField{
			{
				Name:        "disabled",
				Required:    false,
				Type:        "integer",
				Description: "All games returned if not passed.\n\tPossible value: \n\t\"0\" - show only enabled games\n\t\"1\" - show only disabled games",
				Example:     "0",
			},
			{
				Name:        "producer",
				Required:    false,
				Type:        "string",
				Description: "Producer name",
			},
			{
				Name:        "provider",
				Required:    false,
				Type:        "string",
				Description: "Provider name",
			},
			{
				Name:        "promo_freespins",
				Required:    false,
				Type:        "integer",
				Description: "Possible value \"0\" or \"1\".\n\t\"1\" - show games with freebets,\n\t\"0\" - show games without freebets",
				Example:     "0",
			},
			{
				Name:        "game_id",
				Required:    false,
				Type:        "string",
				Description: "Search by unique game id (game code). Exact match only is shown",
			},
			{
				Name:     "title",
				Required: false,
				Type:     "string",
				Description: "Full-text search by game title (game name). Any matches are shown.\n\tFor example:\n\t" +
					"\"title=lucky\" will return all games whose titles contain the word \"Lucky\".\n\t" +
					"\"title=aztec fire\" will return all games whose titles contain these 2 words \"Aztec\" and “Fire”.",
				Example: "lucky",
			},
			{
				Name:        "limit",
				Required:    false,
				Type:        "integer",
				Description: "The number of entries that will be returned.\n\tIf not specified, returns the first 1000 entries",
				Example:     "10",
			},
			{
				Name:        "offset",
				Required:    false,
				Type:        "integer",
				Description: "The number of entries that will be skipped",
				Example:     "0",
			},
		},
		Example: `{
  "disabled": 0,
  "producer": "bgaming",
  "promo_freespins": "1",
  "limit": "10",
  "offset": "0"
}`,
		Notes: []string{
			`Fields "login" and "password" and "cm" are required by API, but CLI injects them from config.`,
		},
	})
}

func PrintGameLaunchBodySchema() {
	console.PrintSchema(console.Schema{
		Title: "Game: launch request body",
		Intro: "Creates a real-money game session URL for a player.",
		Fields: []console.SchemaField{
			{
				Name:        "token",
				Required:    true,
				Type:        "string",
				Description: "Player authentication token issued by the operator",
				Example:     "eyJhbGciOiJSUzI1NiJ9...",
			},
			{
				Name:        "game_code",
				Required:    true,
				Type:        "string",
				Description: "Unique game identifier code",
				Example:     "bgaming_lucky_streak",
			},
			{
				Name:        "country",
				Required:    false,
				Type:        "string",
				Description: "Player's country in ISO 3166-1 alpha-2 format",
				Example:     "US",
			},
			{
				Name:        "platform",
				Required:    false,
				Type:        "string",
				Description: "Player's platform.\n\tPossible values:\n\t\"desktop\" - desktop browser\n\t\"mobile\" - mobile browser",
				Example:     "desktop",
			},
			{
				Name:        "exit_url",
				Required:    false,
				Type:        "string",
				Description: "URL to redirect the player when they click the exit/lobby button",
				Example:     "https://example.com/lobby",
			},
			{
				Name:        "cashier_url",
				Required:    false,
				Type:        "string",
				Description: "URL to redirect the player when they click the cashier/deposit button",
				Example:     "https://example.com/deposit",
			},
			{
				Name:        "lang",
				Required:    false,
				Type:        "string",
				Description: "Game interface language in ISO 639-1 format",
				Example:     "en",
			},
		},
		Example: `{
  "token": "eyJhbGciOiJSUzI1NiJ9...",
  "game_code": "bgaming_lucky_streak",
  "country": "US",
  "platform": "desktop",
  "exit_url": "https://example.com/lobby",
  "cashier_url": "https://example.com/deposit",
  "lang": "en"
}`,
		Notes: []string{
			`Fields "mode" and "wl_code" are required by API, but CLI injects them from config.`,
		},
	})
}

func PrintGameDemoLaunchBodySchema() {
	console.PrintSchema(console.Schema{
		Title: "Game: demo launch request body",
		Intro: "Creates a demo (free-play) game session URL. No player authentication required.",
		Fields: []console.SchemaField{
			{
				Name:        "game_code",
				Required:    true,
				Type:        "string",
				Description: "Unique game identifier code",
				Example:     "bgaming_lucky_streak",
			},
			{
				Name:        "country",
				Required:    false,
				Type:        "string",
				Description: "Player's country in ISO 3166-1 alpha-2 format",
				Example:     "US",
			},
			{
				Name:        "platform",
				Required:    false,
				Type:        "string",
				Description: "Player's platform.\n\tPossible values:\n\t\"desktop\" - desktop browser\n\t\"mobile\" - mobile browser",
				Example:     "desktop",
			},
			{
				Name:        "exit_url",
				Required:    false,
				Type:        "string",
				Description: "URL to redirect the player when they click the exit/lobby button",
				Example:     "https://example.com/lobby",
			},
			{
				Name:        "cashier_url",
				Required:    false,
				Type:        "string",
				Description: "URL to redirect the player when they click the cashier/deposit button",
				Example:     "https://example.com/deposit",
			},
			{
				Name:        "lang",
				Required:    false,
				Type:        "string",
				Description: "Game interface language in ISO 639-1 format",
				Example:     "en",
			},
		},
		Example: `{
  "game_code": "bgaming_lucky_streak",
  "country": "US",
  "platform": "desktop",
  "exit_url": "https://example.com/lobby",
  "lang": "en"
}`,
		Notes: []string{
			`Fields "mode", "wl_code" and "token" are required by API, but CLI injects them automatically.`,
		},
	})
}
