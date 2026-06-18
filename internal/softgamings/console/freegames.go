package console

import "apikit/internal/kit/console"

func PrintFreegamesGrantQuerySchema() {
	console.PrintSchema(console.Schema{
		Title: "Freegames: grant query params",
		Intro: "Grants a freegames bonus to a player. Endpoint: GET /{KEY}/Freerounds/Add/",
		Fields: []console.SchemaField{
			{
				Name:        "Login",
				Required:    true,
				Type:        "string",
				Description: "Player wallet ID (UUID)",
				Example:     "550e8400-e29b-41d4-a716-446655440000",
			},
			{
				Name:        "Operator",
				Required:    true,
				Type:        "string",
				Description: "Provider name",
				Example:     "Belatra",
			},
			{
				Name:        "Games",
				Required:    true,
				Type:        "string",
				Description: "External game ID",
				Example:     "lucky_streak",
			},
			{
				Name:        "Count",
				Required:    true,
				Type:        "integer",
				Description: "Number of free rounds to grant",
				Example:     "10",
			},
			{
				Name:        "Expire",
				Required:    true,
				Type:        "string",
				Description: "Bonus expiration date in format \"yyyy-MM-dd HH:mm:ss\"",
				Example:     "2026-12-31 23:59:59",
			},
			{
				Name:        "TID",
				Required:    true,
				Type:        "string",
				Description: "Wallet bonus balance ID (UUID). Used for hash calculation",
				Example:     "550e8400-e29b-41d4-a716-446655440001",
			},
			{
				Name:        "Country",
				Required:    false,
				Type:        "string",
				Description: "Player's country in ISO 3166-1 alpha-2 format",
				Example:     "US",
			},
		},
		Example: `{
  "Login": "550e8400-e29b-41d4-a716-446655440000",
  "Operator": "Belatra",
  "Games": "lucky_streak",
  "Count": 10,
  "Expire": "2026-12-31 23:59:59",
  "TID": "550e8400-e29b-41d4-a716-446655440001",
  "Country": "US"
}`,
		Notes: []string{
			`Field "Hash" is required by API, but CLI generates and injects it automatically.`,
		},
	})
}

func PrintFreegamesCancelQuerySchema() {
	console.PrintSchema(console.Schema{
		Title: "Freegames: cancel query params",
		Intro: "Cancels a freegames bonus for a player. Endpoint: GET /{KEY}/Freerounds/Remove/",
		Fields: []console.SchemaField{
			{
				Name:        "Login",
				Required:    true,
				Type:        "string",
				Description: "Player wallet ID (UUID)",
				Example:     "550e8400-e29b-41d4-a716-446655440000",
			},
			{
				Name:        "Operator",
				Required:    true,
				Type:        "string",
				Description: "Provider name",
				Example:     "Belatra",
			},
			{
				Name:        "ExtID",
				Required:    true,
				Type:        "string",
				Description: "Wallet bonus balance ID (UUID)",
				Example:     "550e8400-e29b-41d4-a716-446655440001",
			},
		},
		Example: `{
  "Login": "550e8400-e29b-41d4-a716-446655440000",
  "Operator": "Belatra",
  "ExtID": "550e8400-e29b-41d4-a716-446655440001"
}`,
		Notes: []string{
			`Fields "TID" and "Hash" are required by API, but CLI generates and injects them automatically.`,
		},
	})
}
