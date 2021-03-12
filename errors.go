package mediator

import "errors"

var InvalidHandlerForRequest = errors.New("invalid handler for request")

type NotFoundHandlerError struct {
}

func (err *NotFoundHandlerError) Error() string {
	return "handler not found in container"
}
