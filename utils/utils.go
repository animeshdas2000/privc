package utils

type Response struct {
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"error_message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type TokenRequestPayload struct {
	Id   string            `json:"id"`
	Data map[string]string `json:"data"`
}
