package gores

import "github.com/fikri240794/gocerr"

type ResponseErrorVM struct {
	Message     string                  `json:"message"`
	ErrorFields []*ResponseErrorFieldVM `json:"error_fields,omitempty"`
}

func NewResponseErrorVM() *ResponseErrorVM {
	return &ResponseErrorVM{}
}

func (vm *ResponseErrorVM) SetMessage(message string) *ResponseErrorVM {
	vm.Message = message

	return vm
}

func (vm *ResponseErrorVM) AddErrorFields(errorFields ...*ResponseErrorFieldVM) *ResponseErrorVM {
	vm.ErrorFields = append(vm.ErrorFields, errorFields...)

	return vm
}

func (vm *ResponseErrorVM) mapFromCustomError(err gocerr.Error) *ResponseErrorVM {
	vm = vm.SetMessage(err.Message)

	if len(err.ErrorFields) > 0 {
		for i := 0; i < len(err.ErrorFields); i++ {
			vm = vm.AddErrorFields(NewResponseErrorFieldVM(err.ErrorFields[i].Field, err.ErrorFields[i].Message))
		}
	}

	return vm
}

func (vm *ResponseErrorVM) ParseError(err error) *ResponseErrorVM {
	var (
		customError   gocerr.Error
		isCustomError bool
	)

	if err == nil {
		return vm
	}

	customError, isCustomError = err.(gocerr.Error)

	if isCustomError {
		vm = vm.mapFromCustomError(customError)
	} else {
		vm = vm.SetMessage(err.Error())
	}

	return vm
}
