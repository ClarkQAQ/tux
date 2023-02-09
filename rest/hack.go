package rest

import "context"

const (
	RequestHack HackType = iota << 1
	ResponseHack
)

type HackType int8

type Hackor interface {
	setName(name string)

	Name() string
	Invoke(ctx context.Context)
}

// 群组元数据
type Hackot struct {
	name     string
	hackType HackType
	handler  func(ctx context.Context)
}

func Hack(htp HackType, handler func(ctx context.Context)) *Hackot {
	return &Hackot{
		hackType: htp,
		handler:  handler,
	}
}

func (m *Hackot) setName(name string) {
	m.name = name
}

func (m *Hackot) Name() string {
	return m.name
}

func (m *Hackot) Invoke(ctx context.Context) {
	m.handler(ctx)
}
