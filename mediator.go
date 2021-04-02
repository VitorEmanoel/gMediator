package mediator

import "context"

var globalMediator *ContextMediator

type ContextMediator struct {
	Container	Container
}

type Options struct {
	Func func(request IRequest)
}

func WithContext(ctx context.Context) Options {
	return Options{
		Func: func(request IRequest) {
			request.WithContext(ctx)
		},
	}
}

type Mediator interface {
	Send(request IRequest, options ...Options) (interface{}, error)
	GetContainer() Container
}

func NewMediator() Mediator{
	mediator := &ContextMediator{Container: NewContainer()}
	globalMediator = mediator
	return mediator
}

func (m *ContextMediator) GetContainer() Container {
	return m.Container
}

func (m *ContextMediator) Send(request IRequest, options ...Options) (interface{}, error) {
	for _, option := range options {
		option.Func(request)
	}

	return m.Container.ExecuteRequest(request)
}

func Send(request IRequest, options ...Options) (interface{}, error) {
	if globalMediator == nil {
		return nil, NotExistsMediator
	}
	return globalMediator.Send(request, options...)
}