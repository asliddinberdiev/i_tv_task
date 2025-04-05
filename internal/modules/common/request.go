package common

type RequestID struct {
	ID uint `json:"id" validate:"required,numeric,gte=1"`
}

type RequestSearch struct {
	Search string `json:"search"`
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
}

