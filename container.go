package mediator

import (
	"log"
	"reflect"
)

type Container interface {
	Inject(name string, data interface{})
	RegisterRequest(request Request, handler interface{})
	ExecuteRequest(request Request) (interface{}, error)
}

type ContainerContext struct {
	Handlers	map[reflect.Type]interface{}
	InjectValues      map[string]interface{}
}

func (c *ContainerContext) injectValues(data interface{}) {
	var dataType = reflect.TypeOf(data).Elem()
	var dataValue = reflect.ValueOf(data)
	for i := 0; i < dataType.NumField(); i++ {
		var field = dataType.Field(i)
		var injectName = field.Tag.Get("inject")
		if injectName == "" {
			continue
		}
		injectValue, exists := c.InjectValues[injectName]
		if !exists {
			continue
		}
		injectType := reflect.TypeOf(injectValue)
		if injectType.AssignableTo(field.Type) {
			dataValue.Elem().Field(i).Set(reflect.ValueOf(injectValue))
		}
	}
}

func (c *ContainerContext) Inject(name string, data interface{}) {
	c.InjectValues[name] = data
}

func NewContainer() Container {
	return &ContainerContext{
		Handlers: make(map[reflect.Type]interface{}),
		InjectValues: make(map[string]interface{}),
	}
}

func (c *ContainerContext) RegisterRequest(request Request, handler interface{}) {
	var requestType = reflect.TypeOf(request)
	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}
	c.injectValues(handler)
	c.Handlers[requestType] = handler
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
	var reflectValueHandler = reflect.ValueOf(handler)
	log.Println(reflectValueHandler.NumMethod())
	handleMethod := reflectValueHandler.MethodByName("Handle")
	var values = handleMethod.Call([]reflect.Value{reflect.ValueOf(request)})
	if len(values) != 2 {
		return nil, InvalidHandlerForRequest
	}
	value := values[0].Interface()
	errValue := values[1].Interface()
	err, ok := errValue.(error)
	if errValue != nil && !ok {
		return nil, InvalidHandlerForRequest
	}
	return value, err
}