package mediator

import "reflect"

type CallbackType string

const (
	Before CallbackType = "BEFORE"
	After CallbackType = "AFTER"
)

type Callback struct {
	Type    CallbackType
	Func    interface{}
}

func validateFunc(callbackFunc interface{}) bool {
	funcType := reflect.TypeOf(callbackFunc)
	return funcType.Kind() == reflect.Func
}

func NewCallback(callbackType CallbackType, callbackFunc interface{}) *Callback {
	if !validateFunc(callbackFunc) {
		return nil
	}
	return &Callback{Type: callbackType, Func: callbackFunc}
}
