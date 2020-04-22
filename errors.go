package mediator

type NotFoundHandlerError struct {
}

func (err *NotFoundHandlerError) Error() string {
	return "handler not found in container"
}
