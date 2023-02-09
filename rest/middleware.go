package rest

import "context"

type Middlewareor interface {
	setName(name string)

	Name() string
	Invoke(ctx context.Context)
}

// 群组元数据
type Middlewareot struct {
	name    string
	handler func(ctx context.Context)
}

func Middleware(handler func(ctx context.Context)) *Middlewareot {
	return &Middlewareot{
		handler: handler,
	}
}

func (m *Middlewareot) setName(name string) {
	m.name = name
}

func (m *Middlewareot) Name() string {
	return m.name
}

func (m *Middlewareot) Invoke(ctx context.Context) {
	m.handler(ctx)
}
