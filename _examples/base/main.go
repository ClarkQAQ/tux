package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ClarkQAQ/tux"
)

func main() {
	t := tux.New()

	t.Use(func(c *tux.Context) {
		v := time.Now()
		c.Next()
		fmt.Printf("status: %d, url: %s, time: %s\n", c.Status(), c.Req.RequestURI, time.Since(v))
	})

	t.Use(func(c *tux.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.String(http.StatusInternalServerError, "Internal Server Error")
				c.End()
			}
		}()

		c.Next()
	})

	// Hello World
	t.Get("/", func(c *tux.Context) {
		c.String(200, "Hello, World!")
	})

	// 恐慌测试
	t.Get("/recover", func(c *tux.Context) {
		panic("panic")
	})

	// 阻塞测试
	t.Get("/block", func(c *tux.Context) {
		time.Sleep(5 * time.Second)
		c.String(200, "Hello, World!")
	})

	// 方法路由测试
	t.Object("/test", &TestApi{})

	// 导出路由
	for _, v := range t.ExportRoute() {
		fmt.Printf("path: %s, method: %#v\n", v.Path, v.Value)
	}

	if _, e := t.ServeAddr(":8080"); e != nil {
		panic(fmt.Sprintf("serve error: %s", e))
	}
}

type TestApi struct{}

func (t *TestApi) Get(c *tux.Context) {
	c.String(200, "Hello, Get!")
}

func (t *TestApi) Post(c *tux.Context) {
	c.String(200, "Hello, Post!")
}
