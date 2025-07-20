package gores

import (
	"errors"
	"net/http"
	"testing"

	"github.com/fikri240794/gocerr"
)

// someStruct is a test data structure used for testing generic ResponseVM functionality
type someStruct struct {
	SomeField string
}

// testResponseVMEquality performs deep equality comparison between two ResponseVM instances.
// This helper function ensures comprehensive testing of all fields including nested structures.
func testResponseVMEquality(t *testing.T, expected, actual *ResponseVM[*someStruct]) {
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

	// Test HTTP status code
	if expected.Code != actual.Code {
		t.Errorf("expected code is %d, got %d", expected.Code, actual.Code)
	}

	// Test error field using dedicated helper
	testResponseErrorVMEquality(t, expected.Error, actual.Error)

	// Test data field with nil handling
	if expected.Data == nil && actual.Data != nil {
		t.Errorf("expected data is nil, got not nil")
		return
	}

	if expected.Data != nil && actual.Data == nil {
		t.Errorf("expected data is not nil, got nil")
		return
	}

	if expected.Data == nil && actual.Data == nil {
		return
	}

	// Test data field content
	if expected.Data.SomeField != actual.Data.SomeField {
		t.Errorf("expected data some field is %s, got %s", expected.Data.SomeField, actual.Data.SomeField)
	}
}

func TestResponseVM(t *testing.T) {
	testCases := []struct {
		Name     string
		Expected *ResponseVM[*someStruct]
		Actual   *ResponseVM[*someStruct]
	}{
		{
			Name:     "NewResponseVM",
			Expected: &ResponseVM[*someStruct]{},
			Actual:   NewResponseVM[*someStruct](),
		},
		{
			Name: "SetCode",
			Expected: &ResponseVM[*someStruct]{
				Code: http.StatusOK,
			},
			Actual: NewResponseVM[*someStruct]().
				SetCode(http.StatusOK),
		},
		{
			Name: "SetCode_ClientError",
			Expected: &ResponseVM[*someStruct]{
				Code: http.StatusBadRequest,
			},
			Actual: NewResponseVM[*someStruct]().
				SetCode(http.StatusBadRequest),
		},
		{
			Name: "SetCode_ServerError",
			Expected: &ResponseVM[*someStruct]{
				Code: http.StatusInternalServerError,
			},
			Actual: NewResponseVM[*someStruct]().
				SetCode(http.StatusInternalServerError),
		},
		{
			Name: "SetData",
			Expected: &ResponseVM[*someStruct]{
				Data: &someStruct{
					SomeField: "some field",
				},
			},
			Actual: NewResponseVM[*someStruct]().
				SetData(&someStruct{
					SomeField: "some field",
				}),
		},
		{
			Name: "SetData_EmptyStruct",
			Expected: &ResponseVM[*someStruct]{
				Data: &someStruct{},
			},
			Actual: NewResponseVM[*someStruct]().
				SetData(&someStruct{}),
		},
		{
			Name: "SetError",
			Expected: &ResponseVM[*someStruct]{
				Error: &ResponseErrorVM{
					ErrorFields: make([]*ResponseErrorFieldVM, 0),
				},
			},
			Actual: NewResponseVM[*someStruct]().
				SetError(NewResponseErrorVM()),
		},
		{
			Name: "SetError_WithMessage",
			Expected: &ResponseVM[*someStruct]{
				Error: &ResponseErrorVM{
					Message:     "test error",
					ErrorFields: make([]*ResponseErrorFieldVM, 0),
				},
			},
			Actual: NewResponseVM[*someStruct]().
				SetError(NewResponseErrorVM().SetMessage("test error")),
		},
		{
			Name:     "SetErrorFromError_ErrNil",
			Expected: &ResponseVM[*someStruct]{},
			Actual: NewResponseVM[*someStruct]().
				SetErrorFromError(nil),
		},
		{
			Name: "SetErrorFromError_CustomError",
			Expected: &ResponseVM[*someStruct]{
				Code: http.StatusBadRequest,
				Error: &ResponseErrorVM{
					Message:     "message",
					ErrorFields: make([]*ResponseErrorFieldVM, 0),
				},
			},
			Actual: NewResponseVM[*someStruct]().
				SetErrorFromError(
					gocerr.New(
						http.StatusBadRequest,
						"message",
					),
				),
		},
		{
			Name: "SetErrorFromError_CustomErrorWithFields",
			Expected: &ResponseVM[*someStruct]{
				Code: http.StatusUnprocessableEntity,
				Error: &ResponseErrorVM{
					Message: "validation failed",
					ErrorFields: []*ResponseErrorFieldVM{
						{
							Field:   "email",
							Message: "email is required",
						},
						{
							Field:   "password",
							Message: "password must be at least 8 characters",
						},
					},
				},
			},
			Actual: NewResponseVM[*someStruct]().
				SetErrorFromError(
					gocerr.New(
						http.StatusUnprocessableEntity,
						"validation failed",
						gocerr.NewErrorField("email", "email is required"),
						gocerr.NewErrorField("password", "password must be at least 8 characters"),
					),
				),
		},
		{
			Name: "SetErrorFromError_StandardError",
			Expected: &ResponseVM[*someStruct]{
				Code: http.StatusInternalServerError,
				Error: &ResponseErrorVM{
					Message:     "message",
					ErrorFields: make([]*ResponseErrorFieldVM, 0),
				},
			},
			Actual: NewResponseVM[*someStruct]().
				SetErrorFromError(
					errors.New("message"),
				),
		},
		{
			Name: "ChainedMethods_SuccessResponse",
			Expected: &ResponseVM[*someStruct]{
				Code: http.StatusCreated,
				Data: &someStruct{
					SomeField: "created resource",
				},
			},
			Actual: NewResponseVM[*someStruct]().
				SetCode(http.StatusCreated).
				SetData(&someStruct{
					SomeField: "created resource",
				}),
		},
		{
			Name: "ChainedMethods_ErrorResponse",
			Expected: &ResponseVM[*someStruct]{
				Code: http.StatusBadRequest,
				Error: &ResponseErrorVM{
					Message: "bad request",
					ErrorFields: []*ResponseErrorFieldVM{
						{
							Field:   "name",
							Message: "name is required",
						},
					},
				},
			},
			Actual: NewResponseVM[*someStruct]().
				SetCode(http.StatusBadRequest).
				SetErrorFromError(
					gocerr.New(
						http.StatusBadRequest,
						"bad request",
						gocerr.NewErrorField("name", "name is required"),
					),
				),
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			testResponseVMEquality(
				t,
				testCases[i].Expected,
				testCases[i].Actual,
			)
		})
	}
}

// TestResponseVM_MethodChaining tests that all methods return the same instance for proper chaining
func TestResponseVM_MethodChaining(t *testing.T) {
	vm := NewResponseVM[*someStruct]()

	// Test that all methods return the same instance
	result1 := vm.SetCode(http.StatusOK)
	if result1 != vm {
		t.Error("SetCode should return the same instance for method chaining")
	}

	result2 := vm.SetData(&someStruct{SomeField: "test"})
	if result2 != vm {
		t.Error("SetData should return the same instance for method chaining")
	}

	result3 := vm.SetError(NewResponseErrorVM())
	if result3 != vm {
		t.Error("SetError should return the same instance for method chaining")
	}

	result4 := vm.SetErrorFromError(errors.New("test"))
	if result4 != vm {
		t.Error("SetErrorFromError should return the same instance for method chaining")
	}
}
