package dto

// /admservice

type FreegamesBetAmountsRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Command  string `json:"cm" validate:"required"`
	WlCode   string `json:"wlcode" validate:"required"`

	Currency string `json:"currency"`
	Game     string `json:"game,omitempty"`
	Producer string `json:"producer,omitempty"`
}

// freegames grant (CLI: freegames grant) — create_offer + add_offer params.
// login/password/wlcode are injected from the config.

type FreegamesGrantRequest struct {
	Offer    string `json:"offer" validate:"required"`
	Game     string `json:"game" validate:"required"`
	Spins    int    `json:"spins" validate:"required"`
	Totalbet int    `json:"totalbet" validate:"required"` // per-spin total bet, in cents
	Currency string `json:"currency" validate:"required"`
	Wlid     string `json:"wlid" validate:"required"` // player id (add_offer)

	StartLive       string `json:"startLive,omitempty"`       // create_offer start; DD/MM/YYYY HH:MM UTC
	EndDateAddOffer string `json:"endDateAddOffer,omitempty"` // create_offer end; DD/MM/YYYY HH:MM UTC
	Lifetime        string `json:"lifetime,omitempty"`        // add_offer expiry; DD/MM/YYYY HH:MM UTC
	CmId            string `json:"cmId,omitempty"`            // add_offer idempotency key (uuid)
}

// freegames cancel (CLI: freegames cancel) — remove_offer by gift id.

type FreegamesCancelRequest struct {
	GiftId string `json:"giftId" validate:"required"`
}
