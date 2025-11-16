package dto

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorDTO struct {
	Error Error `json:"error"`
}

func NewErrorResponse(code, message string) ErrorDTO {
	return ErrorDTO{
		Error: Error{
			Code:    code,
			Message: message,
		},
	}
}
