# 内存泄漏

## golang中内存泄漏的原因。
* 存在数据结构引用未被释放，导致gc不掉

* goroutine 飞了，占用数据不能被释放
    Go语言是带内存自动回收的特性，因此内存一般不会泄漏。但是Goroutine确存在泄漏的情况，同时泄漏的Goroutine引用的内存同样无法被回收。

```
func main() {
    ch := func() <-chan int {
        ch := make(chan int)
        go func() {
            for i := 0; ; i++ {
                ch <- i
            }
        } ()
        return ch
    }()

    for v := range ch {
        fmt.Println(v)
        if v == 5 {
            break
        }
    }
}
```
* 上面的程序中后台Goroutine向管道输入自然数序列，main函数中输出序列。但是当break跳出for循环的时候，后台Goroutine就处于无法被回收的状态了。

* 我们可以通过context包来避免这个问题：

```
func main() {
    ctx, cancel := context.WithCancel(context.Background())

    ch := func(ctx context.Context) <-chan int {
        ch := make(chan int)
        go func() {
            for i := 0; ; i++ {
                select {
                case <- ctx.Done():
                    return
                case ch <- i:
                }
            }
        } ()
        return ch
    }(ctx)

    for v := range ch {
        fmt.Println(v)
        if v == 5 {
            cancel()
            break
        }
    }
}
```
当main函数在break跳出循环时，通过调用cancel()来通知后台Goroutine退出，这样就避免了Goroutine的泄漏。


* 6.golang内存泄漏场景：

a. 永远处于阻塞状态的goroutine，这将导致这些 goroutine 中使用的许多代码块永远无法进行垃圾收集，如下场景案例：
如果将以下函数作为 goroutine 的启动函数并将一个 nil channel 参数传递给它, 则 goroutine 将永远阻塞. Go 运行时认为 goroutine 仍然存活, 所以为 s 分配的内存块将永远不会被收集.
func k(c <-chan bool) {
    s := make([]int64, 1e6)
    if <-c { // 如果 c 为 nil, 这里将永远阻塞
        _ = s
        // 使用 s, ...
    }
}
b.终结器(Finalizers)：为循环引用组内的成员设置 finalizer 可能会阻止为这个循环引用组分配的所有内存块被收集。
在下列函数被调用并退出之后, 为 x和 y 分配的内存块不保证在未来会被垃圾收集器回收.
func memoryLeaking() {
    type T struct {
        v [1<<20]int
        t *T
    }

    var finalizer = func(t *T) {
         fmt.Println("finalizer called")
    }
    
    var x, y T
    
    // SetFinalizer 会使 x 逃逸到堆上.
    runtime.SetFinalizer(&x, finalizer)
    
    // 以下语句将导致 x 和 y 变得无法收集.
    x.t, y.t = &y, &x // y 也逃逸到了 堆上.
}

## 排查方法

* pprof、火焰图检查

