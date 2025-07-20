// Package main demonstrates comprehensive usage of the gores library
// for standardized HTTP response handling in REST API applications.
// This example covers scenarios from simple success responses to complex error handling.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fikri240794/gocerr"
	"github.com/fikri240794/gores"
)

// User represents a sample data structure for demonstration purposes
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	IsActive bool   `json:"is_active"`
}

// EmptyResponse represents responses with no data payload
type EmptyResponse struct{}

// SimpleUserList represents a simple collection of users for demonstration
type SimpleUserList struct {
	Users []*User `json:"users"`
	Count int     `json:"count"`
}

func main() {
	fmt.Println("=== GORES Library Demonstration ===")
	fmt.Println("This program demonstrates all features of the gores library")
	fmt.Println("including success responses, error handling, and method chaining.")
	fmt.Println()

	// Demonstration 1: Simple Success Response
	fmt.Println("1. Simple Success Response:")
	demonstrateSimpleSuccess()

	// Demonstration 2: Success Response with Complex Data
	fmt.Println("\n2. Success Response with Complex Data:")
	demonstrateComplexData()

	// Demonstration 3: Standard Error Response
	fmt.Println("\n3. Standard Error Response:")
	demonstrateStandardError()

	// Demonstration 4: Custom Error Response
	fmt.Println("\n4. Custom Error Response:")
	demonstrateCustomError()

	// Demonstration 5: Validation Error with Multiple Fields
	fmt.Println("\n5. Validation Error with Multiple Fields:")
	demonstrateValidationError()

	// Demonstration 5.1: Advanced gocerr Features
	fmt.Println("\n5.1. Advanced gocerr Features:")
	demonstrateAdvancedGocerr()

	// Demonstration 6: Method Chaining Patterns
	fmt.Println("\n6. Method Chaining Patterns:")
	demonstrateMethodChaining()

	// Demonstration 7: Different HTTP Status Codes
	fmt.Println("\n7. Different HTTP Status Codes:")
	demonstrateDifferentStatusCodes()

	// Demonstration 8: Empty/No Content Responses
	fmt.Println("\n8. Empty/No Content Responses:")
	demonstrateEmptyResponses()

	// Demonstration 9: Error from Nil (Edge Case)
	fmt.Println("\n9. Error from Nil (Edge Case):")
	demonstrateNilErrorHandling()

	fmt.Println("\n=== GoRes Library Demonstration Complete ===")
	fmt.Println("All core features have been demonstrated successfully!")
}

// demonstrateSimpleSuccess shows basic success response creation
func demonstrateSimpleSuccess() {
	user := &User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Age:      30,
		IsActive: true,
	}

	// Create simple success response
	response := gores.NewResponseVM[*User]().
		SetCode(http.StatusOK).
		SetData(user)

	printResponse("Simple User Response", response)
}

// demonstrateComplexData shows response with complex nested data structures
func demonstrateComplexData() {
	users := []*User{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30, IsActive: true},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25, IsActive: true},
		{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", Age: 35, IsActive: false},
	}

	userList := &SimpleUserList{
		Users: users,
		Count: len(users),
	}

	// Create response with complex data structure
	response := gores.NewResponseVM[*SimpleUserList]().
		SetCode(http.StatusOK).
		SetData(userList)

	printResponse("User List Response", response)
}

// demonstrateStandardError shows handling of standard Go errors
func demonstrateStandardError() {
	// Simulate a standard Go error
	err := errors.New("database connection failed")

	// Create error response from standard error
	response := gores.NewResponseVM[*User]().
		SetErrorFromError(err)

	printResponse("Standard Error Response", response)
}

// demonstrateCustomError shows handling of custom gocerr errors
func demonstrateCustomError() {
	// Create custom error with specific HTTP status
	customErr := gocerr.New(
		http.StatusNotFound,
		"User not found",
	)

	// Create error response from custom error
	response := gores.NewResponseVM[*User]().
		SetErrorFromError(customErr)

	printResponse("Custom Error Response", response)
}

// demonstrateValidationError shows validation errors with field-specific messages
func demonstrateValidationError() {
	// Create validation error with multiple field errors
	validationErr := gocerr.New(
		http.StatusUnprocessableEntity,
		"Validation failed",
		gocerr.NewErrorField("name", "name is required and must be at least 2 characters"),
		gocerr.NewErrorField("email", "email format is invalid"),
		gocerr.NewErrorField("age", "age must be between 18 and 120"),
		gocerr.NewErrorField("password", "password must be at least 8 characters with uppercase, lowercase, and numbers"),
	)

	// Create validation error response
	response := gores.NewResponseVM[*User]().
		SetErrorFromError(validationErr)

	printResponse("Validation Error Response", response)
}

// demonstrateAdvancedGocerr shows advanced features of gocerr
func demonstrateAdvancedGocerr() {
	// Create a comprehensive validation error
	validationErr := gocerr.New(
		http.StatusUnprocessableEntity,
		"User registration failed",
		gocerr.NewErrorField("username", "username must be unique"),
		gocerr.NewErrorField("email", "email format is invalid"),
		gocerr.NewErrorField("phone", "phone number is required"),
	)

	// Demonstrate gocerr utility functions
	fmt.Printf("Error Code: %d\n", gocerr.GetErrorCode(validationErr))
	fmt.Printf("Has Error Fields: %t\n", gocerr.HasErrorFields(validationErr))
	fmt.Printf("Error Field Count: %d\n", gocerr.ErrorFieldCount(validationErr))
	fmt.Printf("Has 'email' field: %t\n", gocerr.HasErrorField(validationErr, "email"))
	fmt.Printf("Email error message: %s\n", gocerr.GetErrorFieldMessage(validationErr, "email"))
	fmt.Printf("Is code 422: %t\n", gocerr.IsErrorCodeEqual(validationErr, http.StatusUnprocessableEntity))

	// Demonstrate parsing
	if parsed, ok := gocerr.Parse(validationErr); ok {
		fmt.Printf("Parsed error successfully: %s\n", parsed.String())
	}

	// Create gores response using the error
	response := gores.NewResponseVM[*User]().
		SetErrorFromError(validationErr)

	printResponse("Advanced gocerr Features", response)
}

// demonstrateMethodChaining shows various method chaining patterns
func demonstrateMethodChaining() {
	// Example 1: Success response with chaining
	successResponse := gores.NewResponseVM[*User]().
		SetCode(http.StatusCreated).
		SetData(&User{
			ID:       100,
			Name:     "New User",
			Email:    "newuser@example.com",
			Age:      28,
			IsActive: true,
		})

	printResponse("Chained Success Response", successResponse)

	// Example 2: Error response with manual error construction
	errorResponse := gores.NewResponseVM[*User]().
		SetCode(http.StatusBadRequest).
		SetError(
			gores.NewResponseErrorVM().
				SetMessage("Manual error construction").
				AddErrorFields(
					gores.NewResponseErrorFieldVM("custom_field", "This is a manually constructed error field"),
				),
		)

	printResponse("Chained Error Response", errorResponse)
}

// demonstrateDifferentStatusCodes shows responses with various HTTP status codes
func demonstrateDifferentStatusCodes() {
	statusCodes := []struct {
		code        int
		description string
		hasData     bool
	}{
		{http.StatusOK, "OK - Success", true},
		{http.StatusCreated, "Created - Resource created successfully", true},
		{http.StatusAccepted, "Accepted - Request accepted for processing", false},
		{http.StatusNoContent, "No Content - Successful with no data", false},
		{http.StatusBadRequest, "Bad Request - Client error", false},
		{http.StatusUnauthorized, "Unauthorized - Authentication required", false},
		{http.StatusForbidden, "Forbidden - Access denied", false},
		{http.StatusNotFound, "Not Found - Resource not found", false},
		{http.StatusConflict, "Conflict - Resource conflict", false},
		{http.StatusInternalServerError, "Internal Server Error - Server error", false},
	}

	for _, sc := range statusCodes {
		var response *gores.ResponseVM[*User]

		if sc.hasData {
			// Create response with data
			response = gores.NewResponseVM[*User]().
				SetCode(sc.code).
				SetData(&User{
					ID:    1,
					Name:  "Sample User",
					Email: "user@example.com",
					Age:   25,
				})
		} else if sc.code >= 400 {
			// Create error response
			response = gores.NewResponseVM[*User]().
				SetCode(sc.code).
				SetError(
					gores.NewResponseErrorVM().
						SetMessage(sc.description),
				)
		} else {
			// Create success response without data
			response = gores.NewResponseVM[*User]().
				SetCode(sc.code)
		}

		printResponse(fmt.Sprintf("%d - %s", sc.code, sc.description), response)
	}
}

// demonstrateEmptyResponses shows responses with no data payload
func demonstrateEmptyResponses() {
	// No content response (successful operation with no data)
	noContentResponse := gores.NewResponseVM[*EmptyResponse]().
		SetCode(http.StatusNoContent)

	printResponse("No Content Response", noContentResponse)

	// Accepted response (operation accepted but not completed)
	acceptedResponse := gores.NewResponseVM[*EmptyResponse]().
		SetCode(http.StatusAccepted)

	printResponse("Accepted Response", acceptedResponse)
}

// demonstrateNilErrorHandling shows edge case handling for nil errors
func demonstrateNilErrorHandling() {
	// Test nil error handling
	response := gores.NewResponseVM[*User]().
		SetCode(http.StatusOK).
		SetErrorFromError(nil). // This should not modify the response
		SetData(&User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
			Age:   30,
		})

	printResponse("Nil Error Handling (should remain success)", response)
}

// printResponse is a helper function to pretty print responses for demonstration
func printResponse(title string, response interface{}) {
	fmt.Printf("--- %s ---\n", title)

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))
}
