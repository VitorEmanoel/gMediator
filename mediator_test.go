package mediator

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type PingRequest struct {
	Request
	Message		string
}

func cleanGlobals() {
	globalMediator = nil
	globalContainer = nil
}

type PingRequestHandler struct {
	SecondMessage       string       `inject:"secondMessage"`
}

func (p *PingRequestHandler) Handle(request *PingRequest) (interface{}, error) {
	return request.Message + " Pong" + p.SecondMessage, nil
}

func TestSimplesMediator(t *testing.T) {
	var mediator = NewMediator()
	mediator.GetContainer().RegisterRequest(&PingRequest{}, &PingRequestHandler{})
	response, err := mediator.Send(&PingRequest{
		Message: "Ping",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Ping Pong", response)
	t.Cleanup(cleanGlobals)
}

func TestMediatorWithInject(t *testing.T) {
	var mediator = NewMediator()
	mediator.GetContainer().Inject("secondMessage", " Pong")
	mediator.GetContainer().RegisterRequest(&PingRequest{}, &PingRequestHandler{})
	response, err := mediator.Send(&PingRequest{
		Message: "Ping",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Ping Pong Pong", response)
	t.Cleanup(cleanGlobals)
}

func TestMediatorWithGlobalRegister(t *testing.T) {
	var mediator = NewMediator()
	RegisterRequest(&PingRequest{}, &PingRequestHandler{})
	response, err := mediator.Send(&PingRequest{
		Message: "Ping",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Ping Pong", response)
	t.Cleanup(cleanGlobals)
}

func TestMediatorWithGlobalRegisterBeforeCreate(t *testing.T) {
	RegisterRequest(&PingRequest{}, &PingRequestHandler{})
	var mediator = NewMediator()
	response, err := mediator.Send(&PingRequest{
		Message: "Ping",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Ping Pong", response)
	t.Cleanup(cleanGlobals)
}

func TestMediatorWithGlobalSend(t *testing.T) {
	NewMediator()
	RegisterRequest(&PingRequest{}, &PingRequestHandler{})
	response, err := Send(&PingRequest{
		Message: "Ping",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Ping Pong", response)
	t.Cleanup(cleanGlobals)
}

func TestMediatorWithGlobalSendWithoutCreateMediator(t *testing.T) {
	RegisterRequest(&PingRequest{}, &PingRequestHandler{})
	response, err := Send(&PingRequest{
		Message: "Ping",
	})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, NotExistsMediator, err)
	t.Cleanup(cleanGlobals)
}

type TestContextPingRequest struct {
	Request
}

type TestContextPingRequestHandler struct {

}

var InvalidPing = errors.New("invalid ping")

func (t *TestContextPingRequestHandler) Handle(request *TestContextPingRequest) (string, error) {
	value := request.Context().Value("test")
	ping, ok := value.(string)
	if !ok {
		return "", InvalidPing
	}
	return ping + " Pong", nil

}

func TestMediatorContext(t *testing.T) {
	NewMediator()
	RegisterRequest(&TestContextPingRequest{}, &TestContextPingRequestHandler{})
	var ctx = context.WithValue(context.Background(), "test", "Ping")
	response, err := Send(&TestContextPingRequest{}, WithContext(ctx))
	assert.Nil(t, err)
	assert.Equal(t, "Ping Pong", response)
	t.Cleanup(cleanGlobals)
}

func ChangeContextCallback(request *TestContextPingRequest) {
	request.WithContext(context.WithValue(request.Context(), "test", "Ping Ping"))
}

func TestMediatorWithBeforeCallback(t *testing.T) {
	var mediator = NewMediator()
	RegisterRequest(&TestContextPingRequest{}, &TestContextPingRequestHandler{})
	mediator.GetContainer().RegisterCallbacks(&TestContextPingRequest{}, NewCallback(Before, ChangeContextCallback))
	var ctx = context.WithValue(context.Background(), "test", "Ping")
	response, err := Send(&TestContextPingRequest{}, WithContext(ctx))
	assert.Nil(t, err)
	assert.Equal(t, "Ping Ping Pong", response)
	t.Cleanup(cleanGlobals)
}

func ChangeResponseCallback(request *TestContextPingRequest, value string, err error) (string, error){
	return value + " Pong", err
}

func TestMediatorWithAfterCallback(t *testing.T) {
	var mediator = NewMediator()
	RegisterRequest(&TestContextPingRequest{}, &TestContextPingRequestHandler{})
	mediator.GetContainer().RegisterCallbacks(&TestContextPingRequest{}, NewCallback(After, ChangeResponseCallback))
	var ctx = context.WithValue(context.Background(), "test", "Ping")
	response, err := Send(&TestContextPingRequest{}, WithContext(ctx))
	assert.Nil(t, err)
	assert.Equal(t, "Ping Pong Pong", response)
	t.Cleanup(cleanGlobals)
}

func TestMediatorWithAfterAndBeforeCallback(t *testing.T) {
	var mediator = NewMediator()
	RegisterRequest(&TestContextPingRequest{}, &TestContextPingRequestHandler{})
	mediator.GetContainer().RegisterCallbacks(&TestContextPingRequest{}, NewCallback(After, ChangeResponseCallback),  NewCallback(Before, ChangeContextCallback))
	var ctx = context.WithValue(context.Background(), "test", "Ping")
	response, err := Send(&TestContextPingRequest{}, WithContext(ctx))
	assert.Nil(t, err)
	assert.Equal(t, "Ping Ping Pong Pong", response)
	t.Cleanup(cleanGlobals)
}

func TestMediatorWithAfterCallbackGlobalRegister(t *testing.T) {
	RegisterRequest(&TestContextPingRequest{}, &TestContextPingRequestHandler{})
	RegisterCallbacks(&TestContextPingRequest{}, NewCallback(After, ChangeResponseCallback), NewCallback(Before, ChangeContextCallback))
	NewMediator()
	var ctx = context.WithValue(context.Background(), "test", "Ping")
	response, err := Send(&TestContextPingRequest{}, WithContext(ctx))
	assert.Nil(t, err)
	assert.Equal(t, "Ping Ping Pong Pong", response)
	t.Cleanup(cleanGlobals)
}