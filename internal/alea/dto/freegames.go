package dto

// freeSpinsCurrencies (CLI: freegames bet-amounts) — GraphQL query variables.

type FreegamesBetAmountsVariables struct {
	JurisdictionCode string   `json:"jurisdictionCode" validate:"required"`
	GameId           int      `json:"gameId" validate:"required"`
	Currencies       []string `json:"currencies,omitempty"`
}

// freegames grant (CLI: freegames grant) — POST /api/v1/free-spins body (signed).

type FreegamesGrantRequest struct {
	CasinoPlayerId   string      `json:"casinoPlayerId" validate:"required"`
	CasinoBonusId    string      `json:"casinoBonusId" validate:"required"`
	CasinoCampaignId string      `json:"casinoCampaignId,omitempty"`
	Country          string      `json:"country" validate:"required"`
	Currency         string      `json:"currency" validate:"required"`
	SoftwareId       int         `json:"softwareId" validate:"required"`
	Level            int         `json:"level" validate:"required"`
	Amount           int         `json:"amount" validate:"required"`
	Games            []GrantGame `json:"games" validate:"required,min=1,dive"`
	StartAt          string      `json:"startAt,omitempty"`
	ExpireAt         string      `json:"expireAt,omitempty"`
}

type GrantGame struct {
	Id     int    `json:"id" validate:"required"`
	Device string `json:"device,omitempty"`
}

// freegames cancel (CLI: freegames cancel) — DELETE /api/v1/free-spins/{casinoBonusId} (signed).

type FreegamesCancelRequest struct {
	CasinoBonusId string `json:"casinoBonusId" validate:"required"`
}
