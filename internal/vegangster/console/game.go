package console

import "apikit/internal/kit/console"

func PrintGameLaunchBodySchema() {
	console.PrintSchema(console.Schema{
		Title: "Game: launch request body",
		Intro: "Creates a real-money game session URL for a player. Requires player authentication token.",
		Fields: []console.SchemaField{
			{
				Name:        "player_id",
				Required:    true,
				Type:        "string",
				Description: "Unique player identifier in the operator's system",
				Example:     "player123",
			},
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
				Name:        "platform",
				Required:    true,
				Type:        "string",
				Description: "Player's platform.\n\tPossible values:\n\t\"desktop\" - desktop browser\n\t\"mobile\" - mobile browser",
				Example:     "desktop",
			},
			{
				Name:        "currency",
				Required:    true,
				Type:        "string",
				Description: "Player's currency in ISO 4217 format",
				Example:     "USD",
			},
			{
				Name:        "lang",
				Required:    true,
				Type:        "string",
				Description: "Game interface language in ISO 639-1 format",
				Example:     "en",
			},
			{
				Name:        "country",
				Required:    true,
				Type:        "string",
				Description: "Player's country in ISO 3166-1 alpha-2 format",
				Example:     "US",
			},
			{
				Name:        "ip",
				Required:    true,
				Type:        "string",
				Description: "Player's IP address",
				Example:     "192.168.1.1",
			},
			{
				Name:        "lobby_url",
				Required:    false,
				Type:        "string",
				Description: "URL to redirect the player when they click the lobby/exit button",
				Example:     "https://example.com/lobby",
			},
			{
				Name:        "deposit_url",
				Required:    false,
				Type:        "string",
				Description: "URL to redirect the player when they click the deposit button",
				Example:     "https://example.com/deposit",
			},
			{
				Name:        "player_nick",
				Required:    false,
				Type:        "string",
				Description: "Player's display nickname shown in the game",
				Example:     "Lucky Player",
			},
		},
		Example: `{
  "player_id": "player123",
  "token": "eyJhbGciOiJSUzI1NiJ9...",
  "game_code": "bgaming_lucky_streak",
  "platform": "desktop",
  "currency": "USD",
  "lang": "en",
  "country": "US",
  "ip": "192.168.1.1",
  "lobby_url": "https://example.com/lobby",
  "deposit_url": "https://example.com/deposit",
  "player_nick": "Lucky Player"
}`,
		Notes: []string{
			`Fields "operator_id" and "brand_id" are required by API, but CLI injects them from config.`,
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
				Name:        "platform",
				Required:    true,
				Type:        "string",
				Description: "Player's platform.\n\tPossible values:\n\t\"desktop\" - desktop browser\n\t\"mobile\" - mobile browser",
				Example:     "desktop",
			},
			{
				Name:        "currency",
				Required:    true,
				Type:        "string",
				Description: "Currency used in demo mode in ISO 4217 format",
				Example:     "USD",
			},
			{
				Name:        "lang",
				Required:    false,
				Type:        "string",
				Description: "Game interface language in ISO 639-1 format",
				Example:     "en",
			},
			{
				Name:        "country",
				Required:    false,
				Type:        "string",
				Description: "Player's country in ISO 3166-1 alpha-2 format",
				Example:     "US",
			},
			{
				Name:        "ip",
				Required:    false,
				Type:        "string",
				Description: "Player's IP address",
				Example:     "192.168.1.1",
			},
			{
				Name:        "lobby_url",
				Required:    false,
				Type:        "string",
				Description: "URL to redirect the player when they click the lobby/exit button",
				Example:     "https://example.com/lobby",
			},
		},
		Example: `{
  "game_code": "bgaming_lucky_streak",
  "platform": "desktop",
  "currency": "USD",
  "lang": "en",
  "country": "US",
  "ip": "192.168.1.1",
  "lobby_url": "https://example.com/lobby"
}`,
		Notes: []string{
			`Fields "operator_id" and "brand_id" are required by API, but CLI injects them from config.`,
		},
	})
}
