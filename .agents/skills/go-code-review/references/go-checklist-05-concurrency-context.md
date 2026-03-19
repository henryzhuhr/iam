# Go 代码审查清单 05: 并发、上下文与底层能力

本文件聚焦 channel、goroutine、context、atomic、Mutex、随机数、unsafe 和 cgo。

## 05. 并发、上下文与底层能力

### 05.1 channel

- 【强制】channel 通常应无缓冲或容量为 1，其他容量要解释。
- 【强制】channel 必须 `make` 初始化。
- 【强制】多生产者场景下要明确关闭时机，禁止重复关闭。
- 【建议】限制读写权限。

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

优秀示例：

```go
func producer(ch chan<- int) {
	ch <- 1
	close(ch)
}

func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}
```

### 05.2 goroutine

- 【强制】goroutine 必须可等待退出或取消，不能泄漏。
- 【强制】协程池要限制最大并发数。
- 【强制】不要在 `init()` 中启动 goroutine。

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

优秀示例：

```go
done := make(chan struct{})
go func() {
	defer close(done)
	doWork()
}()
<-done
```

### 05.3 context

- 【强制】`context.Context` 始终为第一个参数。
- 【强制】不要传 `nil` context。
- 【强制】`WithTimeout`、`WithDeadline`、`WithCancel` 返回的 `cancel` 必须执行。
- 【建议】RPC、数据库、外部 IO 要设置超时。

优秀示例：

```go
func slowOperationWithTimeout(ctx context.Context) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	return slowOperation(ctx)
}
```

### 05.4 Mutex

- 【强制】`sync.Mutex` 不可复制，含锁结构体传值和拷贝都要谨慎。

### 05.5 atomic / rand / unsafe / cgo

- 【强制】原子变量必须始终原子化读写。
- 【建议】如果项目已采用 `go.uber.org/atomic`，优先沿用其类型安全封装。
- 【强制】安全相关随机数使用 `crypto/rand`。
- 【强制】禁止随意使用 `unsafe`。
- 【强制】`cgo` 场景要避免把 Go 内存长期交给 C 持有。

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
