package tux

import (
	"path"
	"reflect"
	"strings"
)

// 路由组
type Group struct {
	tux        *Tux          // Tux
	prefix     string        // 组前缀
	field      *Group        // 上级组
	middleware []HandlerFunc // 组定义中间件
}

func (g *Group) NewGroup(prefix string) *Group {
	return &Group{
		tux:    g.tux,
		prefix: path.Join(g.prefix, prefix),
		field:  g,
	}
}

func (g *Group) upMiddlewares() (handles []HandlerFunc) {
	if g.field != nil {
		return append(g.field.upMiddlewares(), g.middleware...)
	}

	return g.middleware
}

func (g *Group) GroupTx(prefix string, f func(g *Group)) {
	f(g.NewGroup(prefix))
}

func (g *Group) Use(middleware ...HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *Group) Method(method, part string, handler HandlerFunc) {
	g.tux.tree.Set(strings.ToUpper(method)+"@"+path.Join(g.prefix, part),
		append(g.upMiddlewares(), handler))
}

func (g *Group) Get(part string, handler HandlerFunc) {
	g.Method("GET", part, handler)
}

func (g *Group) Post(part string, handler HandlerFunc) {
	g.Method("POST", part, handler)
}

func (g *Group) Put(part string, handler HandlerFunc) {
	g.Method("PUT", part, handler)
}

func (g *Group) Delete(part string, handler HandlerFunc) {
	g.Method("DELETE", part, handler)
}

func (g *Group) Patch(part string, handler HandlerFunc) {
	g.Method("PATCH", part, handler)
}

func (g *Group) Head(part string, handler HandlerFunc) {
	g.Method("HEAD", part, handler)
}

func (g *Group) Options(part string, handler HandlerFunc) {
	g.Method("OPTIONS", part, handler)
}

func (g *Group) Any(part string, handler HandlerFunc) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for _, method := range methods {
		g.Method(method, part, handler)
	}
}

func (g *Group) Object(part string, object interface{}) {
	// 通过反射获取结构体
	v := reflect.ValueOf(object)
	t := v.Type()

	if v.Kind() == reflect.Struct {
		newValue := reflect.New(t)
		newValue.Elem().Set(v)
		v = newValue
		t = v.Type()
	}

	g.GroupTx(part, func(r *Group) {
		for i := 0; i < v.NumMethod(); i++ {
			methodName := strings.ToUpper(t.Method(i).Name)

			switch methodName {
			case "MIDDLEWARE":
				if h, ok := v.Method(i).Interface().(func(*Context)); ok {
					r.Use(h)
				}
			default:
				if h, ok := v.Method(i).Interface().(func(*Context)); ok {
					r.Method(methodName, "/", h)
				}
			}
		}
	})
}
