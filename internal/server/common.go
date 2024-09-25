package server

type Response struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
