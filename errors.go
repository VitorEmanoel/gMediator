package mediator

import "errors"

var InvalidHandlerForRequest = errors.New("invalid handler for request")

var NotFoundHandler = errors.New("not found handler")

var NotExistsMediator = errors.New("not exists mediator")

var NotExistsContainer = errors.New("not exists container")