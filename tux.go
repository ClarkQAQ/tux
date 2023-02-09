package tux

import (
	"net"
	"net/http"
	"sync"

	"github.com/ClarkQAQ/tux/tree"
)

type Tux struct {
	*Group                         // 路由组
	tree   *tree.Tree[HandlerList] // 路由树

	contextPool *sync.Pool
}

func New() *Tux {
	tux := &Tux{}

	tux.Group = &Group{tux, "/", nil, nil}
	tux.tree = &tree.Tree[HandlerList]{}
	tux.contextPool = &sync.Pool{
		New: func() interface{} {
			return tux.newContext()
		},
	}

	return tux
}

func (tux *Tux) ExportRoute() []*tree.ExportValue[HandlerList] {
	return tux.tree.ExportTreeMethon("")
}

func defaultNotFound(c *Context) {
	c.Writer.WriteHeader(http.StatusNotFound)
	c.Writer.Write([]byte("404 NOT FOUND:" + c.Req.URL.Path))
}

func (tux *Tux) Handle(c *Context) {
	c.handlerList, c.vpath = tux.tree.Get(c.Req.Method + "@" + c.Req.URL.Path)
	if len(c.handlerList) == 0 {
		c.handlerList = append(tux.middleware, defaultNotFound)
	}

	c.Next()
}

func (tux *Tux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := tux.contextPool.Get().(*Context)
	defer func(t *Tux, ctx *Context) {
		ctx.reset()
		t.contextPool.Put(c)
	}(tux, c)

	c.use(w, req)

	tux.Handle(c)

	if c.index < -1 {
		panic(nil)
	}

	if c.exported {
		return
	}

	w.WriteHeader(c.Writer.status)

	if _, e := w.Write(c.Writer.body.Bytes()); e != nil {
		panic(e)
	}
}

func (tux *Tux) ServeAddr(addr string) (*http.Server, error) {
	net, e := net.Listen("tcp", addr)
	if e != nil {
		return nil, e
	}

	return tux.ServeListener(net)
}

func (tux *Tux) ServeListener(l net.Listener) (*http.Server, error) {
	http := &http.Server{Handler: tux}

	return http, http.Serve(l)
}
