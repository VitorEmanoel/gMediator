package mediator

import "testing"

type PingRequest struct {
	Request
	Message		string
}

type PingRequestHandler struct {
}

func (p *PingRequestHandler) Handle(request Request) (interface{}, error) {
	var pingRequest = request.(PingRequest)
	return pingRequest.Message + " Pong", nil
}

func TestNewMediator(t *testing.T) {
	var mediator = NewMediator()
	mediator.GetContainer().RegisterRequest(PingRequest{}, &PingRequestHandler{})
	response, err := mediator.Send(PingRequest{
		Message: "Teste",
	})
	if err != nil {
		t.Error("Error in send mediator. Error: ", err.Error())
		return
	}
	t.Log(response)
}
