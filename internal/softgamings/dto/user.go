package dto

// /User/GetUserData/

type GetUserQuery struct {
	Login string `json:"Login,omitempty" validate:"required"`
	TID   string `json:"TID,omitempty"`
	Hash  string `json:"Hash,omitempty"`
}
