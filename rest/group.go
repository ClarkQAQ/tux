package rest

import (
	"path"
	"reflect"
	"strings"
)

type Groupor interface {
	export(prevPrefix string, prevMiddlewareor []Middlewareor) ([]Methodor, error)
}

// 群组元数据
type Groupot struct {
	prefix  string
	tags    []string
	structs []any
}

func Group(prefix string, structs ...any) *Groupot {
	return &Groupot{
		prefix:  prefix,
		structs: structs,
	}
}

func (g *Groupot) Tags(ss ...string) *Groupot {
	g.tags = append(g.tags, ss...)
	return g
}

func (g *Groupot) export(prevPrefix string, prevMiddlewareor []Middlewareor) ([]Methodor, error) {
	mrs := []Methodor{}
	prefix := path.Join(prevPrefix, g.prefix)
	mws := append([]Middlewareor{}, prevMiddlewareor...)

	for i := 0; i < len(g.structs); i++ {
		if e := rangeStructMethod(g.structs[i], func(v reflect.Value, m reflect.Method) error {
			if v, ok := v.Interface().(func() Middlewareor); ok {
				mor := v()
				mor.setName(m.Name)

				mws = append(mws, mor)
			}

			return nil
		}); e != nil {
			return nil, e
		}

		if e := rangeStructMethod(g.structs[i], func(v reflect.Value, m reflect.Method) error {
			switch v := v.Interface().(type) {
			case func() Groupor:
				tmrs, e := v().export(prefix, mws)
				if e != nil {
					return e
				}

				mrs = append(mrs, tmrs...)
			case func() Methodor:
				mor := v()
				mor.setMethod(strings.ToUpper(m.Name))
				mor.joinPath(prefix)
				mor.joinMiddleware(mws...)

				mrs = append(mrs, mor)
			}

			return nil
		}); e != nil {
			return nil, e
		}
	}

	return mrs, nil
}
