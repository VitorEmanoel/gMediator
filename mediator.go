package mediator

type ContextMediator struct {
	Container	Container
}

type Mediator interface {
	Send(request Request) (interface{}, error)
	GetContainer() Container
}

func NewMediator() Mediator{
	return &ContextMediator{Container:NewContainer()}
}

func (m *ContextMediator) GetContainer() Container {
	return m.Container
}

func (m *ContextMediator) Send(request Request) (interface{}, error) {
	return m.Container.ExecuteRequest(request)
}
