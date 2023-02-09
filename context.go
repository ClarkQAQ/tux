package tux

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"sync"
)

type HandlerFunc func(*Context)

type HandlerList []HandlerFunc

type HandlerWriter struct {
	status int           // 状态码
	header http.Header   // headers
	body   *bytes.Buffer // body
}

// Handler 上下文
type Context struct {
	tux *Tux // Tux

	exported bool                // 原始 writer 是否被导出
	writer   http.ResponseWriter // 原始的 writer
	Req      *http.Request       // 公开请求
	Writer   *HandlerWriter      // 公开响应

	ctxStore     map[string]interface{} // 上下文存储
	ctxStoreLock *sync.RWMutex          // 上下文存储锁

	vpath       []uint8     // 路径
	index       int         // 当前执行的处理函数索引
	handlerList HandlerList // 处理函数列表
}

func newWriterBufferContext() *HandlerWriter {
	return &HandlerWriter{

		status: 200,
		body:   &bytes.Buffer{},
	}
}

func (w *HandlerWriter) reset() {
	w.status = defaultStatusCode
	w.header = nil
	w.body.Reset()
}

func (w *HandlerWriter) Header() http.Header {
	return w.header
}

func (w *HandlerWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *HandlerWriter) WriteHeader(code int) {
	w.status = code
}

func (tux *Tux) newContext() *Context {
	c := &Context{
		tux: tux,

		ctxStore:     make(map[string]interface{}),
		ctxStoreLock: &sync.RWMutex{},
		index:        -1,
	}

	c.Writer = newWriterBufferContext()
	return c
}

func (c *Context) reset() {
	c.exported = false
	c.writer = nil
	c.Req = nil
	c.Writer.reset()

	c.vpath = nil
	c.ctxStore = nil
	c.index = -1
	c.handlerList = nil
}

func (c *Context) use(writer http.ResponseWriter, req *http.Request) {
	c.writer = writer
	c.Req = req

	c.Writer.header = writer.Header()
}

func (c *Context) Set(key string, value interface{}) {
	c.ctxStoreLock.Lock()
	defer c.ctxStoreLock.Unlock()

	if c.ctxStore == nil {
		c.ctxStore = make(map[string]interface{})
	}

	c.ctxStore[key] = value
}

func (c *Context) Get(key string) (value interface{}, ok bool) {
	c.ctxStoreLock.RLock()
	defer c.ctxStoreLock.RUnlock()

	if c.ctxStore == nil {
		return nil, false
	}

	value, ok = c.ctxStore[key]
	return
}

func (c *Context) Delete(key string) {
	c.ctxStoreLock.Lock()
	defer c.ctxStoreLock.Unlock()

	if c.ctxStore == nil {
		return
	}

	delete(c.ctxStore, key)
}

// func (c *Context) Param(key string) string {
// 	if v, ok := c.Params[key]; ok {
// 		return v
// 	}

// 	return ""
// }

// URL Query
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// POST Form Value
func (c *Context) PostForm(key string) string {
	c.Req.ParseForm()
	return c.Req.FormValue(key)
}

// 获取当前请求的Cookit
func (c *Context) Cookie(key string) string {
	if v, _ := c.Req.Cookie(key); v != nil {
		return v.Value
	}

	return ""
}

func (c *Context) ReqBody() (body []byte) {
	if c.Req.Body != nil {
		body, _ = io.ReadAll(c.Req.Body)
	}
	c.Req.Body = io.NopCloser(bytes.NewBuffer(body))
	return body
}

// 设置状态码
// 也可以获取当前设置的状态码
// 加了魔法, 可以重复设置状态码
func (c *Context) Status(code ...int) int {
	if len(code) > 0 {
		c.Writer.WriteHeader(code[0])
	}

	return c.Writer.status
}

// 设置 header 的简单封装
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 占位符填充输出
// 内部封装了fmt.Sprintf
// 默认content-type为text/html
func (c *Context) Sprintf(code int, format string, values ...any) {
	c.Status(code)
	if c.Writer.Header().Get(HeaderContentType) == "" {
		c.SetHeader(HeaderContentType, "text/html;  charset=utf-8")
	}

	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 输出字符串
// 默认content-type为text/html
func (c *Context) String(code int, format string) {
	c.Status(code)
	if c.Writer.Header().Get(HeaderContentType) == "" {
		c.SetHeader(HeaderContentType, "text/html;  charset=utf-8")
	}

	c.Writer.Write([]byte(format))
}

func (c *Context) JSON(code int, obj any) {
	c.SetHeader(HeaderContentType, "application/json; charset=utf-8")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) XML(code int, obj any) {
	c.SetHeader(HeaderContentType, "application/xml; charset=utf-8")
	c.Status(code)
	encoder := xml.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 输出字节数据
// 默认content-type为application/octet-stream
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	if c.Writer.Header().Get(HeaderContentType) == "" {
		c.SetHeader(HeaderContentType, "application/octet-stream")
	}

	c.Writer.Write(data)
}

func (c *Context) File(code int, ffs fs.FS, filename string) {
	c.Status(code)
	if c.Writer.Header().Get(HeaderContentType) == "" {
		c.SetHeader(HeaderContentType, "text/html;  charset=utf-8")
	}

	b, e := fs.ReadFile(ffs, filename)
	if e != nil {
		http.Error(c.Writer, e.Error(), 500)
		return
	}

	c.Writer.Write(b)
}

func (c *Context) WriteTo(w io.Writer) (int64, error) {
	return c.Writer.body.WriteTo(w)
}

// Hijacker 接口封装
// 用于将writer转换为bufio.ReadWriter
// 一旦调用此函数, 其他输出函数将无效
// 并且此函数关闭时将调用End方法结束请求
func (c *Context) Hijacker(f func(bufrw *bufio.ReadWriter) error) error {
	hj, ok := c.writer.(http.Hijacker)
	if !ok {
		return ErrResponseAlreadySent
	}

	c.exported = true

	conn, bufrw, e := hj.Hijack()
	if e != nil {
		return e
	}

	defer conn.Close()

	if e := f(bufrw); e != nil {
		return e
	}

	c.End()
	return nil
}

func (c *Context) Response() (http.ResponseWriter, error) {
	if c.exported {
		return nil, ErrWriterAlreadyExported
	}

	c.exported = true
	return c.writer, nil
}

func (c *Context) ResponseExported() bool {
	return c.exported
}

// 循环执行下一个HandlerFunc
func (c *Context) Next() {
	for c.index++; c.index > -1 && c.index < len(c.handlerList); c.index++ {
		func() {
			defer func() {
				if r := recover(); r != nil || c.index < -1 {
					panic(r)
				}
			}()

			c.handlerList[c.index](c)
		}()
	}
}

// handler 调用索引
func (c *Context) Index() int {
	return c.index
}

// 清除缓冲区
func (c *Context) Clean() {
	c.Writer.reset()
	c.Writer.header = c.writer.Header()
}

// 结束请求, 并跳过后续的HandlerFunc
// 将正常输出缓冲区内容
func (c *Context) End() {
	c.index = len(c.handlerList)
	panic(nil)
}

// 重置请求, 并跳过后续的HandlerFunc
// 将无任何输出, 并且浏览器显示连接已重置
// 但是仍然有响应头 "HTTP 1.1 400 Bad Request\r\nConnection: close"
func (c *Context) Close() {
	c.index = -100
	panic(nil)
}
