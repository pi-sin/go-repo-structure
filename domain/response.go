package domain

// Response represent the response struct
type Response struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
}
