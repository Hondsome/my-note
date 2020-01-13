# 深入理解 Golang HTTP Timeout

* 深入理解 Golang HTTP Timeout

* **背景**

    > 前段时间，线上服务器因为部分微服务提供的 HTTP API 响应慢，而我们没有用正确的姿势来处理 HTTP 超时（当然，还有少量 RPC 超时）, 同时我们也没有服务降级策略和容灾机制，导致服务相继挂掉😂。服务降级和容灾需要一段时间的架构改造，但是以正确的姿势使用 HTTP 超时确是马上可以习得的。

* **超时的本质**

    > 所有的 Timeout 都构建于 Golang 提供的 Set[Read|Write]Deadline 原语之上。

## 服务器超时 server timeout

* ReadTimout 包括了TCP 消耗的时间，可以一定程度预防慢客户端和意外断开的客户端占用文件描述符
对于 https请求，ReadTimeout 包括了 TLS 握手的时间；WriteTimeout 包括了 TLS握手、读取 Header 的时间（虚线部分）, 而 http 请求只包括读取 body 和写 response 的时间。
此外，http.ListenAndServe, http.ListenAndServeTLS and http.Serve 等方法都没有设置超时，且无法设置超时。因此不适合直接用来提供公网服务。正确的姿势是：

```golang
package main

import (
    "net/http"
    "time"
)

func main() {
    server := &amp;http.Server{
        Addr:         ":8081",
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 5 * time.Second,
    }

    http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("hi"))
    })

    server.ListenAndServe()
}
```

## 客户端超时 client timeout

* http.Client 会自动跟随重定向（301, 302）, 重定向时间也会记入 http.Client.Timeout, 这点一定要注意。
取消一个 http request 有两种方式：

* Request.Cancel
Context (Golang >= 1.7.0)
后一种因为可以传递 parent context, 因此可以做级联 cancel, 效果更佳。

* 代码示例：

```golang
ctx, cancel := context.WithCancel(context.TODO())  // or parant context
timer := time.AfterFunc(5*time.Second, func() {  
    cancel()
})

req, err := http.NewRequest("GET", "http://httpbin.org/range/2048?duration=8&amp;chunk_size=256", nil)  
if err != nil {  
    log.Fatal(err)
}
req = req.WithContext(ctx)  

ctx, cancel := context.WithCancel(context.TODO())  // or parant context
timer := time.AfterFunc(5*time.Second, func() {  
    cancel()
})
 
req, err := http.NewRequest("GET", "http://httpbin.org/range/2048?duration=8&amp;chunk_size=256", nil)  
if err != nil {  
    log.Fatal(err)
}
req = req.WithContext(ctx)  
```

* Credits
The complete guide to Go net/http timeouts
Go middleware for net.Conn tracking (Prometheus/trace)