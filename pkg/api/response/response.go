package response

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func OK(data any) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

func Error(msg string) *Response {
	return &Response{
		Success: false,
		Message: msg,
	}
}
