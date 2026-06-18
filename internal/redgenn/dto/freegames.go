package dto

// /admservice

type FreegamesBetLevelsRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Command  string `json:"cm" validate:"required"`
	WlCode   string `json:"wlcode" validate:"required"`

	Currency string `json:"currency"`
	Game     string `json:"game,omitempty"`
	Producer string `json:"producer,omitempty"`
}
