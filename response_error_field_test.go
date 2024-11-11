package gores

import (
	"testing"
)

func TestNewResponseErrorFieldVM(t *testing.T) {
	field := "exampleField"
	message := "exampleMessage"

	vm := NewResponseErrorFieldVM(field, message)

	if vm.Field != field {
		t.Errorf("Expected field is %s, got %s", field, vm.Field)
	}

	if vm.Message != message {
		t.Errorf("Expected message is %s, got %s", message, vm.Message)
	}
}
