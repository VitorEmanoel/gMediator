package mediator

import "reflect"

type Container interface {
	RegisterRequest(request Request, handler RequestHandler)
	ExecuteRequest(request Request) (interface{}, error)
}

type ContainerContext struct {
	Handlers	map[reflect.Type]*RequestHandler
}

func NewContainer() Container {
	return &ContainerContext{Handlers: make(map[reflect.Type]*RequestHandler)}
}

func (c *ContainerContext) RegisterRequest(request Request, handler RequestHandler) {
	var requestType = reflect.TypeOf(request)
	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}
	c.Handlers[requestType] = &handler
}

func (c *ContainerContext) ExecuteRequest(request Request) (interface{}, error){
	var requestType = reflect.TypeOf(request)
	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}
	handler, ok := c.Handlers[requestType]
	if !ok {
		return nil, &NotFoundHandlerError{}
	}
	var objectHandler = *handler
	return objectHandler.Handle(request)
}