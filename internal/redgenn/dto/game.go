package dto

// /admservice

type GameListRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Command  string `json:"cm" validate:"required"`

	Disabled       *int8  `json:"disabled,omitempty"`
	Producer       string `json:"producer,omitempty"`
	Provider       string `json:"provider,omitempty"`
	PromoFreespins *int8  `json:"promo_freespins,omitempty"`
	GameId         string `json:"game_id,omitempty"`
	Title          string `json:"title,omitempty"`
	Limit          *int   `json:"limit,omitempty"`
	Offset         *int   `json:"offset,omitempty"`
}

type GameLaunchRequest struct {
	Mode     string `json:"mode" validate:"required"`
	WlCode   string `json:"wl_code" validate:"required"`
	Token    string `json:"token" validate:"required"`
	GameCode string `json:"game_code" validate:"required"`

	Country    string `json:"country,omitempty"`
	Platform   string `json:"platform,omitempty"`
	ExitUrl    string `json:"exit_url,omitempty"`
	CashierUrl string `json:"cashier_url,omitempty"`
	Language   string `json:"lang,omitempty"`
}
