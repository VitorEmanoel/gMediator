package mediator

import "testing"

type PingRequest struct {
	Request
	Message		string
}

type Test struct {

}

type PingRequestHandler struct {
	Test       string       `inject:"testando"`
}

func (p *PingRequestHandler) Handle(request PingRequest) (interface{}, error) {
	return request.Message + " Pong " + p.Test, nil
}

func TestNewMediator(t *testing.T) {
	var mediator = NewMediator()
	mediator.GetContainer().Inject("testando", "testeee")
	mediator.GetContainer().RegisterRequest(PingRequest{}, &PingRequestHandler{})
	response, err := mediator.Send(PingRequest{
		Message: "Ping",
	})
	if err != nil {
		t.Error("Error in send mediator. Error: ", err.Error())
		return
	}
	t.Log(response)
}
