package dto

// freeSpinsCurrencies (CLI: freegames bet-amounts) — GraphQL query variables.

type FreegamesBetAmountsVariables struct {
	JurisdictionCode string   `json:"jurisdictionCode" validate:"required"`
	GameId           int      `json:"gameId" validate:"required"`
	Currencies       []string `json:"currencies,omitempty"`
}
