package console

import "apikit/internal/kit/console"

func PrintGameLaunchQuerySchema() {
	console.PrintSchema(console.Schema{
		Title: "Game: launch query params",
		Intro: "Creates a real-money game session URL for a player. Endpoint: GET /{KEY}/User/AuthHTML/",
		Fields: []console.SchemaField{
			{
				Name:        "Login",
				Required:    true,
				Type:        "string",
				Description: "Player wallet ID (UUID)",
				Example:     "550e8400-e29b-41d4-a716-446655440000",
			},
			{
				Name:        "Password",
				Required:    true,
				Type:        "string",
				Description: "Player ID (UUID)",
				Example:     "550e8400-e29b-41d4-a716-446655440001",
			},
			{
				Name:        "Currency",
				Required:    true,
				Type:        "string",
				Description: "Player's currency in ISO 4217 format",
				Example:     "USD",
			},
			{
				Name:        "ExtParam",
				Required:    true,
				Type:        "string",
				Description: "Game session ID (UUID)",
				Example:     "550e8400-e29b-41d4-a716-446655440002",
			},
			{
				Name:        "System",
				Required:    true,
				Type:        "string",
				Description: "External provider ID",
				Example:     "789",
			},
			{
				Name:        "Page",
				Required:    true,
				Type:        "string",
				Description: "External game ID",
				Example:     "123",
			},
			{
				Name:        "UserAutoCreate",
				Required:    false,
				Type:        "integer",
				Description: "Whether to auto-create the player if not exists.\n\tPossible values:\n\t\"1\" - auto-create\n\t\"0\" - do not auto-create",
				Example:     "1",
			},
			{
				Name:        "IsMobile",
				Required:    false,
				Type:        "integer",
				Description: "Player's device type.\n\tPossible values:\n\t\"1\" - mobile\n\t\"0\" - desktop",
				Example:     "0",
			},
			{
				Name:        "UserIP",
				Required:    false,
				Type:        "string",
				Description: "Player's real IP address. Required for licensing.\n\tIf not provided, server IP is used as fallback.",
				Example:     "192.168.1.1",
			},
		},
		Example: `{
  "Login": "550e8400-e29b-41d4-a716-446655440000",
  "Password": "550e8400-e29b-41d4-a716-446655440001",
  "Currency": "USD",
  "ExtParam": "550e8400-e29b-41d4-a716-446655440002",
  "System": "789",
  "Page": "123",
  "UserAutoCreate": 1,
  "IsMobile": 0,
  "UserIP": "192.168.1.1"
}`,
		Notes: []string{
			`Fields "TID" and "Hash" are required by API, but CLI generates and injects them automatically.`,
		},
	})
}

func PrintGameDemoLaunchQuerySchema() {
	console.PrintSchema(console.Schema{
		Title: "Game: demo launch query params",
		Intro: "Creates a demo (free-play) game session URL. Endpoint: GET /{KEY}/User/AuthHTML/",
		Fields: []console.SchemaField{
			{
				Name:        "System",
				Required:    true,
				Type:        "string",
				Description: "External provider ID",
				Example:     "634",
			},
			{
				Name:        "Page",
				Required:    true,
				Type:        "string",
				Description: "External game ID",
				Example:     "235",
			},
			{
				Name:        "UserAutoCreate",
				Required:    false,
				Type:        "integer",
				Description: "Whether to auto-create the player if not exists.\n\tPossible values:\n\t\"1\" - auto-create\n\t\"0\" - do not auto-create",
				Example:     "1",
			},
			{
				Name:        "IsMobile",
				Required:    false,
				Type:        "integer",
				Description: "Player's device type.\n\tPossible values:\n\t\"1\" - mobile\n\t\"0\" - desktop",
				Example:     "0",
			},
			{
				Name:        "UserIP",
				Required:    false,
				Type:        "string",
				Description: "Player's real IP address. Required for licensing.\n\tIf not provided, server IP is used as fallback.",
				Example:     "192.168.1.1",
			},
		},
		Example: `{
  "System": "634",
  "Page": "235",
  "UserAutoCreate": 1,
  "IsMobile": 0,
  "UserIP": "192.168.1.1"
}`,
		Notes: []string{
			`Fields "Login", "Password" and "Demo" are required by API, but CLI injects them automatically.`,
			`Fields "TID" and "Hash" are required by API, but CLI generates and injects them automatically.`,
		},
	})
}
