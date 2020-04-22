package mediator

type Request interface {
}


type RequestHandler interface {
	Handle(request Request) (interface{}, error)
}
