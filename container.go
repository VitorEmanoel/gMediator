package mediator

import (
	"reflect"
)

var globalContainer *ContainerContext

var registeredHandlers = make(map[reflect.Type]interface{})
var registeredCallbacks = make(map[reflect.Type][]*Callback)

func RegisterRequest(request IRequest, handler interface{}) {
	if globalContainer == nil{
		var requestType = reflect.TypeOf(request)
		if requestType.Kind() == reflect.Ptr {
			requestType = requestType.Elem()
		}
		registeredHandlers[requestType] = handler
		return
	}
	globalContainer.RegisterRequest(request, handler)
}

func RegisterCallbacks(request IRequest, callbacks ...*Callback) {
	if globalContainer == nil {
		var requestType = reflect.TypeOf(request)
		if requestType.Kind() == reflect.Ptr {
			requestType = requestType.Elem()
		}
		localCallbacks := registeredCallbacks[requestType]
		localCallbacks = append(localCallbacks, callbacks...)
		registeredCallbacks[requestType] = localCallbacks
		return
	}
	globalContainer.RegisterCallbacks(request, callbacks...)
}

type Handler struct {
	Handler     interface{}
	Callbacks   []*Callback
}

type Container interface {
	Inject(name string, data interface{})
	RegisterRequest(request IRequest, handler interface{})
	RegisterCallbacks(request IRequest, callback ...*Callback)
	ExecuteRequest(request IRequest) (interface{}, error)
}

type ContainerContext struct {
	Handlers            map[reflect.Type]*Handler
	InjectValues        map[string]interface{}
}

func (c *ContainerContext) injectValues(data interface{}, inject string) {
	var dataType = reflect.TypeOf(data).Elem()
	var dataValue = reflect.ValueOf(data)
	for i := 0; i < dataType.NumField(); i++ {
		var field = dataType.Field(i)
		var injectName = field.Tag.Get("inject")
		if injectName == "" || (inject != "" && inject != injectName){
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
	for _, handler := range c.Handlers {
		c.injectValues(handler.Handler, name)
	}
}

func (c *ContainerContext) register(requestType reflect.Type, handler interface{}) {
	c.injectValues(handler, "")
	handlerValue := Handler{
		Handler:   handler,
	}
	c.Handlers[requestType] = &handlerValue
}

func (c *ContainerContext) RegisterRequest(request IRequest, handler interface{}) {
	var requestType = reflect.TypeOf(request)
	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}
	c.register(requestType, handler)
}

func (c *ContainerContext) RegisterCallbacks(request IRequest, callbacks ...*Callback) {
	var requestType = reflect.TypeOf(request)
	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}
	handler, ok := c.Handlers[requestType]
	if !ok {
		return
	}
	handler.Callbacks = append(handler.Callbacks, callbacks...)
}


func (c *ContainerContext) callCallbacks(handler *Handler, callbackType CallbackType, values... reflect.Value) []reflect.Value{
	for _, callback := range handler.Callbacks {
		if callback.Type == callbackType {
			var callbackValue = reflect.ValueOf(callback.Func)
			return callbackValue.Call(values)
		}
	}
	return nil
}

func (c *ContainerContext) ExecuteRequest(request IRequest) (interface{}, error){
	var requestType = reflect.TypeOf(request)
	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}
	handler, ok := c.Handlers[requestType]
	if !ok {
		return nil, NotFoundHandler
	}
	var reflectValueHandler = reflect.ValueOf(handler.Handler)
	handleMethod := reflectValueHandler.MethodByName("Handle")
	c.callCallbacks(handler, Before, reflect.ValueOf(request))
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
	finalValues := c.callCallbacks(handler, After, reflect.ValueOf(request), values[0], values[1])
	if finalValues != nil {
		value = finalValues[0].Interface()
		errValue := finalValues[1].Interface()
		err, ok = errValue.(error)
		if errValue != nil && !ok {
			return nil, InvalidHandlerForRequest
		}
	}
	return value, err
}

func NewContainer() Container {
	var ctx = &ContainerContext{
		Handlers: make(map[reflect.Type]*Handler),
		InjectValues: make(map[string]interface{}),
	}
	for request, handler := range registeredHandlers {
		ctx.register(request, handler)
		if _, exists := registeredCallbacks[request]; !exists {
			continue
		}
		ctx.Handlers[request].Callbacks = append(ctx.Handlers[request].Callbacks, registeredCallbacks[request]...)
	}
	globalContainer = ctx
	return ctx
}
