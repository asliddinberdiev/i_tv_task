package common

type Response struct {
	Status  uint16      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseWithList struct {
	Status  uint16        `json:"status"`
	Message string        `json:"message"`
	Data    interface{} `json:"data"`
	Total   uint64        `json:"total"`
}

type ResponseID struct {
	Status  uint16 `json:"status"`
	Message string `json:"message"`
	ID      uint   `json:"id"`
}

type ResponseError struct {
	Status  uint16 `json:"status"`
	Message string `json:"message"`
}
