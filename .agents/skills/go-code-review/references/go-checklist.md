# Go 代码审查清单

## 使用方式

- 优先检查【强制】项，它们更容易对应真实缺陷、回归或线上风险。
- 【建议】项不要机械当成高严重问题，除非仓库自身已有明确约束。
- 报告 findings 时，尽量说明是命中规范条款，还是结合代码路径的推理结论。
- 如果仓库规范与本清单冲突，以仓库显式约定为准。

## 1. 风格规范

### 1.1 遵从惯例

- 【建议】遵从 Go 社区惯例，不要引入反直觉写法。

### 1.2 文件长度

- 【建议】单文件尽量不超过 800 行；超长文件通常意味着职责混杂，review 成本高。

### 1.3 缩进、括号、空格

- 【强制】统一使用 `gofmt`。
- 【强制】tab 缩进，左大括号不换行。

优秀示例：

```go
var i int = 1 + 2
v := []float64{1.0, 2.0, 3.0}[i-i]
fmt.Printf("%f\n", v+1)
```

### 1.4 标识符命名

- 【强制】标识符遵循 `MixedCaps/mixedCaps`。
- 【强制】导出标识符大写开头，非导出标识符小写开头。
- 【强制】错误变量以 `err/Err` 开头，错误类型以 `Error` 结尾。
- 【强制】receiver 命名应短、稳定，不用 `this`、`self`。
- 【强制】标识符不要重复包名语义。

反例：

```go
filterHandlerMap
uidSlice
uidArray
uidSliceSlice
```

优秀示例：

```go
opToHandler
uids
uids
classesUids
```

优秀示例：

```go
type ExitError struct {
	// ...
}

func (ri *ResearchInfo) Load() {}
func (w *ReportWriter) Write() {}
```

### 1.5 编码格式

- 【强制】代码文件必须为 UTF-8。
- 【强制】换行符使用 Unix 风格。

## 2. 注释

### 2.1 通用规则

- 【建议】单行注释与内容之间留一个空格。
- 【建议】导出名字要有注释。
- 【建议】优先使用 `//`，`/* ... */` 只用于包级文档。
- 【建议】注释符合 godoc。

### 2.2 声明注释

- 【建议】包、函数、方法、结构体、接口、全局变量、类型别名都应有清晰注释。

优秀示例：

```go
// NewAttrModel，属性数据层操作类的工厂方法。
// 返回值：属性操作类指针。
func NewAttrModel(ctx *common.Context) *AttrModel {
	// ...
}

// HasPrefix 返回 true，如果 name 包含指定 prefix。
func HasPrefix(name string, prefixes []string) bool {
	// ...
}
```

优秀示例：

```go
// User，用户实例，定义了用户的基础信息。
type User struct {
	Username string // 用户名
	Email    string // 邮箱
}
```

## 3. 错误处理

### 3.1 返回错误

- 【强制】错误统一作为最后一个返回值返回。
- 【强制】遵循 fail fast。
- 【建议】接收 `context.Context` 的函数通常也应返回 `error`。

### 3.2 日志

- 【强制】不要只打日志不返回错误，除非当前层就是最终消费错误的边界。
- 【建议】避免同一错误多层重复记录。

### 3.3 处理方式

- 【强制】不要吞掉错误。
- 【强制】出现失败后不要盲目继续执行。

### 3.4 错误字符串

- 【强制】错误字符串简洁稳定，避免首字母大写和末尾标点。

反例：

```go
errors.New("File Not Found.")
```

优秀示例：

```go
errors.New("file not found")
```

### 3.5 类型断言

- 【强制】类型断言优先使用 `v, ok := ...`。

反例：

```go
s := i.(string)
```

优秀示例：

```go
s, ok := i.(string)
if !ok {
	return errors.New("unexpected type")
}
```

### 3.6 panic / recover

- 【强制】`panic` 必须在当前 goroutine 内捕获。
- 【强制】自行启动的 goroutine 要考虑 `recover`。

优秀示例：

```go
go func() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v", r)
		}
	}()

	doWork()
}()
```

## 4. 资源管理

### 4.1 主动关闭资源

- 【强制】文件、连接、响应体、锁等资源必须释放。
- 【强制】`defer` 关闭资源时，要判断关闭错误是否重要。
- 【建议】避免在长循环里无控制地 `defer`。

反例：

```go
for _, name := range files {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
}
```

优秀示例：

```go
for _, name := range files {
	f, err := os.Open(name)
	if err != nil {
		return err
	}

	if err := handleFile(f); err != nil {
		_ = f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
}
```

### 4.2 内存泄漏

- 【建议】从大切片或长字符串上截取小片段时，关注底层内存未释放问题。
- 【建议】大流量场景下，map 删除元素后不会自动收缩，必要时重建。

反例：

```go
func head(data []byte) []byte {
	return data[:16]
}
```

优秀示例：

```go
func head(data []byte) []byte {
	out := make([]byte, 16)
	copy(out, data[:16])
	return out
}
```

### 4.3 并发

- 【强制】不要在闭包中直接捕获循环变量。
- 【强制】`map` 和 `slice` 不是并发安全的，多 goroutine 读写必须同步。

反例：

```go
ints := []int{1, 2, 3}
for _, i := range ints {
	go func() {
		fmt.Printf("%v\n", i)
	}()
}
```

优秀示例：

```go
ints := []int{1, 2, 3}
for _, i := range ints {
	go func(i int) {
		fmt.Printf("%v\n", i)
	}(i)
}
```

优秀示例：

```go
ints := []int{1, 2, 3}
for _, i := range ints {
	i := i
	go func() {
		fmt.Printf("%v\n", i)
	}()
}
```

## 5. 项目组织

### 5.1 单元测试

- 【强制】项目应提供单元测试。
- 【强制】并发相关代码测试时应考虑 `go test -race`。

### 5.2 模块化

- 【建议】项目应包含 `go.mod`，模块名、包名和目录结构保持清晰一致。

### 5.3 包导入

- 【强制】导入分组、别名和空白导入要合理，避免未使用导入和无意义别名。

## 6. 函数

### 6.1 参数

- 【建议】相同类型参数相邻。
- 【建议】优先值传递；需要修改对象、对象很大、或含同步原语时再考虑指针。
- 【建议】不要给 `map`、`slice`、`chan`、`interface` 传指针。
- 【建议】优先用切片而不是数组作为参数。

### 6.2 返回值

- 【强制】不要返回多个仅用于流程控制的状态。
- 【强制】返回值个数不超过 3 个。

反例：

```go
isContinue, retCode := p.processUnity()
```

优秀示例：

```go
retCode := p.processUnity()
```

### 6.3 实现

- 【建议】把 `map` 或 `slice` 存入结构体时，考虑浅拷贝和外部修改风险。

反例：

```go
func (d *Driver) SetTrips(trips []Trip) {
	d.trips = trips
}
```

优秀示例：

```go
func (d *Driver) SetTrips(trips []Trip) {
	d.trips = append([]Trip(nil), trips...)
}
```

### 6.4 分组

- 【建议】构造函数靠前，导出函数在前，同一 receiver 靠拢，工具函数放末尾。

反例：

```go
func (s *something) Cost() {
	return calcCost(s.weights)
}

type something struct{ /* ... */ }

func calcCost(n []int) int { return 0 }

func (s *something) Stop() {}

func newSomething() *something { return &something{} }
```

优秀示例：

```go
type something struct{ /* ... */ }

func newSomething() *something { return &something{} }

func (s *something) Cost() {
	return calcCost(s.weights)
}

func (s *something) Stop() {}

func calcCost(n []int) int { return 0 }
```

## 7. 语言特性

### 7.1 变量声明

- 【建议】变量就近声明并缩小作用域。
- 【建议】局部变量优先用 `:=`，零值局部变量优先用 `var`。
- 【建议】相似声明使用分组。

反例：

```go
const a = 1
const b = 2

var x = 1
var y = 2
```

优秀示例：

```go
const (
	a = 1
	b = 2
)

var (
	x = 1
	y = 2
)
```

反例：

```go
func foo() {
	var s = "foo"
}
```

优秀示例：

```go
func foo() {
	s := "foo"
}
```

反例：

```go
err := file.Chmod(0664)
if err != nil {
	return err
}
```

优秀示例：

```go
if err := file.Chmod(0664); err != nil {
	return err
}
```

### 7.2 语句

- 【强制】禁止变量遮蔽。
- 【强制】指针解引用前确认非 nil。
- 【强制】外部输入参与数值运算时，注意溢出、截断和符号错误。

### 7.3 魔数

- 【强制】除 `0` 和 `1` 外，避免直接使用魔数。
- 【强制】重复出现的字符串应提取常量。

### 7.4 for / range / switch / goto

- 【建议】不要持有循环变量地址。
- 【强制】只需要 `range` 的 key 时丢弃 value，只需要 value 时把 key 写成 `_`。
- 【强制】遍历大结构体切片优先按下标访问。
- 【强制】`switch` 必须有 `default`。
- 【强制】业务代码禁止使用 `goto`。

反例：

```go
for _, v := range ints {
	fmt.Println(&v)
}
```

优秀示例：

```go
for _, v := range ints {
	v := v
	fmt.Println(&v)
}
```

优秀示例：

```go
for key := range m {
	if key.expired() {
		delete(m, key)
	}
}
```

优秀示例：

```go
sum := 0
for _, v := range array {
	sum += v
}
```

反例：

```go
for _, item := range largeStructs {
	item.Age = 10
}
```

优秀示例：

```go
for i := range largeStructs {
	largeStructs[i].Age = 10
}
```

### 7.5 结构体

- 【强制】字段排序考虑内存对齐。
- 【强制】结构体拷贝要考虑浅拷贝问题。
- 【强制】结构体初始化使用带字段名的字面量。
- 【建议】结构体指针初始化优先 `&T{}`。

反例：

```go
user := User{"alice", 18}
```

优秀示例：

```go
user := User{
	Name: "alice",
	Age:  18,
}
```

### 7.6 slice

- 【强制】`slice` 应合理初始化，必要时预估容量。
- 【强制】索引前确保不越界。
- 【强制】`copy` 前确认目标长度足够。
- 【建议】空切片优先 `var s []T`。
- 【建议】零长度返回值优先 `nil`，空判断统一 `len(s) == 0`。

反例：

```go
nums := []int{}
```

优秀示例：

```go
nums := make([]int, 0, n)
```

反例：

```go
src := []int{0, 1, 2}
var dst []int
copy(dst, src)
```

优秀示例：

```go
src := []int{0, 1, 2}
dst := make([]int, len(src))
copy(dst, src)
```

反例：

```go
if x == "" {
	return []int{}
}

func isEmpty(s []string) bool {
	return s == nil
}
```

优秀示例：

```go
if x == "" {
	return nil
}

func isEmpty(s []string) bool {
	return len(s) == 0
}
```

### 7.7 map

- 【强制】固定元素列表直接用字面量，其他场景用 `make`。
- 【强制】可预估容量时显式给容量。
- 【强制】判断 key 是否存在要用 `v, ok := m[k]`。
- 【强制】禁止向 nil map 写入。

反例：

```go
m := make(map[T1]T2, 3)
m[k1] = v1
m[k2] = v2
m[k3] = v3
```

优秀示例：

```go
m := map[T1]T2{
	k1: v1,
	k2: v2,
	k3: v3,
}
```

### 7.8 channel

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

### 7.9 goroutine

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

### 7.10 context

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

### 7.11 init

- 【强制】避免复杂 `init()` 逻辑，尤其不要在 `init()` 中启动后台 goroutine。

### 7.12 interface

- 【建议】接口应尽量小，由消费者定义。

### 7.13 Mutex

- 【强制】`sync.Mutex` 不可复制，含锁结构体传值和拷贝都要谨慎。

### 7.14 序列化

- 【强制】序列化字段、tag、默认值语义必须准确，避免“看起来生效、实际无效”。

优秀示例：

```go
type Response struct {
	Message string `json:"message"`
}
```

### 7.15 time

- 【强制】时间运算、时区和超时处理要明确。

优秀示例：

```go
ctx, cancel := context.WithTimeout(parent, 100*time.Millisecond)
defer cancel()
```

### 7.16 atomic / rand / unsafe / cgo

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

## 8. 国际化

- 【强制】中英文文本使用对应语言的标点规范。
- 【强制】如果项目采用统一翻译函数，如 `T()`，则所有待翻译字符串必须走统一入口。
- 【强制】翻译词条不应直接使用格式化后的字符串。

## 9. 工具检查

- 【强制】每个 Go 模块都应配置并通过 `golangci-lint`。
- 【强制】至少关注这些高价值检查器：`bodyclose`、`errcheck`、`gocritic`、`gocyclo`、`gofmt`、`goimports`、`gomnd`、`gosec`、`gosimple`、`govet`、`ineffassign`、`nolintlint`、`revive`、`staticcheck`、`stylecheck`、`typecheck`、`unused`。
- 【建议】把工具结果和代码推理分开陈述：工具证明的问题直接引用，工具未覆盖的问题补推理链路。

## Review 输出建议

- 优先输出命中【强制】条款且会造成真实 bug、竞态、资源泄漏、安全问题、协议不一致的问题。
- 命中【建议】条款但只影响风格时，除非仓库明确要求，否则降级处理。
