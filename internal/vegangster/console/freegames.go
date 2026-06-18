package console

import "apikit/internal/kit/console"

func PrintFreegamesGrantBodySchema() {
	console.PrintSchema(console.Schema{
		Title: "Freegames: grant request body",
		Intro: "Request body must be a JSON object. Operator and Brand are taken from --config, so you don't need to pass them in --body.",
		Fields: []console.SchemaField{
			{
				Name:        "player_id",
				Required:    true,
				Type:        "string",
				Description: "Unique Player ID in the Operator’s system.",
				Example:     "64651509b8c355917ec34421",
			},
			{
				Name:        "country",
				Required:    true,
				Type:        "string",
				Description: "Player's country ISO 3166-1 code. Following enum contains all supported values.",
				Example:     "US",
			},
			{
				Name:        "ip",
				Required:    true,
				Type:        "string",
				Description: "Player's IP address",
				Example:     "145.22.35.62",
			},
			{
				Name:     "reference",
				Required: true,
				Type:     "string",
				Example:  "Ref_1_fsefgrh",
			},
			{
				Name:        "game_code",
				Required:    true,
				Type:        "string",
				Description: "Unique game identifier in Vegangster system in form of a string. Can be obtained from operator/v1/game/list endpoint.",
				Example:     "playson.100_joker_staxx",
			},
			{
				Name:        "rounds",
				Required:    true,
				Type:        "integer",
				Description: "Amount of rounds for a freegame session.",
				Example:     "10",
			},
			{
				Name:        "rounds_bet",
				Required:    true,
				Type:        "integer",
				Description: "Bet amount for each round of a freegame session.\n\tSubunits are used: if currency is 'USD' then rounds_bet=100 equals to a round bet of 1 USD",
				Example:     "100",
			},
			{
				Name:     "currency",
				Required: true,
				Type:     "string",
				Description: "ISO 4217 currency code. Following enum contains all currencies supported by our system.\n\t " +
					"Note that native game play support with these currencies varies per provider.\n\t " +
					"Please contact us to know which provider supports which currencies.",
				Example: "USD",
			},
			{
				Name:        "end_date",
				Required:    true,
				Type:        "date",
				Description: "ISO 8601 Extended datetime format (YYYY-MM-DDThh:mm:ss).",
				Example:     "2023-05-11T12:22:35",
			},
			{
				Name:        "offer_end_date",
				Required:    true,
				Type:        "date",
				Description: "ISO 8601 Extended datetime format (YYYY-MM-DDThh:mm:ss).",
				Example:     "2023-05-11T12:22:35",
			},
			{
				Name:        "start_date",
				Required:    false,
				Type:        "date",
				Description: "ISO 8601 Extended datetime format (YYYY-MM-DDThh:mm:ss).",
				Example:     "2023-05-11T12:22:35",
			},
		},
		Example: `{
  "player_id": "64651509b8c355917ec34421",
  "country": "US",
  "ip": "145.22.35.62",
  "reference": "Ref_1_fsefgrh",
  "game_code": "the_best_game",
  "rounds": 10,
  "rounds_bet": 100,
  "currency": "USD",
  "end_date": "2023-05-11T12:22:35",
  "offer_end_date": "2023-05-11T12:22:35",
  "start_date": "2023-05-11T12:22:35"
}`,
		Notes: []string{
			`Fields "operator_id" and "brand_id" are required by API, but CLI injects them from config.`,
		},
	})
}

func PrintFreegamesCancelBodySchema() {
	console.PrintSchema(console.Schema{
		Title: "Freegames: cancel request body",
		Intro: "Request body must be a JSON object. Operator and Brand are taken from --config, so you don't need to pass them in --body.",
		Fields: []console.SchemaField{
			{
				Name:        "id",
				Required:    true,
				Type:        "string",
				Description: "Id of the freegame offer generated on Vegangster side.",
				Example:     "64e4c3900e812e194e6e3767",
			},
		},
		Example: `{
  "id": "64e4c3900e812e194e6e3767"
}`,
		Notes: []string{
			`Fields "operator_id" and "brand_id" are required by API, but CLI injects them from config.`,
		},
	})
}

func PrintFreegamesBetAmountsBodySchema() {
	console.PrintSchema(console.Schema{
		Title: "Freegames: bet-amounts request body",
		Intro: "Request body must be a JSON object. Operator and Brand are taken from --config, so you don't need to pass them in --body.",
		Fields: []console.SchemaField{
			{
				Name:        "game_code",
				Required:    true,
				Type:        "string",
				Description: "Unique game identifier in Vegangster system. Can be obtained from /operator/v1/game/list.",
				Example:     "playson.100_joker_staxx",
			},
			{
				Name:        "currency",
				Required:    false,
				Type:        "string (ISO-4217)",
				Description: "ISO 4217 currency code. Supported values depend on provider/currency support.",
				Example:     "EUR",
			},
			{
				Name:        "country",
				Required:    false,
				Type:        "string (ISO-3166-1 alpha-2)",
				Description: "Player country code (2-letter).",
				Example:     "US",
			},
		},
		Example: `{
  "game_code": "playson.100_joker_staxx",
  "currency": "EUR",
  "country": "US"
}`,
		Notes: []string{
			`Fields "operator_id" and "brand_id" are required by API, but CLI injects them from config.`,
		},
	})
}

func PrintFreegamesBetAmountListBodySchema() {
	console.PrintSchema(console.Schema{
		Title: "Freegames: bet-amount list request body",
		Intro: "Request body must be a JSON object. Operator and Brand are taken from --config, so you don't need to pass them in --body.",
		Fields: []console.SchemaField{
			{
				Name:        "game_codes",
				Required:    true,
				Type:        "array",
				Description: "List of unique game identifier in Vegangster system. Can be obtained from /operator/v1/game/list.",
				Example:     "[\"playson.100_joker_staxx\"]",
			},
			{
				Name:        "currencies",
				Required:    false,
				Type:        "string (ISO-4217)",
				Description: "List of ISO 4217 currency code. Supported values depend on provider/currency support.",
				Example:     "[\"EUR\"]",
			},
		},
		Example: `{
  "game_codes": [
	"playson.100_joker_staxx"
  ],
  "currencies": ["EUR"]
}`,
		Notes: []string{
			`Fields "operator_id" and "brand_id" are required by API, but CLI injects them from config.`,
		},
	})
}
