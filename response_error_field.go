package gores

// ResponseErrorFieldVM represents a field-specific error in API validation responses.
// It provides detailed information about which field caused an error and why.
// This structure is commonly used for form validation and request parameter errors.
type ResponseErrorFieldVM struct {
	Field   string `json:"field"`   // The name of the field that caused the error
	Message string `json:"message"` // Human-readable error message for this field
}

// NewResponseErrorFieldVM creates a new field error with the specified field name and message.
// This constructor provides a convenient way to create field-specific errors with validation.
// Both field and message parameters are required for meaningful error reporting.
func NewResponseErrorFieldVM(field, message string) *ResponseErrorFieldVM {
	return &ResponseErrorFieldVM{
		Field:   field,
		Message: message,
	}
}
