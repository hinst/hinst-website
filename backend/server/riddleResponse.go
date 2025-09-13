package server

type riddleResponse struct {
	Id      int `json:"id"`
	Product int `json:"product"`
	Steps   int `json:"steps"`
	Limit   int `json:"limit"`
}
