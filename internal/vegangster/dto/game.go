package dto

// /operator/v1/game/list

type GameListRequest struct {
	OperatorID string `json:"operator_id" validate:"required"`
	BrandID    string `json:"brand_id" validate:"required"`
}

// /operator/v1/game/url

type GameLaunchRequest struct {
	OperatorID string `json:"operator_id"  validate:"required"`
	BrandID    string `json:"brand_id"  validate:"required"`
	PlayerId   string `json:"player_id"  validate:"required"`
	Token      string `json:"token"  validate:"required"`
	GameCode   string `json:"game_code"  validate:"required"`
	Platform   string `json:"platform"  validate:"required"`
	Currency   string `json:"currency"  validate:"required"`
	Language   string `json:"lang,omitempty"  validate:"required"`
	Country    string `json:"country,omitempty"  validate:"required"`
	IP         string `json:"ip,omitempty"  validate:"required"`
	LobbyUrl   string `json:"lobby_url,omitempty"`
	DepositUrl string `json:"deposit_url,omitempty"`
	PlayerNick string `json:"player_nick,omitempty"`
}

// /operator/v1/game/demo/url

type GameDemoLaunchRequest struct {
	OperatorID string `json:"operator_id" validate:"required"`
	BrandID    string `json:"brand_id" validate:"required"`
	GameCode   string `json:"game_code" validate:"required"`
	Platform   string `json:"platform" validate:"required"`
	Currency   string `json:"currency" validate:"required"`
	Language   string `json:"lang,omitempty"`
	Country    string `json:"country,omitempty"`
	IP         string `json:"ip,omitempty"`
	LobbyUrl   string `json:"lobby_url,omitempty"`
}
