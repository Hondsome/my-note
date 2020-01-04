1.gc垃圾回收算法：标记-清除法

基本原理：从根（包括全局指针以及goroutine栈上指针）出发，标记可达节点为灰色，然后是继续标记灰色节点可达的节点为灰色，同时自身由灰色变为黑色，不断重复，剩下所有的白色节点为垃圾回收节点

2.减轻gc负担的根本方法是：减少对象数量，特别是减少小对象的大量创建

3.goroutine调度原理：半抢占式（在函数入口处插入协作代码）的协作调度，可以理解为公平协作调度和抢占式调度的结合体

GPM模型，一个M对应一个P对应一个多个G的队列，外加一个G的全局队列
G被调度的时机有以下几种场景：

a. 简单抢占式调度：一旦某个G中出现死循环或永久循环的代码逻辑，那么G将永久占用分配给它的P和M，位于同一个P中的其他G将得不到调度，出现“饿死”的情况，在Go 1.2中实现了抢占式调度来解决这个问题，原理则是在每个函数或方法的入口，加上一段额外的代码，让runtime有机会检查是否需要执行抢占调度（是否需要调度的判断标准是由独立的sysmon线程来维护的）。这种解决方案只能说局部解决了“饿死”问题，对于没有函数调用，纯算法循环计算的G，scheduler依然无法抢占。
b. channel阻塞、锁阻塞（ atomic, mutex, 或者 channel）等阻塞场景：如果G被阻塞在某个channel操作或者G被阻塞到同步互斥操作，也就是 Lock()，Unlock() 等锁的情况下，G会被放置到某个wait队列中，原来的M不会阻塞，会尝试运行下一个runnable的G，当被放在wait队列中的G不再阻塞时，可以将G重新放到队列中等待调度执行
c. 异步system call 例如network I/O情况下的调度：当G进行network I/O操作时，G会被放置到某个wait队列中，而M会尝试运行下一个runnable的G，当被放在wait队列中的G不再阻塞时，可以将G重新放到队列中等待调度执行
c. 同步system call例如文件io操作阻塞情况下的调度：如果G被阻塞在某个system call操作上，那么不光G会阻塞，执行该G的M也会解绑P(实质是被sysmon抢走了)，与G一起进入sleep状态。如果此时有idle的M，则P与其绑定继续执行其他G；如果没有idle M，但仍然有其他G要去执行，那么就会创建一个新M。当阻塞在syscall上的G完成syscall调用后，G会去尝试获取一个可用的P，如果没有可用的P，那么G会被标记为runnable，之前的那个sleep的M将再次进入sleep。Go语言完全是自己封装的系统调用，所以在封装系统调用的时候，可以做不少手脚，也就是进入系统调用的时候执行entersyscall，退出后又执行exitsyscall函数。 也只有封装了entersyscall的系统调用才有可能触发重新调度。
4.slice与数组：slice是引用类型，数组是值类型，slice可以动态扩容，减少动态内存拷贝可以预先分配

5.channel原理：golang参考CSP模型实现的类似管道的结构，本质是一个循环队列+锁+两个Goroutine等待队列，channel是不同goroutine间通信的消息通道

channel引起的panic场景：

a. 关闭一个未初始化(nil) 的 channel ；
b. 重复关闭同一个 channel；
c. 向一个已关闭的 channel 中发送消息；
从一个已关闭的 channel 中读取消息永远不会阻塞，并且会返回一个为 false 的 ok标记，可以用它来判断 channel 是否关闭。
单向 channel 一般是在声明时会用到，防止传递进函数的channel被随意读取获取写入，限制操作，比如

func foo(ch chan<- int) <-chan int {...}
chan<- int 表示一个只可写入的 channel，<-chan int 表示一个只可读取的 channel。上面这个函数约定了 foo 内只能从向 ch 中写入数据，返回只一个只能读取的 channel，这样在方法声明时约定可以防止 channel 被滥用，这种约定在编译期间就确定下来了。
channel的实现原理：channel内部主要组成：一个环形数组实现的队列，用于存储消息元素；两个链表实现的 goroutine 等待队列，用于存储阻塞在 recv 和 send 操作上的 goroutine；一个互斥锁，用于各个属性变动的同步。

6.golang内存泄漏场景：

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
7.interface与nil

go 中的接口分为两种 
var Xxxx interface{}
type Xxxx  interface {
    Show()
}
底层结构：
type eface struct {      //空接口
    _type *_type         //类型信息
    data  unsafe.Pointer //指向数据的指针(go语言中特殊的指针类型unsafe.Pointer类似于c语言中的void*)
}
type iface struct {      //带有方法的接口
    tab  *itab           //存储type信息还有结构实现方法的集合
    data unsafe.Pointer  //指向数据的指针(go语言中特殊的指针类型unsafe.Pointer类似于c语言中的void*)
}
data指向了nil 并不代表interface 是nil
8.golang内存分配原理

new与make最终调用的都是mallocgc函数来进行分配内存，在启动时向操作系统预分配大块内存，内存分配以TCMalloc的结构进行设计, 不同大小的对象以不同class的span进行分配,减少内存碎片, 每个P都有cache减少锁冲突，同时还有zero分配采用同一个全局对象, tiny分配器加快小对象分配的优化。

9.切片会导致整个底层数组被锁定

切片会导致整个底层数组被锁定，底层数组无法释放内存。如果底层数组较大会对内存产生很大的压力。

func main() {
    headerMap := make(map[string][]byte)

    for i := 0; i < 5; i++ {
        name := "/path/to/file"
        data, err := ioutil.ReadFile(name)
        if err != nil {
            log.Fatal(err)
        }
        headerMap[name] = data[:1]
    }

    // do some thing
}
解决的方法是将结果克隆一份，这样可以释放底层的数组：

func main() {
    headerMap := make(map[string][]byte)

    for i := 0; i < 5; i++ {
        name := "/path/to/file"
        data, err := ioutil.ReadFile(name)
        if err != nil {
            log.Fatal(err)
        }
        headerMap[name] = append([]byte{}, data[:1]...)
    }

    // do some thing
}

作者：凯文不上班
链接：https://www.jianshu.com/p/33669e270139
来源：简书
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。