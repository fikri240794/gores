package gores

import (
	"errors"
	"net/http"
	"testing"

	"github.com/fikri240794/gocerr"
)

type someStruct struct {
	SomeField string
}

func testResponseVMEquality(t *testing.T, expected, actual *ResponseVM[*someStruct]) {
	if expected == nil && actual != nil {
		t.Errorf("expected is nil, got not nil")
	}

	if expected != nil && actual == nil {
		t.Errorf("expected not nil, got nil")
	}

	if expected == nil && actual == nil {
		return
	}

	if expected.Code != actual.Code {
		t.Errorf("expected code is %d, got %d", expected.Code, actual.Code)
	}

	testResponseErrorVMEquality(t, expected.Error, actual.Error)

	if expected.Data == nil && actual.Data != nil {
		t.Errorf("expected data is nil, got not nil")
	}

	if expected.Data != nil && actual.Data == nil {
		t.Errorf("expected data is not nil, got nil")
	}

	if expected.Data == nil && actual.Data == nil {
		return
	}

	if expected.Data.SomeField != actual.Data.SomeField {
		t.Errorf("expected data some field is %s, got %s", expected.Data.SomeField, actual.Data.SomeField)
	}
}

func TestResponseVM(t *testing.T) {
	var testCases []struct {
		Name     string
		Expected *ResponseVM[*someStruct]
		Actual   *ResponseVM[*someStruct]
	} = []struct {
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
			Name: "SetError",
			Expected: &ResponseVM[*someStruct]{
				Error: &ResponseErrorVM{},
			},
			Actual: NewResponseVM[*someStruct]().
				SetError(NewResponseErrorVM()),
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
					Message: "message",
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
			Name: "SetErrorFromError_StandardError",
			Expected: &ResponseVM[*someStruct]{
				Code: http.StatusInternalServerError,
				Error: &ResponseErrorVM{
					Message: "message",
				},
			},
			Actual: NewResponseVM[*someStruct]().
				SetErrorFromError(
					errors.New("message"),
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
