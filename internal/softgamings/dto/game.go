package dto

// /User/AuthHTML/

type GameLaunchQuery struct {
	Login          string `json:"Login,omitempty" validate:"required"`
	Password       string `json:"Password,omitempty" validate:"required"`
	Currency       string `json:"Currency,omitempty" validate:"required"`
	ExtParam       string `json:"ExtParam,omitempty" validate:"required"`
	UserAutoCreate int8   `json:"UserAutoCreate,omitempty"`
	System         string `json:"System,omitempty" validate:"required"`
	Page           string `json:"Page,omitempty" validate:"required"`
	UserIP         string `json:"UserIP,omitempty"`
	IsMobile       int8   `json:"IsMobile,omitempty"`
	TID            string `json:"TID,omitempty"`
	Hash           string `json:"Hash,omitempty"`
}

type GameDemoLaunchQuery struct {
	Login          string `json:"Login,omitempty" validate:"required"`
	Password       string `json:"Password,omitempty" validate:"required"`
	Demo           int8   `json:"Demo,omitempty" validate:"required"`
	UserAutoCreate int8   `json:"UserAutoCreate,omitempty"`
	System         string `json:"System,omitempty" validate:"required"`
	Page           string `json:"Page,omitempty" validate:"required"`
	UserIP         string `json:"UserIP,omitempty"`
	IsMobile       int8   `json:"IsMobile,omitempty"`
	TID            string `json:"TID,omitempty"`
	Hash           string `json:"Hash,omitempty"`
}
