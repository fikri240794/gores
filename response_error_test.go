package gores

import (
	"errors"
	"net/http"
	"testing"

	"github.com/fikri240794/gocerr"
)

func testResponseErrorVMEquality(t *testing.T, expected, actual *ResponseErrorVM) {
	if expected == nil && actual != nil {
		t.Errorf("expected is nil, got not nil")
	}

	if expected != nil && actual == nil {
		t.Errorf("expected not nil, got nil")
	}

	if expected == nil && actual == nil {
		return
	}

	if expected.Message != actual.Message {
		t.Errorf("expected message is %s, got %s", expected.Message, actual.Message)
	}

	if len(expected.ErrorFields) != len(actual.ErrorFields) {
		t.Errorf("expected length of error fields is %d, got %d", len(expected.ErrorFields), len(actual.ErrorFields))
	} else {
		for i := range expected.ErrorFields {
			if expected.ErrorFields[i] != nil && actual.ErrorFields[i] == nil {
				t.Errorf("expected error fields item is not nil, but got nil")
			}

			if expected.ErrorFields[i] == nil && actual.ErrorFields[i] != nil {
				t.Errorf("expected error fields item is nil, but got not nil")
			}

			if expected.ErrorFields[i] == nil && actual.ErrorFields[i] == nil {
				continue
			}

			if expected.ErrorFields[i].Field != actual.ErrorFields[i].Field {
				t.Errorf("expected error fields item field is %s, got %s", expected.ErrorFields[i].Field, actual.ErrorFields[i].Field)
			}

			if expected.ErrorFields[i].Message != actual.ErrorFields[i].Message {
				t.Errorf("expected error fields item message is %s, got %s", expected.ErrorFields[i].Message, actual.ErrorFields[i].Message)
			}
		}
	}
}

func TestResponseErrorVM(t *testing.T) {
	var testCases []struct {
		Name     string
		Expected *ResponseErrorVM
		Actual   *ResponseErrorVM
	} = []struct {
		Name     string
		Expected *ResponseErrorVM
		Actual   *ResponseErrorVM
	}{
		{
			Name:     "NewResponseErrorVM",
			Expected: &ResponseErrorVM{},
			Actual:   NewResponseErrorVM(),
		},
		{
			Name: "SetMessage",
			Expected: &ResponseErrorVM{
				Message: "message",
			},
			Actual: NewResponseErrorVM().
				SetMessage("message"),
		},
		{
			Name:     "ParseError_ErrNil",
			Expected: &ResponseErrorVM{},
			Actual: NewResponseErrorVM().
				ParseError(nil),
		},
		{
			Name: "ParseError_CustomError",
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
			Name: "ParseError_StandardError",
			Expected: &ResponseErrorVM{
				Message: "message",
			},
			Actual: NewResponseErrorVM().
				ParseError(errors.New("message")),
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
