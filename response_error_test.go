package gores

import (
	"errors"
	"net/http"
	"testing"

	"github.com/fikri240794/gocerr"
)

// testResponseErrorVMEquality performs deep equality comparison between ResponseErrorVM instances.
// This helper function ensures comprehensive testing of all fields including error field slices.
func testResponseErrorVMEquality(t *testing.T, expected, actual *ResponseErrorVM) {
	t.Helper() // Mark as test helper for better error reporting

	// Test nil pointer combinations
	if expected == nil && actual != nil {
		t.Errorf("expected is nil, got not nil")
		return
	}

	if expected != nil && actual == nil {
		t.Errorf("expected not nil, got nil")
		return
	}

	if expected == nil && actual == nil {
		return
	}

	// Test message field
	if expected.Message != actual.Message {
		t.Errorf("expected message is %s, got %s", expected.Message, actual.Message)
	}

	// Test error fields slice length
	if len(expected.ErrorFields) != len(actual.ErrorFields) {
		t.Errorf("expected length of error fields is %d, got %d", len(expected.ErrorFields), len(actual.ErrorFields))
		return
	}

	// Test each error field in the slice
	for i := range expected.ErrorFields {
		expectedField := expected.ErrorFields[i]
		actualField := actual.ErrorFields[i]

		if expectedField != nil && actualField == nil {
			t.Errorf("expected error fields item at index %d is not nil, but got nil", i)
			continue
		}

		if expectedField == nil && actualField != nil {
			t.Errorf("expected error fields item at index %d is nil, but got not nil", i)
			continue
		}

		if expectedField == nil && actualField == nil {
			continue
		}

		if expectedField.Field != actualField.Field {
			t.Errorf("expected error fields item field at index %d is %s, got %s", i, expectedField.Field, actualField.Field)
		}

		if expectedField.Message != actualField.Message {
			t.Errorf("expected error fields item message at index %d is %s, got %s", i, expectedField.Message, actualField.Message)
		}
	}
}

func TestResponseErrorVM(t *testing.T) {
	testCases := []struct {
		Name     string
		Expected *ResponseErrorVM
		Actual   *ResponseErrorVM
	}{
		{
			Name: "NewResponseErrorVM",
			Expected: &ResponseErrorVM{
				ErrorFields: make([]*ResponseErrorFieldVM, 0),
			},
			Actual: NewResponseErrorVM(),
		},
		{
			Name: "SetMessage",
			Expected: &ResponseErrorVM{
				Message:     "message",
				ErrorFields: make([]*ResponseErrorFieldVM, 0),
			},
			Actual: NewResponseErrorVM().
				SetMessage("message"),
		},
		{
			Name: "SetMessage_EmptyString",
			Expected: &ResponseErrorVM{
				Message:     "",
				ErrorFields: make([]*ResponseErrorFieldVM, 0),
			},
			Actual: NewResponseErrorVM().
				SetMessage(""),
		},
		{
			Name: "AddErrorFields_SingleField",
			Expected: &ResponseErrorVM{
				ErrorFields: []*ResponseErrorFieldVM{
					{
						Field:   "email",
						Message: "email is required",
					},
				},
			},
			Actual: NewResponseErrorVM().
				AddErrorFields(NewResponseErrorFieldVM("email", "email is required")),
		},
		{
			Name: "AddErrorFields_MultipleFields",
			Expected: &ResponseErrorVM{
				ErrorFields: []*ResponseErrorFieldVM{
					{
						Field:   "email",
						Message: "email is required",
					},
					{
						Field:   "password",
						Message: "password is too short",
					},
					{
						Field:   "age",
						Message: "age must be positive",
					},
				},
			},
			Actual: NewResponseErrorVM().
				AddErrorFields(
					NewResponseErrorFieldVM("email", "email is required"),
					NewResponseErrorFieldVM("password", "password is too short"),
					NewResponseErrorFieldVM("age", "age must be positive"),
				),
		},
		{
			Name: "AddErrorFields_ChainedCalls",
			Expected: &ResponseErrorVM{
				ErrorFields: []*ResponseErrorFieldVM{
					{
						Field:   "name",
						Message: "name is required",
					},
					{
						Field:   "email",
						Message: "email format is invalid",
					},
				},
			},
			Actual: NewResponseErrorVM().
				AddErrorFields(NewResponseErrorFieldVM("name", "name is required")).
				AddErrorFields(NewResponseErrorFieldVM("email", "email format is invalid")),
		},
		{
			Name: "ParseError_ErrNil",
			Expected: &ResponseErrorVM{
				ErrorFields: make([]*ResponseErrorFieldVM, 0),
			},
			Actual: NewResponseErrorVM().
				ParseError(nil),
		},
		{
			Name: "ParseError_CustomError_NoFields",
			Expected: &ResponseErrorVM{
				Message:     "message",
				ErrorFields: make([]*ResponseErrorFieldVM, 0),
			},
			Actual: NewResponseErrorVM().
				ParseError(
					gocerr.New(
						http.StatusBadRequest,
						"message",
					),
				),
		},
		{
			Name: "ParseError_CustomError_WithSingleField",
			Expected: &ResponseErrorVM{
				Message: "message",
				ErrorFields: []*ResponseErrorFieldVM{
					{
						Field:   "field",
						Message: "message",
					},
				},
			},
			Actual: NewResponseErrorVM().
				ParseError(
					gocerr.New(
						http.StatusBadRequest,
						"message",
						gocerr.NewErrorField("field", "message"),
					),
				),
		},
		{
			Name: "ParseError_CustomError_WithMultipleFields",
			Expected: &ResponseErrorVM{
				Message: "validation failed",
				ErrorFields: []*ResponseErrorFieldVM{
					{
						Field:   "username",
						Message: "username is taken",
					},
					{
						Field:   "email",
						Message: "email is invalid",
					},
					{
						Field:   "password",
						Message: "password is too weak",
					},
				},
			},
			Actual: NewResponseErrorVM().
				ParseError(
					gocerr.New(
						http.StatusUnprocessableEntity,
						"validation failed",
						gocerr.NewErrorField("username", "username is taken"),
						gocerr.NewErrorField("email", "email is invalid"),
						gocerr.NewErrorField("password", "password is too weak"),
					),
				),
		},
		{
			Name: "ParseError_StandardError",
			Expected: &ResponseErrorVM{
				Message:     "message",
				ErrorFields: make([]*ResponseErrorFieldVM, 0),
			},
			Actual: NewResponseErrorVM().
				ParseError(errors.New("message")),
		},
		{
			Name: "ChainedMethods_CompleteErrorResponse",
			Expected: &ResponseErrorVM{
				Message: "custom message",
				ErrorFields: []*ResponseErrorFieldVM{
					{
						Field:   "manual_field",
						Message: "manually added error",
					},
				},
			},
			Actual: NewResponseErrorVM().
				SetMessage("custom message").
				AddErrorFields(NewResponseErrorFieldVM("manual_field", "manually added error")),
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			testResponseErrorVMEquality(
				t,
				testCases[i].Expected,
				testCases[i].Actual,
			)
		})
	}
}

// TestResponseErrorVM_MethodChaining tests that all methods return the same instance for proper chaining
func TestResponseErrorVM_MethodChaining(t *testing.T) {
	vm := NewResponseErrorVM()

	// Test that all methods return the same instance
	result1 := vm.SetMessage("test message")
	if result1 != vm {
		t.Error("SetMessage should return the same instance for method chaining")
	}

	result2 := vm.AddErrorFields(NewResponseErrorFieldVM("test", "test message"))
	if result2 != vm {
		t.Error("AddErrorFields should return the same instance for method chaining")
	}

	result3 := vm.ParseError(errors.New("test"))
	if result3 != vm {
		t.Error("ParseError should return the same instance for method chaining")
	}
}

// TestResponseErrorVM_SliceInitialization tests that ErrorFields slice is properly initialized
func TestResponseErrorVM_SliceInitialization(t *testing.T) {
	vm := NewResponseErrorVM()

	// Test that ErrorFields is not nil
	if vm.ErrorFields == nil {
		t.Error("ErrorFields should be initialized, not nil")
	}

	// Test that ErrorFields has zero length but is not nil
	if len(vm.ErrorFields) != 0 {
		t.Errorf("ErrorFields should have zero length, got %d", len(vm.ErrorFields))
	}

	// Test that we can append to the slice without panics
	vm.AddErrorFields(NewResponseErrorFieldVM("test", "test message"))

	if len(vm.ErrorFields) != 1 {
		t.Errorf("After adding one field, length should be 1, got %d", len(vm.ErrorFields))
	}
}
