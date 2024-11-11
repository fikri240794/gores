package gores

import (
	"net/http"

	"github.com/fikri240794/gocerr"
)

type ResponseVM[T comparable] struct {
	Code  int              `json:"code"`
	Error *ResponseErrorVM `json:"error,omitempty"`
	Data  T                `json:"data,omitempty"`
}

func NewResponseVM[T comparable]() *ResponseVM[T] {
	return &ResponseVM[T]{}
}

func (vm *ResponseVM[T]) SetCode(code int) *ResponseVM[T] {
	vm.Code = code

	return vm
}

func (vm *ResponseVM[T]) SetData(data T) *ResponseVM[T] {
	vm.Data = data

	return vm
}

func (vm *ResponseVM[T]) SetError(err *ResponseErrorVM) *ResponseVM[T] {
	vm.Error = err

	return vm
}

func (vm *ResponseVM[T]) SetErrorFromError(err error) *ResponseVM[T] {
	var (
		customError   gocerr.Error
		isCustomError bool
	)

	if err == nil {
		return vm
	}

	vm = vm.SetCode(http.StatusInternalServerError)

	customError, isCustomError = err.(gocerr.Error)

	if isCustomError {
		vm = vm.SetCode(customError.Code)
	}

	vm.Error = NewResponseErrorVM().
		ParseError(err)

	return vm
}
