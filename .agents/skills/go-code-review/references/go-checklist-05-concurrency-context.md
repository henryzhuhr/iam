# Go 代码审查清单 05: 并发、上下文与底层能力

本文件聚焦 channel、goroutine、context、atomic、锁、随机数、unsafe 和 cgo。并发问题通常比风格问题更值得优先报告，因为它们在测试里不容易稳定复现，却可能造成线上故障。

## 05. 并发、上下文与底层能力

### 05.1 channel

- channel 的缓冲大小应该服务于协议、背压和吞吐需求，不存在“通常只能是无缓冲或 1”这样的通用规则。
- nil channel 既可能是 bug，也可能是 `select` 状态机里用来动态关闭某个 case 的技巧；review 要看生命周期是否清楚，而不是看到 nil 就判错。
- 多生产者、多消费者和关闭 channel 的所有权必须明确，尤其要避免重复关闭、发送到已关闭 channel、永远没人接收或永远没人关闭。
- 方向限定的 `chan<-` / `<-chan` 有助于表达契约，但它是接口设计建议，不是所有局部变量都必须这么写。

反例：

```go
var ch chan int
go func() {
	ch <- 1
}()
```

优秀示例：

```go
ch := make(chan int)
go func() {
	ch <- 1
}()
fmt.Println(<-ch)
```

### 05.2 goroutine

- 每个 goroutine 都应该有清楚的退出条件，或者明确说明它为何需要伴随进程长期存活。
- 是否需要等待 goroutine 结束，要结合调用语义判断；请求作用域、测试、批处理和优雅退出路径通常都需要等待或取消。
- 无界地为输入数据启动 goroutine、在 `init()` 或包级变量初始化中偷偷启动后台任务，通常都值得重点审查。

优秀示例：

```go
var wg sync.WaitGroup
for i := 0; i < n; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		doWork()
	}()
}
wg.Wait()
```

### 05.3 context

- `context.Context` 应从调用链上传递下来，不要在请求处理中随意替换成 `context.Background()` 或丢弃取消信号。
- 对新设计的函数，`context.Context` 通常放第一个参数；如果是在实现既有接口，就要先服从接口定义。
- `WithTimeout`、`WithDeadline`、`WithCancel` 创建出的 `cancel` 只要由当前函数持有，就应在合适时机调用，以释放关联资源。
- RPC、数据库、外部 HTTP、消息队列和后台任务通常都值得显式 timeout 或 deadline，但具体阈值要依赖业务语义。

优秀示例：

```go
func slowOperationWithTimeout(ctx context.Context) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	return slowOperation(ctx)
}
```

### 05.4 Mutex 与 atomic

- `sync.Mutex`、`sync.RWMutex`、`sync.WaitGroup`、`sync.Once` 以及原子类型在首次使用后都不应再被复制。
- 同一个状态如果一部分路径使用 atomic、一部分路径直接普通读写，通常就是数据竞争或内存模型错误。
- atomic 适合保护独立变量；如果多个字段需要一起保持不变量，锁通常比原子变量更安全、更容易推理。

反例：

```go
type foo struct {
	running int32
}

func (f *foo) isRunning() bool {
	return f.running == 1
}
```

优秀示例：

```go
type foo struct {
	running atomic.Bool
}

func (f *foo) isRunning() bool {
	return f.running.Load()
}
```

### 05.5 rand / unsafe / cgo

- 安全相关随机数必须使用 `crypto/rand`；`math/rand` 只适合非安全场景。
- `unsafe` 和 cgo 不是绝对禁区，但它们都要求更强的不变量说明：内存所有权、生命周期、对齐、可移动性和并发约束必须清楚。
- 只要 `unsafe` 或 cgo 代码缺少边界说明、测试或封装隔离，就值得提高审查强度。
