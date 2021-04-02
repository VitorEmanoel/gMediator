package mediator

import (
	"context"
)

type IRequest interface {
	Context()                           context.Context
	WithContext(ctx context.Context)
}

type Request struct {
	Ctx     context.Context
}

func (r *Request) Context() context.Context{
	return r.Ctx
}

func (r *Request) WithContext(ctx context.Context) {
	r.Ctx = ctx
}


type RequestHandler interface {
	Handle(ctx context.Context, request Request) (interface{}, error)
}
