package gores

type ResponseErrorFieldVM struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewResponseErrorFieldVM(field string, message string) *ResponseErrorFieldVM {
	return &ResponseErrorFieldVM{
		Field:   field,
		Message: message,
	}
}
