package models

// Response is the response struct, that will be sent back to the user
type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
