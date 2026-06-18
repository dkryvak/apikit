package dto

// /operator/v1/freegames/grant

type FreegamesGrantRequest struct {
	OperatorID   string `json:"operator_id" validate:"required"`
	BrandID      string `json:"brand_id" validate:"required"`
	PlayerID     string `json:"player_id" validate:"required"`
	Country      string `json:"country" validate:"required"`
	Ip           string `json:"ip"`
	Reference    string `json:"reference" validate:"required"`
	GameCode     string `json:"game_code" validate:"required"`
	Rounds       int    `json:"rounds" validate:"required"`
	RoundsBet    int    `json:"rounds_bet" validate:"required"`
	Currency     string `json:"currency" validate:"required"`
	EndDate      string `json:"end_date" validate:"required"`
	OfferEndDate string `json:"offer_end_date" validate:"required"`
	StartDate    string `json:"start_date"`
}

// /operator/v1/freegames/cancel

type FreegamesCancelRequest struct {
	OperatorID string `json:"operator_id" validate:"required"`
	BrandID    string `json:"brand_id" validate:"required"`
	ID         string `json:"id" validate:"required"`
}

// /operator/v1/freegames/bet-amounts

type FreegamesBetAmountsRequest struct {
	OperatorID string `json:"operator_id" validate:"required"`
	BrandID    string `json:"brand_id" validate:"required"`
	GameCode   string `json:"game_code,omitempty"`
	Currency   string `json:"currency,omitempty"`
	Country    string `json:"country,omitempty"`
}

// /operator/v1/freegames/bet-amount/list

type FreegamesBetAmountListRequest struct {
	OperatorID string `json:"operator_id" validate:"required"`
	BrandID    string `json:"brand_id" validate:"required"`
	GameCodes  string `json:"game_codes,omitempty"`
	Currencies string `json:"currencies,omitempty"`
}
