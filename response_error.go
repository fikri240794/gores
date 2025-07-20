package gores

import "github.com/fikri240794/gocerr"

// ResponseErrorVM represents error information in standardized API responses.
// It contains a human-readable error message and optional field-specific errors.
// This structure provides detailed error context for client applications.
type ResponseErrorVM struct {
	Message     string                  `json:"message"`                // Primary error message
	ErrorFields []*ResponseErrorFieldVM `json:"error_fields,omitempty"` // Field-specific validation errors
}

// NewResponseErrorVM creates a new instance of ResponseErrorVM with initialized empty fields.
// This constructor ensures proper memory allocation and prevents nil pointer issues.
func NewResponseErrorVM() *ResponseErrorVM {
	return &ResponseErrorVM{
		ErrorFields: make([]*ResponseErrorFieldVM, 0), // Pre-allocate slice to avoid nil slice issues
	}
}

// SetMessage sets the primary error message for the response.
// This method follows the fluent API pattern for method chaining.
// The message should be user-friendly and provide clear information about the error.
func (vm *ResponseErrorVM) SetMessage(message string) *ResponseErrorVM {
	vm.Message = message
	return vm
}

// AddErrorFields appends one or more field-specific errors to the error response.
// This method is optimized for performance by pre-calculating required capacity
// to minimize slice reallocations when adding multiple fields.
// It handles both single field additions and bulk operations efficiently.
func (vm *ResponseErrorVM) AddErrorFields(errorFields ...*ResponseErrorFieldVM) *ResponseErrorVM {
	// Pre-allocate capacity if current slice would need expansion
	currentLen := len(vm.ErrorFields)
	newLen := currentLen + len(errorFields)

	// Optimize memory allocation by ensuring sufficient capacity
	if cap(vm.ErrorFields) < newLen {
		// Create new slice with appropriate capacity to avoid multiple reallocations
		newSlice := make([]*ResponseErrorFieldVM, currentLen, newLen)
		copy(newSlice, vm.ErrorFields)
		vm.ErrorFields = newSlice
	}

	// Append all error fields efficiently
	vm.ErrorFields = append(vm.ErrorFields, errorFields...)
	return vm
}

// mapFromCustomError efficiently extracts error information from gocerr.Error types.
// This method uses gocerr helper functions for safer and more maintainable error field extraction.
// It optimizes performance by leveraging gocerr's optimized field access methods.
func (vm *ResponseErrorVM) mapFromCustomError(customErr gocerr.Error) *ResponseErrorVM {
	vm.Message = customErr.Message

	// Use gocerr.GetErrorFields for safe and efficient error field extraction
	if errorFields := gocerr.GetErrorFields(customErr); len(errorFields) > 0 {
		// Pre-allocate slice with exact capacity to avoid reallocations
		responseFields := make([]*ResponseErrorFieldVM, 0, len(errorFields))

		// Convert each gocerr.ErrorField to ResponseErrorFieldVM efficiently
		for i := range errorFields {
			responseField := NewResponseErrorFieldVM(
				errorFields[i].Field,
				errorFields[i].Message,
			)
			responseFields = append(responseFields, responseField)
		}

		vm.ErrorFields = responseFields
	}

	return vm
}

// ParseError automatically extracts error information from any Go error type.
// It leverages gocerr.Parse for robust error type detection and processing.
// This method provides the main error parsing logic used throughout the library.
// For nil errors, it returns early to avoid unnecessary processing.
func (vm *ResponseErrorVM) ParseError(err error) *ResponseErrorVM {
	// Early return for nil errors to optimize performance
	if err == nil {
		return vm
	}

	// Use gocerr.Parse for robust custom error detection and extraction
	if customError, ok := gocerr.Parse(err); ok {
		return vm.mapFromCustomError(customError)
	}

	// Handle standard Go errors
	return vm.SetMessage(err.Error())
}
