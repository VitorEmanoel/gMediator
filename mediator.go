package mediator

import "context"

type ContextMediator struct {
	Container	Container
}

type Mediator interface {
	Send(ctx context.Context, request Request) (interface{}, error)
	GetContainer() Container
}

func NewMediator() Mediator{
	return &ContextMediator{Container:NewContainer()}
}

func (m *ContextMediator) GetContainer() Container {
	return m.Container
}

func (m *ContextMediator) Send(ctx context.Context, request Request) (interface{}, error) {
	return m.Container.ExecuteRequest(ctx, request)
}
