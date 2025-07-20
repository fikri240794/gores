// Package gores provides standardized HTTP response structures for REST API development.
// It offers type-safe response models with consistent error handling and JSON serialization.
package gores

import (
	"net/http"

	"github.com/fikri240794/gocerr"
)

// ResponseVM represents a standardized HTTP response structure with generic data support.
// It provides a consistent format for API responses including status codes, error information, and data payload.
// The generic type T allows for type-safe data handling while maintaining flexibility.
type ResponseVM[T comparable] struct {
	Code  int              `json:"code"`            // HTTP status code
	Error *ResponseErrorVM `json:"error,omitempty"` // Error details if any
	Data  T                `json:"data,omitempty"`  // Response payload data
}

// NewResponseVM creates a new instance of ResponseVM with zero values.
// This is the preferred way to initialize a response struct to ensure proper memory allocation.
func NewResponseVM[T comparable]() *ResponseVM[T] {
	return &ResponseVM[T]{}
}

// SetCode sets the HTTP status code for the response.
// This method uses method chaining pattern for fluent API design.
// The code should be a valid HTTP status code (e.g., 200, 400, 500).
func (vm *ResponseVM[T]) SetCode(code int) *ResponseVM[T] {
	vm.Code = code
	return vm
}

// SetData sets the data payload for the response.
// This method accepts any type that satisfies the comparable constraint.
// The data will be serialized as JSON in the response body.
func (vm *ResponseVM[T]) SetData(data T) *ResponseVM[T] {
	vm.Data = data
	return vm
}

// SetError sets the error information for the response.
// This method allows manual error assignment when you have a pre-constructed ResponseErrorVM.
// For automatic error parsing from Go errors, use SetErrorFromError instead.
func (vm *ResponseVM[T]) SetError(err *ResponseErrorVM) *ResponseVM[T] {
	vm.Error = err
	return vm
}

// SetErrorFromError automatically processes a Go error and sets appropriate response fields.
// It leverages gocerr helper functions for robust error handling and code extraction.
// For nil errors, the method returns early without modifications for performance.
// For gocerr.Error types, it extracts the custom HTTP status code and error fields.
// For standard errors, it defaults to HTTP 500 Internal Server Error.
func (vm *ResponseVM[T]) SetErrorFromError(err error) *ResponseVM[T] {
	// Early return for nil errors to avoid unnecessary processing
	if err == nil {
		return vm
	}

	// Default to internal server error for safety
	vm.Code = http.StatusInternalServerError

	// Use gocerr.GetErrorCode for safe error code extraction
	if errorCode := gocerr.GetErrorCode(err); errorCode != 0 {
		// Override with custom error code if available
		vm.Code = errorCode
	}

	// Parse error details efficiently using enhanced ParseError method
	vm.Error = NewResponseErrorVM().ParseError(err)

	return vm
}
