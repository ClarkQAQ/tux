package rest

import (
	"context"
	"fmt"
	"path"
)

// 端点接口
type Methodor interface {
	newInputPointer() interface{}
	setMethod(method string)
	joinPath(prefix string)
	joinMiddleware(mws ...Middlewareor)
	exportString() string

	Invoke(ctx context.Context, input interface{}) (interface{}, error)
}

// 端点结构
type Methodot[i, o any] struct {
	method     string         // 请求方法
	path       string         // 请求路径
	middleware []Middlewareor // 中间件
	handler    func(ctx context.Context, input *i) (o, error)

	tags        []string // 标签
	summary     string   // 摘要
	description string   // 描述
}

// 新建端点
func Method[req, res any](handler func(ctx context.Context, input *req) (res, error)) *Methodot[req, res] {
	return &Methodot[req, res]{
		handler: handler,
	}
}

func (m *Methodot[i, o]) Handler(handler func(ctx context.Context, input *i) (o, error)) *Methodot[i, o] {
	m.handler = handler
	return m
}

func (m *Methodot[i, o]) Tags(ss ...string) *Methodot[i, o] {
	m.tags = append(m.tags, ss...)
	return m
}

func (m *Methodot[i, o]) Summary(summary string) *Methodot[i, o] {
	m.summary = summary
	return m
}

func (m *Methodot[i, o]) Description(desc string) *Methodot[i, o] {
	m.description = desc
	return m
}

func (m *Methodot[req, res]) Invoke(ctx context.Context, input interface{}) (interface{}, error) {
	inp, ok := input.(*req)
	if !ok {
		return nil, fmt.Errorf("%w of input: %T, expected: %T", ErrInvalidType, input, new(req))
	}

	return m.handler(ctx, inp)
}

func (m *Methodot[req, res]) newInputPointer() interface{} {
	return new(req)
}

func (m *Methodot[req, res]) setMethod(method string) {
	m.method = method
}

func (m *Methodot[req, res]) joinPath(prefix string) {
	m.path = path.Join(m.path, prefix)
}

func (m *Methodot[req, res]) joinMiddleware(mws ...Middlewareor) {
	m.middleware = append(m.middleware, mws...)
}

func (m *Methodot[req, res]) exportString() string {
	mwNames := ""
	for i, mw := range m.middleware {
		mwNames += fmt.Sprintf("%d: %s ", i, mw.Name())
	}

	return fmt.Sprintf("Method: %s - Path: %s Handler: %T Middleware: %v",
		m.method, m.path, m.handler, mwNames)
}
