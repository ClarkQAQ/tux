package tree

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	defaultSeparator = "/" // 默认路径分隔符
)

type Value struct{}

func TestRouteTree_Get(t *testing.T) {
	var tree = &Tree[*Value]{}
	tree.Set("/api/qaq/:id", &Value{})
	tree.Set("/api/qaq/:id/qwq", &Value{})
	tree.Set("/api/qaq/:id/qwq/:name", &Value{})
	tree.Set("/api/asterisk/*path", &Value{})

	testGetHandler := func(path string) {
		handler, vpath := tree.Get(path)
		if handler == nil {
			t.Fatalf("handler must not be nil, path: %s", path)
		}

		fmt.Printf("test get: %s, vpath: %s\n", path, vpath)

	}

	testGetHandler("/api/qaq/1")
	testGetHandler("/api/qaq/1/qwq")
	testGetHandler("/api/qaq/1/qwq/2")
	testGetHandler("/api/asterisk/qwq/oxo")

	{
		m := tree.ExportTreeMethon("")
		if len(m) < 1 {
			t.Fatal("export must not be empty")
		}

		for _, v := range m {
			fmt.Printf("path: %s, value: %#v\n", v.Path, v.Value)
		}
	}

}

type RandomString struct {
	mu     sync.Mutex
	r      *rand.Rand
	layout string
}

var (
	Numeric = &RandomString{
		layout: "0123456789",
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
		mu:     sync.Mutex{},
	}
)

func (c *RandomString) Generate(n int) []byte {
	c.mu.Lock()
	var b = make([]byte, n)
	var length = len(c.layout)
	for i := 0; i < n; i++ {
		var idx = c.r.Intn(length)
		b[i] = c.layout[idx]
	}
	c.mu.Unlock()
	return b
}

func (c *RandomString) Intn(n int) int {
	c.mu.Lock()
	x := c.r.Intn(n)
	c.mu.Unlock()
	return x
}

// fork by: https://github.com/lxzan/uRouter/blob/main/trie_test.go#L68
func BenchmarkRouteTree_Get(b *testing.B) {
	var count = 1024
	var segmentLen = 2
	var tree = &Tree[*Value]{}
	var r = Numeric
	for i := 0; i < count; i++ {
		var idx = r.Intn(4)
		var list []string
		for j := 0; j < 4; j++ {
			var ele = string(r.Generate(segmentLen))
			if j == idx {
				ele = ":" + ele
			}
			list = append(list, ele)
		}
		tree.Set(strings.Join(list, defaultSeparator), &Value{})
	}

	var paths []string
	for i := 0; i < count; i++ {
		var path = r.Generate(12)
		path[0], path[3], path[6], path[9] = defaultSeparator[0], defaultSeparator[0], defaultSeparator[0], defaultSeparator[0]
		paths = append(paths, string(path))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var path = paths[i&(count-1)]
		tree.Get(path)
	}
}

func BenchmarkRouteTree_Set(b *testing.B) {
	var count = 1024
	var tree = &Tree[*Value]{}
	var r = Numeric
	var m = &Value{}

	var paths []string
	for i := 0; i < count; i++ {
		var path = r.Generate(12)
		path[0], path[3], path[6], path[9] = defaultSeparator[0], defaultSeparator[0], defaultSeparator[0], defaultSeparator[0]
		paths = append(paths, string(path))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var path = paths[i&(count-1)]
		tree.Set(path, m)
	}
}
