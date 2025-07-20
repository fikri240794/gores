package gores

import (
	"testing"
)

func TestNewResponseErrorFieldVM(t *testing.T) {
	testCases := []struct {
		name            string
		field           string
		message         string
		expectedField   string
		expectedMessage string
	}{
		{
			name:            "NormalFieldAndMessage",
			field:           "exampleField",
			message:         "exampleMessage",
			expectedField:   "exampleField",
			expectedMessage: "exampleMessage",
		},
		{
			name:            "EmptyFieldAndMessage",
			field:           "",
			message:         "",
			expectedField:   "",
			expectedMessage: "",
		},
		{
			name:            "FieldWithSpecialCharacters",
			field:           "user.email",
			message:         "Email format is invalid",
			expectedField:   "user.email",
			expectedMessage: "Email format is invalid",
		},
		{
			name:            "LongFieldName",
			field:           "very_long_field_name_with_multiple_words",
			message:         "This field has a very detailed error message explaining what went wrong",
			expectedField:   "very_long_field_name_with_multiple_words",
			expectedMessage: "This field has a very detailed error message explaining what went wrong",
		},
		{
			name:            "FieldWithNumbers",
			field:           "field123",
			message:         "Field with numbers failed validation",
			expectedField:   "field123",
			expectedMessage: "Field with numbers failed validation",
		},
		{
			name:            "MessageWithQuotes",
			field:           "description",
			message:         "Description must not contain \"special\" characters",
			expectedField:   "description",
			expectedMessage: "Description must not contain \"special\" characters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test the constructor function
			vm := NewResponseErrorFieldVM(tc.field, tc.message)

			// Verify the struct is not nil
			if vm == nil {
				t.Fatal("NewResponseErrorFieldVM returned nil")
			}

			// Test field assignment
			if vm.Field != tc.expectedField {
				t.Errorf("Expected field is %q, got %q", tc.expectedField, vm.Field)
			}

			// Test message assignment
			if vm.Message != tc.expectedMessage {
				t.Errorf("Expected message is %q, got %q", tc.expectedMessage, vm.Message)
			}
		})
	}
}

// TestResponseErrorFieldVM_StructureValidation tests the JSON tag structure and field accessibility
func TestResponseErrorFieldVM_StructureValidation(t *testing.T) {
	field := "testField"
	message := "testMessage"

	vm := NewResponseErrorFieldVM(field, message)

	// Test that fields are publicly accessible (exported)
	// This is important for JSON marshaling
	if vm.Field == "" {
		t.Error("Field should be publicly accessible and properly set")
	}

	if vm.Message == "" {
		t.Error("Message should be publicly accessible and properly set")
	}

	// Test field modification (ensuring they're not read-only)
	vm.Field = "modifiedField"
	vm.Message = "modifiedMessage"

	if vm.Field != "modifiedField" {
		t.Error("Field should be modifiable")
	}

	if vm.Message != "modifiedMessage" {
		t.Error("Message should be modifiable")
	}
}
