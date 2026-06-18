package dto

// /{KEY}/Freerounds/Add/

type FreegamesGrantQuery struct {
	Login    string `json:"Login,omitempty" validate:"required"`
	Operator string `json:"Operator,omitempty" validate:"required"`
	Games    string `json:"Games,omitempty" validate:"required"`
	Count    int16  `json:"Count,omitempty" validate:"required"`
	Expire   string `json:"Expire,omitempty" validate:"required,datetime=2006-01-02 15:04:05"`
	Country  string `json:"Country,omitempty"`
	TID      string `json:"TID,omitempty" validate:"required"`
	Hash     string `json:"Hash,omitempty"`
}

// /{KEY}/Freerounds/Remove/

type FreegamesCancelQuery struct {
	Login    string `json:"Login,omitempty" validate:"required"`
	Operator string `json:"Operator,omitempty" validate:"required"`
	ExtID    string `json:"ExtID,omitempty" validate:"required"`
	TID      string `json:"TID,omitempty"`
	Hash     string `json:"Hash,omitempty"`
}
