package errors

import "core/types/httptypes"

type CalError struct {
	BaseError error
	Status *httptypes.HttpResponseStatus
}

func (e *CalError) Error() string {
	return e.BaseError.Error()
}