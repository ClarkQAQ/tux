<p align="center">
  <h3 align="center">Tux</h3>
  <p align="center">
    计划成为优秀的 Rest API 专用框架 (Tux 只是因为那只企鹅)
    <br />
  </p>
</p>

**目前 Tux 处于开发阶段, Rest 模块功能并未完善**

### 目标和上手指南

#### 目标

做这个是想做一个专门用于 Rest API 的框架, 因为在之前听到了不少吐槽现在 Gin + Swagger 开发 Rest 框架的声音, 所以就想做一个专门用于 Rest API 的框架, 用 Golang 的风格去约束并集成 Swagger 以及抽象 Rest 接口, 让开发者能专注于`业务逻辑/Service`的开发. 

不过由于拖延症以及最近工期的原因一直没有时间去写, 所以现在只是一个简单的 Web 框架, 也是昨天在 V2ex 看到 `uRouter` 的 Benchmark 代码后用他的代码测试出这个项目的路由树性能好像还不错, 才想起来我好像写了小半年一直没发布...

#### 上手指南

目前整个 Router 和 Context 都是 Gin 风格, 所以有一定的 Gin 使用经验的话, 上手会比较快. 不过我还没写内置 Middleware, 所以目前只能自己写 Middleware, 但是我会尽快补上的.

最简化的 Web 服务:

```go

t := tux.New()

t.Get("/", func (c *tux.Context) {
    c.String(http.StatusOK, "Hello World")
})

t.ServeAddr(":8080")

```

目前还有从郭叔叔的 Gf 框架仿写的 "Rest" 功能, 通过一个结构方法来定义一组 Method

```go

type TestApi struct{}

func (t *TestApi) Get(c *tux.Context) {
	c.String(200, "Hello, Get!")
}

func (t *TestApi) Post(c *tux.Context) {
	c.String(200, "Hello, Post!")
}


t.Object("/test", &TestApi{})

```

如果还是**没看明白的话**这里有一个完整的例子: [example](https://github.com/ClarkQAQ/tux/tree/master/_example/base)


### 愉快的 Benchmark 时间

#### 使用 uRouter 的 Tree Benchmark 测试结果以及同平台对比

> 本来想测测 Gin 的但是他好像在重复 Prefix 的时候会直接 panic, 所以就只能对比 uRouter 了...

```txt
// tux
go test -benchmem -run=^$ -bench ^(BenchmarkRouteTree_Get|BenchmarkRouteTree_Set)$ tux/tree

goos: linux
goarch: amd64
pkg: tux/tree
cpu: AMD Ryzen 7 5800H with Radeon Graphics         
BenchmarkRouteTree_Get-16    	126268011	         8.905 ns/op	       0 B/op	       0 allocs/op
BenchmarkRouteTree_Set-16    	12803094	        93.66 ns/op	       1 B/op	       0 allocs/op
PASS
ok  	tux/tree	4.429s


// github.com/lxzan/uRouter
go test -benchmem -run=^$ -bench ^BenchmarkRouteTree_Get$ github.com/lxzan/uRouter

goos: linux
goarch: amd64
pkg: github.com/lxzan/uRouter
cpu: AMD Ryzen 7 5800H with Radeon Graphics         
BenchmarkRouteTree_Get-16    	 4818972	       436.4 ns/op	      80 B/op	       1 allocs/op
PASS
ok  	github.com/lxzan/uRouter	2.366s
```


#### 使用 Gin 的 Benchmark 测试结果以及同平台对比

> 勉强摸到了 Gin 的屁股, 但是仍然无法超过 Gin...明明 Tree 比 Gin 快不少的...有时间看看你不能再去 Context 砍几刀吧...

```txt
// tux
go test -benchmem -run=^$ -bench ^(BenchmarkOneRoute|BenchmarkRecoveryMiddleware|BenchmarkLoggerMiddleware|BenchmarkManyHandlers|Benchmark5Params|BenchmarkOneRouteJSON|BenchmarkOneRoutePrintf|BenchmarkOneRouteSet|BenchmarkOneRouteString|BenchmarkManyRoutesFist|BenchmarkManyRoutesLast|Benchmark404|Benchmark404Many)$ tux

goos: linux
goarch: amd64
pkg: tux
cpu: AMD Ryzen 7 5800H with Radeon Graphics         
BenchmarkOneRoute-16              	22545145	        52.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkRecoveryMiddleware-16    	18259539	        56.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoggerMiddleware-16      	 3436378	       415.3 ns/op	      44 B/op	       3 allocs/op
BenchmarkManyHandlers-16          	 2737478	       457.1 ns/op	      44 B/op	       3 allocs/op
Benchmark5Params-16               	 9419943	       159.9 ns/op	      48 B/op	       1 allocs/op
BenchmarkOneRouteJSON-16          	 3760737	       309.7 ns/op	      32 B/op	       2 allocs/op
BenchmarkOneRoutePrintf-16        	 4094658	       328.3 ns/op	      96 B/op	       2 allocs/op
BenchmarkOneRouteSet-16           	 3351001	       336.3 ns/op	     336 B/op	       2 allocs/op
BenchmarkOneRouteString-16        	11172078	        92.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkManyRoutesFist-16        	23754633	        50.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkManyRoutesLast-16        	22015323	        55.11 ns/op	       0 B/op	       0 allocs/op
Benchmark404-16                   	11578872	       108.8 ns/op	       8 B/op	       1 allocs/op
Benchmark404Many-16               	 9436174	       116.7 ns/op	       8 B/op	       1 allocs/op
PASS
ok  	tux	19.257s


// github.com/gin-gonic/gin
go test -benchmem -run=^$ -bench ^(BenchmarkOneRoute|BenchmarkRecoveryMiddleware|BenchmarkLoggerMiddleware|BenchmarkManyHandlers|Benchmark5Params|BenchmarkOneRouteJSON|BenchmarkOneRouteHTML|BenchmarkOneRouteSet|BenchmarkOneRouteString|BenchmarkManyRoutesFist|BenchmarkManyRoutesLast|Benchmark404|Benchmark404Many)$ github.com/gin-gonic/gin

goos: linux
goarch: amd64
pkg: github.com/gin-gonic/gin
cpu: AMD Ryzen 7 5800H with Radeon Graphics         
BenchmarkOneRoute-16              	39173424	        30.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkRecoveryMiddleware-16    	33598090	        35.66 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoggerMiddleware-16      	  743853	      1524 ns/op	     220 B/op	       8 allocs/op
BenchmarkManyHandlers-16          	 1000000	      1834 ns/op	     220 B/op	       8 allocs/op
Benchmark5Params-16               	17871084	        69.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkOneRouteJSON-16          	 4122009	       388.7 ns/op	      48 B/op	       3 allocs/op
BenchmarkOneRouteHTML-16          	  945630	      2294 ns/op	     256 B/op	       9 allocs/op
BenchmarkOneRouteSet-16           	 3967929	       308.7 ns/op	     336 B/op	       2 allocs/op
BenchmarkOneRouteString-16        	 5524639	       193.7 ns/op	      48 B/op	       1 allocs/op
BenchmarkManyRoutesFist-16        	36296466	        30.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkManyRoutesLast-16        	37410177	        31.74 ns/op	       0 B/op	       0 allocs/op
Benchmark404-16                   	25092416	        44.37 ns/op	       0 B/op	       0 allocs/op
Benchmark404Many-16               	23788837	        50.29 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/gin-gonic/gin	20.164s
```